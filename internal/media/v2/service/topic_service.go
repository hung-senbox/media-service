package service

import (
	"context"
	"fmt"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/redis"
	"media-service/pkg/constants"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicService interface {
	CreateParentTopic(ctx context.Context, req request.CreateTopicRequest) (*model.Topic, error)
	UpdateParentTopic(ctx context.Context, req request.UpdateTopicRequest) error
	GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error)
	GetParentTopics4Web(ctx context.Context) ([]response.TopicResponse, error)
}

type topicService struct {
	userGateway  gateway.UserGateway
	fileGateway  gateway.FileGateway
	redisService *redis.RedisService
	topicRepo    repository.TopicRepository
}

func NewTopicService(topicRepo repository.TopicRepository, fileGw gateway.FileGateway, redisService *redis.RedisService, userGw gateway.UserGateway) TopicService {
	return &topicService{
		topicRepo:    topicRepo,
		redisService: redisService,
		fileGateway:  fileGw,
		userGateway:  userGw,
	}
}

// ------------------- Create Topic -------------------
func (s *topicService) CreateParentTopic(ctx context.Context, req request.CreateTopicRequest) (*model.Topic, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}
	orgID := currentUser.OrganizationAdmin.ID

	inProgress, err := s.redisService.HasAnyUploadInProgress(ctx, orgID)
	if err == nil && inProgress {
		return nil, fmt.Errorf("please wait until the previous upload is completed")
	}

	topic := &model.Topic{
		ID:             primitive.NewObjectID(),
		OrganizationID: orgID,
		IsPublished:    req.IsPublished,
		LanguageConfig: []model.TopicLanguageConfig{},
	}

	// upsert LanguageConfig
	langConfig := model.TopicLanguageConfig{
		LanguageID:  req.LanguageID,
		FileName:    req.FileName,
		Title:       req.Title,
		Note:        req.Note,
		Description: req.Description,
		Images:      []model.TopicImageConfig{},
		Video:       model.TopicVideoConfig{},
		Audio:       model.TopicAudioConfig{},
	}

	if err := s.topicRepo.CreateTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("create topic fail: %w", err)
	}
	if err := s.topicRepo.SetLanguageConfig(ctx, topic.ID.Hex(), langConfig); err != nil {
		return nil, fmt.Errorf("upsert language config fail: %w", err)
	}

	// Tính total task
	totalTasks := 0
	files := []interface{}{
		req.AudioFile, req.VideoFile, req.FullBackgroundFile, req.ClearBackgroundFile,
		req.ClipPartFile, req.DrawingFile, req.IconFile, req.BMFile,
	}
	for _, f := range files {
		if f != nil {
			totalTasks++
		}
	}

	if totalTasks > 0 {
		_ = s.redisService.InitUploadProgress(ctx, orgID, topic.ID.Hex(), totalTasks)
	}

	// Upload async
	go func(orgID, topicID string, req request.CreateTopicRequest) {
		ctxUpload, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		token, _ := ctx.Value(constants.Token).(string)
		ctxUpload = context.WithValue(ctxUpload, constants.Token, token)

		s.uploadFilesAsync(ctxUpload, orgID, topicID, req)
	}(orgID, topic.ID.Hex(), req)

	return topic, nil
}

// ------------------- Upload Async -------------------
func (s *topicService) uploadFilesAsync(ctx context.Context, orgID, topicID string, req request.CreateTopicRequest) {
	decrementTask := func() {
		remaining, _ := s.redisService.DecrementUploadTask(ctx, orgID, topicID)
		if remaining <= 0 {
			_ = s.redisService.SetUploadProgress(ctx, orgID, topicID, 0) // progress = 100%
		}
	}

	if req.AudioFile != nil {
		time.Sleep(4 * time.Second)
		s.uploadAndSaveAudio(ctx, orgID, topicID, req, decrementTask)
	}
	if req.VideoFile != nil {
		time.Sleep(4 * time.Second)
		s.uploadAndSaveVideo(ctx, orgID, topicID, req, decrementTask)
	}
	s.uploadAndSaveImages(ctx, orgID, topicID, req, decrementTask)
}

// --- Upload audio ---
func (s *topicService) uploadAndSaveAudio(ctx context.Context, orgID, topicID string, req request.CreateTopicRequest, decrementTask func()) {
	resp, err := s.fileGateway.UploadAudio(ctx, gw_request.UploadFileRequest{
		File:     req.AudioFile,
		Folder:   "topic_media",
		FileName: req.Title + "_audio",
		Mode:     "private",
	})
	if err != nil {
		_ = s.redisService.SetUploadError(ctx, orgID, topicID, "audio_error", err.Error())
	} else {
		_ = s.topicRepo.SetAudio(ctx, topicID, req.LanguageID, model.TopicAudioConfig{
			AudioKey:  resp.Key,
			LinkUrl:   req.AudioLinkUrl,
			StartTime: req.AudioStart,
			EndTime:   req.AudioEnd,
		})
	}
	decrementTask()
}

// --- Upload video ---
func (s *topicService) uploadAndSaveVideo(ctx context.Context, orgID, topicID string, req request.CreateTopicRequest, decrementTask func()) {
	resp, err := s.fileGateway.UploadVideo(ctx, gw_request.UploadFileRequest{
		File:     req.VideoFile,
		Folder:   "topic_media",
		FileName: req.Title + "_video",
		Mode:     "private",
	})
	if err != nil {
		_ = s.redisService.SetUploadError(ctx, orgID, topicID, "video_error", err.Error())
	} else {
		_ = s.topicRepo.SetVideo(ctx, topicID, req.LanguageID, model.TopicVideoConfig{
			VideoKey:  resp.Key,
			LinkUrl:   req.VideoLinkUrl,
			StartTime: req.VideoStart,
			EndTime:   req.VideoEnd,
		})
	}
	decrementTask()
}

// --- Upload images ---
func (s *topicService) uploadAndSaveImages(ctx context.Context, orgID, topicID string, req request.CreateTopicRequest, decrementTask func()) {
	imageFiles := []struct {
		file *multipart.FileHeader
		link string
		typ  string
	}{
		{req.FullBackgroundFile, req.FullBackgroundLink, "full_background"},
		{req.ClearBackgroundFile, req.ClearBackgroundLink, "clear_background"},
		{req.ClipPartFile, req.ClipPartLink, "clip_part"},
		{req.DrawingFile, req.DrawingLink, "drawing"},
		{req.IconFile, req.IconLink, "icon"},
		{req.BMFile, req.BMLink, "bm"},
	}

	for _, img := range imageFiles {
		time.Sleep(2 * time.Second)
		if img.file == nil {
			continue
		}
		resp, err := s.fileGateway.UploadImage(ctx, gw_request.UploadFileRequest{
			File:      img.file,
			Folder:    "topic_media",
			FileName:  req.Title + "_" + img.typ + "_image",
			ImageName: img.typ,
			Mode:      "private",
		})
		if err != nil {
			_ = s.redisService.SetUploadError(ctx, orgID, topicID, "image_"+img.typ, err.Error())
		} else {
			_ = s.topicRepo.SetImage(ctx, topicID, req.LanguageID, model.TopicImageConfig{
				ImageKey:  resp.Key,
				ImageType: img.typ,
				LinkUrl:   img.link,
			})
		}
		decrementTask()
	}
}

// ------------------- Get upload progress -------------------
func (s *topicService) GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error) {

	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}

	organizationID := currentUser.OrganizationAdmin.ID
	total, _ := s.redisService.GetTotalUploadTask(ctx, organizationID, topicID)
	remaining, _ := s.redisService.GetUploadProgress(ctx, organizationID, topicID)

	progress := 0
	if total > 0 {
		progress = int((total - remaining) * 100 / total)
		if progress > 100 {
			progress = 100
		}
	}

	rawErrors, _ := s.redisService.GetUploadErrors(ctx, organizationID, topicID)

	imageErr := map[string]string{
		"full_background":  rawErrors["image_full_background"],
		"clear_background": rawErrors["image_clear_background"],
		"clip_part":        rawErrors["image_clip_part"],
		"drawing":          rawErrors["image_drawing"],
		"icon":             rawErrors["image_icon"],
		"bm":               rawErrors["image_bm"],
	}

	return &response.GetUploadProgressResponse{
		Progress: progress,
		UploadErrors: map[string]any{
			"audio_error": rawErrors["audio_error"],
			"video_error": rawErrors["video_error"],
			"image_error": imageErr,
		},
	}, nil
}

func (s *topicService) GetParentTopics4Web(ctx context.Context) ([]response.TopicResponse, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}
	orgID := currentUser.OrganizationAdmin.ID
	topics, err := s.topicRepo.GetAllParentByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return mapper.ToTopicResponses(topics), nil

}

func (s *topicService) UpdateParentTopic(ctx context.Context, req request.UpdateTopicRequest) error {
	// Lấy topic hiện tại
	objID, err := primitive.ObjectIDFromHex(req.TopicID)
	if err != nil {
		return fmt.Errorf("invalid topic id: %w", err)
	}

	oldTopic, err := s.topicRepo.GetByID(ctx, objID)
	if err != nil {
		return fmt.Errorf("get topic failed: %w", err)
	}

	// Merge dữ liệu
	for i, lc := range oldTopic.LanguageConfig {
		if lc.LanguageID == req.LanguageID {
			oldTopic.LanguageConfig[i].FileName = req.FileName
			oldTopic.LanguageConfig[i].Title = req.Title
			oldTopic.LanguageConfig[i].Note = req.Node
			oldTopic.LanguageConfig[i].Description = req.Description
		}
	}

	// Update published
	oldTopic.IsPublished = req.IsPublished

	// Call repo update
	return s.topicRepo.UpdateTopic(ctx, oldTopic)
}
