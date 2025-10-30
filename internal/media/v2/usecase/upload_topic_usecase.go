package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

	"media-service/helper"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/repository"
	"media-service/internal/redis"
	"media-service/pkg/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadTopicUseCase interface {
	UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error)
}

type uploadTopicUseCase struct {
	fileGateway  gateway.FileGateway
	redisService *redis.RedisService
	topicRepo    repository.TopicRepository
}

func NewUploadTopicUseCase(topicRepo repository.TopicRepository, fileGw gateway.FileGateway, redisService *redis.RedisService) UploadTopicUseCase {
	return &uploadTopicUseCase{
		topicRepo:    topicRepo,
		redisService: redisService,
		fileGateway:  fileGw,
	}
}

// ------------------- UploadTopic main flow -------------------
func (uc *uploadTopicUseCase) UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser == nil || currentUser.OrganizationAdmin.ID == "" || currentUser.IsSuperAdmin {
		return nil, fmt.Errorf("access denied")
	}
	orgID := currentUser.OrganizationAdmin.ID

	// Check upload in progress
	inProgress, err := uc.redisService.HasAnyUploadInProgress(ctx, orgID)
	if err == nil && inProgress {
		return nil, fmt.Errorf("please wait until the previous upload is completed")
	}

	var topic *model.Topic

	if req.TopicID != "" {
		// Case update existing topic
		topic, err = uc.updateTopicLanguage(ctx, req)
		if err != nil {
			return nil, err
		}
	} else {
		// Case create new topic
		topic, err = uc.createTopicLanguage(ctx, req, orgID)
		if err != nil {
			return nil, err
		}
	}

	// Count valid upload files
	files := []*multipart.FileHeader{
		req.AudioFile, req.VideoFile, req.FullBackgroundFile,
		req.ClearBackgroundFile, req.ClipPartFile,
		req.DrawingFile, req.IconFile, req.BMFile,
	}

	totalTasks := 0
	for _, f := range files {
		if helper.IsValidFile(f) {
			totalTasks++
		}
	}

	if totalTasks > 0 {
		_ = uc.redisService.InitUploadProgress(ctx, orgID, topic.ID.Hex(), totalTasks)
		go uc.uploadFilesAsyncWithContext(ctx, orgID, topic.ID.Hex(), req)
	}

	return topic, nil
}

// ------------------- Upload async -------------------
func (uc *uploadTopicUseCase) uploadFilesAsyncWithContext(ctx context.Context, orgID, topicID string, req request.UploadTopicRequest) {
	ctxUpload := context.Background()

	if token, ok := ctx.Value(constants.Token).(string); ok {
		ctxUpload = context.WithValue(ctxUpload, constants.Token, token)
	}

	uc.uploadFilesAsync(ctxUpload, orgID, topicID, req)
}

func (uc *uploadTopicUseCase) uploadFilesAsync(ctx context.Context, orgID, topicID string, req request.UploadTopicRequest) {
	decrementTask := func() {
		remaining, _ := uc.redisService.DecrementUploadTask(ctx, orgID, topicID)
		if remaining <= 0 {
			_ = uc.redisService.SetUploadProgress(ctx, orgID, topicID, 0)
		}
	}

	if helper.IsValidFile(req.AudioFile) {
		uc.uploadAndSaveAudio(ctx, orgID, topicID, req, decrementTask)
	}

	if helper.IsValidFile(req.VideoFile) {
		uc.uploadAndSaveVideo(ctx, orgID, topicID, req, decrementTask)
	}

	uc.uploadAndSaveImages(ctx, orgID, topicID, req, decrementTask)
}

// ------------------- Upload handlers -------------------
func (uc *uploadTopicUseCase) uploadAndSaveAudio(ctx context.Context, orgID, topicID string, req request.UploadTopicRequest, done func()) {
	defer done()

	// case 1: co file
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		_ = uc.redisService.SetUploadError(ctx, orgID, topicID, "audio_error", err.Error())
		return
	}
	oldAudioKey := uc.getAudioKeyByLanguage(topic, req.LanguageID)
	// case 2: khong co file
	resp, err := uc.fileGateway.UploadAudio(ctx, gw_request.UploadFileRequest{
		File:     req.AudioFile,
		Folder:   "topic_media",
		FileName: req.Title + "_audio",
		Mode:     "private",
	})
	if err != nil {
		_ = uc.redisService.SetUploadError(ctx, orgID, topicID, "audio_error", err.Error())
		return
	}

	_ = uc.topicRepo.SetAudio(ctx, topicID, req.LanguageID, model.TopicAudioConfig{
		AudioKey:  resp.Key,
		LinkUrl:   req.AudioLinkUrl,
		StartTime: req.AudioStart,
		EndTime:   req.AudioEnd,
	})
}

func (uc *uploadTopicUseCase) uploadAndSaveVideo(ctx context.Context, orgID, topicID string, req request.UploadTopicRequest, done func()) {
	defer done()

	resp, err := uc.fileGateway.UploadVideo(ctx, gw_request.UploadFileRequest{
		File:     req.VideoFile,
		Folder:   "topic_media",
		FileName: req.Title + "_video",
		Mode:     "private",
	})
	if err != nil {
		_ = uc.redisService.SetUploadError(ctx, orgID, topicID, "video_error", err.Error())
		return
	}

	_ = uc.topicRepo.SetVideo(ctx, topicID, req.LanguageID, model.TopicVideoConfig{
		VideoKey:  resp.Key,
		LinkUrl:   req.VideoLinkUrl,
		StartTime: req.VideoStart,
		EndTime:   req.VideoEnd,
	})
}

func (uc *uploadTopicUseCase) uploadAndSaveImages(ctx context.Context, orgID, topicID string, req request.UploadTopicRequest, done func()) {
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
		if !helper.IsValidFile(img.file) {
			continue
		}

		resp, err := uc.fileGateway.UploadImage(ctx, gw_request.UploadFileRequest{
			File:      img.file,
			Folder:    "topic_media",
			FileName:  req.Title + "_" + img.typ + "_image",
			ImageName: img.typ,
			Mode:      "private",
		})
		if err != nil {
			_ = uc.redisService.SetUploadError(ctx, orgID, topicID, "image_"+img.typ, err.Error())
			done()
			continue
		}

		_ = uc.topicRepo.SetImage(ctx, topicID, req.LanguageID, model.TopicImageConfig{
			ImageKey:  resp.Key,
			ImageType: img.typ,
			LinkUrl:   img.link,
		})

		done()
	}
}

// ------------------- Topic helpers -------------------
func (uc *uploadTopicUseCase) updateTopicLanguage(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	oldTopic, err := uc.topicRepo.GetByID(ctx, req.TopicID)
	if err != nil {
		return nil, fmt.Errorf("get topic failed: %w", err)
	}

	found := false
	for i, lc := range oldTopic.LanguageConfig {
		if lc.LanguageID == req.LanguageID {
			oldTopic.LanguageConfig[i].FileName = req.FileName
			oldTopic.LanguageConfig[i].Title = req.Title
			oldTopic.LanguageConfig[i].Note = req.Note
			oldTopic.LanguageConfig[i].Description = req.Description
			found = true
			break
		}
	}

	if !found {
		oldTopic.LanguageConfig = append(oldTopic.LanguageConfig, model.TopicLanguageConfig{
			LanguageID:  req.LanguageID,
			FileName:    req.FileName,
			Title:       req.Title,
			Note:        req.Note,
			Description: req.Description,
		})
	}

	oldTopic.IsPublished = req.IsPublished
	return uc.topicRepo.UpdateTopic(ctx, oldTopic)
}

func (uc *uploadTopicUseCase) createTopicLanguage(ctx context.Context, req request.UploadTopicRequest, orgID string) (*model.Topic, error) {
	topic := &model.Topic{
		ID:             primitive.NewObjectID(),
		OrganizationID: orgID,
		IsPublished:    req.IsPublished,
		LanguageConfig: []model.TopicLanguageConfig{},
	}

	newTopic, err := uc.topicRepo.CreateTopic(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("create topic fail: %w", err)
	}

	langConfig := model.TopicLanguageConfig{
		LanguageID:  req.LanguageID,
		FileName:    req.FileName,
		Title:       req.Title,
		Note:        req.Note,
		Description: req.Description,
	}
	if err := uc.topicRepo.SetLanguageConfig(ctx, newTopic.ID.Hex(), langConfig); err != nil {
		return nil, fmt.Errorf("set language config fail: %w", err)
	}

	return newTopic, nil
}

func (uc *uploadTopicUseCase) getAudioKeyByLanguage(topic *model.Topic, languageID uint) string {
	for _, lc := range topic.LanguageConfig {
		if lc.LanguageID == languageID {
			if lc.Audio.AudioKey != "" {
				return lc.Audio.AudioKey
			}
			break
		}
	}
	return ""
}
