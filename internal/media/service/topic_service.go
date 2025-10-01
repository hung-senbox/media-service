package service

import (
	"context"
	"fmt"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/media/dto/request"
	"media-service/internal/media/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/repository"
	"media-service/internal/redis"
	"media-service/pkg/constants"
	"mime/multipart"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicService interface {
	UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error)
	GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error)
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

func (s *topicService) UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	currentUser, err := s.userGateway.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}

	topicID := req.TopicID
	if topicID == "" {
		topicID = primitive.NewObjectID().Hex()
	}

	// Tạo topic
	oid, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return nil, fmt.Errorf("invalid topicID: %w", err)
	}

	topic := &model.Topic{
		ID:             oid,
		OrganizationID: currentUser.OrganizationAdmin.ID,
		ParentID:       req.ParentID,
		IsPublished:    req.IsPublished,
		LanguageConfig: []model.TopicLanguageConfig{}, // khởi tạo rỗng
	}

	// upsert LanguageConfig
	langConfig := model.TopicLanguageConfig{
		LanguageID:  req.LanguageID,
		FileName:    req.FileName,
		Title:       req.Title,
		Note:        req.Note,
		Description: req.Description,
		Images:      []model.TopicImageConfig{},
		Videos:      []model.TopicVideoConfig{},
		Audios:      []model.TopicAudioConfig{},
	}

	if err := s.topicRepo.UploadTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("create topic fail: %w", err)
	}

	if err := s.topicRepo.AddLanguageConfigToTopic(ctx, topic.ID.Hex(), langConfig); err != nil {
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
		_ = s.redisService.InitUploadProgress(ctx, topic.ID.Hex(), totalTasks)
	}

	// Upload async
	go func(topicID string, req request.UploadTopicRequest) {
		ctxUpload, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		token, _ := ctx.Value(constants.Token).(string)
		ctxUpload = context.WithValue(ctxUpload, constants.Token, token)

		s.uploadFilesAsync(ctxUpload, topicID, req)
	}(topic.ID.Hex(), req)

	return topic, nil
}

// --- Upload async ---
func (s *topicService) uploadFilesAsync(ctx context.Context, topicID string, req request.UploadTopicRequest) {
	decrementTask := func() {
		remaining, _ := s.redisService.DecrementUploadTask(ctx, topicID)
		if remaining <= 0 {
			_ = s.redisService.SetUploadProgress(ctx, topicID, 0) // progress = 100%
		}
	}

	if req.AudioFile != nil {
		s.uploadAndSaveAudio(ctx, topicID, req, decrementTask)
	}
	if req.VideoFile != nil {
		s.uploadAndSaveVideo(ctx, topicID, req, decrementTask)
	}
	s.uploadAndSaveImages(ctx, topicID, req, decrementTask)
}

// --- Upload audio ---
func (s *topicService) uploadAndSaveAudio(ctx context.Context, topicID string, req request.UploadTopicRequest, decrementTask func()) {
	if req.AudioFile == nil {
		return
	}
	resp, err := s.fileGateway.UploadAudio(ctx, gw_request.UploadFileRequest{
		File:     req.AudioFile,
		Folder:   "topic_audios",
		FileName: req.Title,
		Mode:     "private",
	})
	if err != nil {
		_ = s.redisService.SetUploadError(ctx, topicID, "audio_error", err.Error())
	} else {
		_ = s.topicRepo.AddAudioToTopic(ctx, topicID, req.LanguageID, req.AudioOldKey, model.TopicAudioConfig{
			AudioKey:  resp.Key,
			LinkUrl:   req.AudioLinkUrl,
			StartTime: req.AudioStart,
			EndTime:   req.AudioEnd,
		})
	}
	decrementTask()
}

// --- Upload video ---
func (s *topicService) uploadAndSaveVideo(ctx context.Context, topicID string, req request.UploadTopicRequest, decrementTask func()) {
	if req.VideoFile == nil {
		return
	}

	resp, err := s.fileGateway.UploadVideo(ctx, gw_request.UploadFileRequest{
		File:     req.VideoFile,
		Folder:   "topic_videos",
		FileName: req.Title,
		Mode:     "private",
	})
	if err != nil {
		_ = s.redisService.SetUploadError(ctx, topicID, "video_error", err.Error())
	} else {
		_ = s.topicRepo.AddVideoToTopic(ctx, topicID, req.LanguageID, req.VideoOldKey, model.TopicVideoConfig{
			VideoKey:  resp.Key,
			LinkUrl:   req.VideoLinkUrl,
			StartTime: req.VideoStart,
			EndTime:   req.VideoEnd,
		})
	}

	decrementTask()
}

// --- Upload images ---
func (s *topicService) uploadAndSaveImages(ctx context.Context, topicID string, req request.UploadTopicRequest, decrementTask func()) {
	imageFiles := []struct {
		file   *multipart.FileHeader
		link   string
		oldKey string
		typ    string
	}{
		{req.FullBackgroundFile, req.FullBackgroundLink, req.FullBackgroundOldKey, "full_background"},
		{req.ClearBackgroundFile, req.ClearBackgroundLink, req.ClearBackgroundOldKey, "clear_background"},
		{req.ClipPartFile, req.ClipPartLink, req.ClipPartOldKey, "clip_part"},
		{req.DrawingFile, req.DrawingLink, req.DrawingOldKey, "drawing"},
		{req.IconFile, req.IconLink, req.IconOldKey, "icon"},
		{req.BMFile, req.BMLink, req.BMOldKey, "bm"},
	}

	for _, img := range imageFiles {
		if img.file == nil {
			continue
		}

		resp, err := s.fileGateway.UploadImage(ctx, gw_request.UploadFileRequest{
			File:      img.file,
			Folder:    "topic_images",
			FileName:  img.file.Filename,
			ImageName: img.typ,
			Mode:      "private",
		})
		if err != nil {
			_ = s.redisService.SetUploadError(ctx, topicID, "image_"+img.typ, err.Error())
		} else {
			// Add hoặc update image dựa trên oldKey và languageID
			_ = s.topicRepo.AddImageToTopic(ctx, topicID, req.LanguageID, img.typ, img.oldKey, model.TopicImageConfig{
				ImageKey:  resp.Key,
				ImageType: img.typ,
				LinkUrl:   img.link,
			})
		}

		decrementTask()
	}
}

// --- Get upload progress ---
func (s *topicService) GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error) {
	total, _ := s.redisService.GetTotalUploadTask(ctx, topicID)
	remaining, _ := s.redisService.GetUploadProgress(ctx, topicID)

	progress := 0
	if total > 0 {
		progress = int((total - remaining) * 100 / total)
		if progress > 100 {
			progress = 100
		}
	}

	rawErrors, _ := s.redisService.GetUploadErrors(ctx, topicID)

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
