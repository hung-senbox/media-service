package service

import (
	"context"
	"fmt"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/media/dto/request"
	"media-service/internal/media/model"
	"media-service/internal/media/repository"
	"media-service/internal/redis"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicService interface {
	UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error)
	GetUploadProgress(ctx context.Context, topicID string) (int64, error)
}

type topicService struct {
	fileGateway  gateway.FileGateway
	redisService *redis.RedisService
	topicRepo    repository.TopicRepository
}

func NewTopicService(topicRepo repository.TopicRepository, fileGateway gateway.FileGateway, redisService *redis.RedisService) TopicService {
	return &topicService{
		topicRepo:    topicRepo,
		redisService: redisService,
		fileGateway:  fileGateway,
	}
}

func (s *topicService) UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	// --- Step 1: Tạo topic trước ---
	topic := &model.Topic{
		ID:          primitive.NewObjectID(),
		IsPublished: req.IsPublished,
		LanguageConfig: []model.TopicLanguageConfig{
			{
				LanguageID:  req.LanguageID,
				FileName:    req.FileName,
				Title:       req.Title,
				Note:        req.Note,
				Description: req.Description,
				Images:      []model.TopicImageConfig{},
				Videos:      []model.TopicVideoConfig{},
				Audios:      []model.TopicAudioConfig{},
			},
		},
	}

	if err := s.topicRepo.UploadTopic(ctx, topic); err != nil {
		return nil, fmt.Errorf("create topic fail: %w", err)
	}

	// --- Step 2: Tính số task cần upload ---
	totalTasks := 0
	if req.AudioFile != nil {
		totalTasks++
	}
	if req.VideoFile != nil {
		totalTasks++
	}
	if req.FullBackgroundFile != nil {
		totalTasks++
	}
	if req.ClearBackgroundFile != nil {
		totalTasks++
	}
	if req.ClipPartFile != nil {
		totalTasks++
	}
	if req.DrawingFile != nil {
		totalTasks++
	}
	if req.IconFile != nil {
		totalTasks++
	}
	if req.BMFile != nil {
		totalTasks++
	}

	if totalTasks > 0 {
		_ = s.redisService.InitUploadProgress(ctx, topic.ID.Hex(), totalTasks)
	}

	// --- Step 3: Upload async ---
	go s.uploadFilesAsync(topic.ID.Hex(), req)

	return topic, nil
}

func (s *topicService) uploadFilesAsync(topicID string, req request.UploadTopicRequest) {
	ctx := context.Background()

	defer func() {
		if r := recover(); r != nil {
			// log error nếu cần
		}
	}()

	// --- audio ---
	if req.AudioFile != nil {
		s.uploadAndSaveAudio(ctx, topicID, req)
		_, _ = s.redisService.DecrementUploadTask(ctx, topicID)
	}

	// --- video ---
	if req.VideoFile != nil {
		s.uploadAndSaveVideo(ctx, topicID, req)
		_, _ = s.redisService.DecrementUploadTask(ctx, topicID)
	}

	// --- images ---
	numImages := s.uploadAndSaveImages(ctx, topicID, req)
	for i := 0; i < numImages; i++ {
		_, _ = s.redisService.DecrementUploadTask(ctx, topicID)
	}
}

// --- Upload & Save Audio ---
func (s *topicService) uploadAndSaveAudio(ctx context.Context, topicID string, req request.UploadTopicRequest) {
	if req.AudioFile == nil {
		return
	}

	uploadReq := gw_request.UploadFileRequest{
		File:     req.AudioFile,
		Folder:   "topic_audios",
		FileName: req.AudioFile.Filename,
		Mode:     "private",
	}
	audioResp, err := s.fileGateway.UploadAudio(ctx, uploadReq)
	if err != nil {
		//log.Errorf("upload audio fail: %v", err)
		return
	}

	_ = s.topicRepo.AddAudioToTopic(ctx, topicID, model.TopicAudioConfig{
		AudioKey:  audioResp.Key,
		LinkUrl:   req.AudioLinkUrl,
		StartTime: req.AudioStart,
		EndTime:   req.AudioEnd,
	})
}

// --- Upload & Save Video ---
func (s *topicService) uploadAndSaveVideo(ctx context.Context, topicID string, req request.UploadTopicRequest) {
	if req.VideoFile == nil {
		return
	}

	uploadReq := gw_request.UploadFileRequest{
		File:     req.VideoFile,
		Folder:   "topic_videos",
		FileName: req.VideoFile.Filename,
		Mode:     "private",
	}
	videoResp, err := s.fileGateway.UploadVideo(ctx, uploadReq)
	if err != nil {
		//log.Errorf("upload video fail: %v", err)
		return
	}

	_ = s.topicRepo.AddVideoToTopic(ctx, topicID, model.TopicVideoConfig{
		VideoKey:  videoResp.Key,
		LinkUrl:   req.VideoLinkUrl,
		StartTime: req.VideoStart,
		EndTime:   req.VideoEnd,
	})
}

// --- Upload & Save Images ---
// return số lượng ảnh được upload thành công
func (s *topicService) uploadAndSaveImages(ctx context.Context, topicID string, req request.UploadTopicRequest) int {
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

	count := 0
	for _, img := range imageFiles {
		if img.file == nil {
			continue
		}

		uploadReq := gw_request.UploadFileRequest{
			File:      img.file,
			Folder:    "topic_images",
			FileName:  img.file.Filename,
			ImageName: img.typ,
			Mode:      "private",
		}
		imgResp, err := s.fileGateway.UploadImage(ctx, uploadReq)
		if err != nil {
			//log.Errorf("upload image (%s) fail: %v", img.typ, err)
			continue
		}

		_ = s.topicRepo.AddImageToTopic(ctx, topicID, model.TopicImageConfig{
			ImageKey:  imgResp.Key,
			ImageType: img.typ,
			LinkUrl:   img.link,
		})
		count++
	}
	return count
}

func (s *topicService) GetUploadProgress(ctx context.Context, topicID string) (int64, error) {
	return s.redisService.GetUploadProgress(ctx, topicID)
}
