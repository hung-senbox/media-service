package service

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/redis"
	"media-service/internal/s3"
	"media-service/pkg/constants"
	"media-service/pkg/uploader"
	"time"
)

type VideoUploaderService interface {
	UploadVideoUploader(ctx context.Context, req request.UploadVideoUploaderRequest) (*model.VideoUploader, error)
	GetUploaderStatus(ctx context.Context, videoUploaderID string) (response.GetUploaderStatusResponse, error)
	GetVideosUploader4Web(ctx context.Context) ([]response.GetVideoUploaderResponse4Web, error)
}

type videoUploaderService struct {
	videoUploaderRepository repository.VideoUploaderRepository
	s3Service               s3.Service
	redisService            *redis.RedisService
	userGateway             gateway.UserGateway
}

func NewVideoUploaderService(videoUploaderRepository repository.VideoUploaderRepository, s3Service s3.Service, redisService *redis.RedisService, userGateway gateway.UserGateway) VideoUploaderService {
	return &videoUploaderService{videoUploaderRepository: videoUploaderRepository, s3Service: s3Service, redisService: redisService, userGateway: userGateway}
}

// ======================================================
// =============== PUBLIC FUNCTION ======================
// ======================================================
func (s *videoUploaderService) UploadVideoUploader(ctx context.Context, req request.UploadVideoUploaderRequest) (*model.VideoUploader, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if !currentUser.IsSuperAdmin {
		return nil, fmt.Errorf("access denied")
	}

	var videoUploader *model.VideoUploader

	// Step 1: tạo record trong MongoDB (insert hoặc update)
	if req.VideoUploaderID != "" {
		existing, err := s.videoUploaderRepository.GetVideoUploaderByID(ctx, req.VideoUploaderID)
		if err != nil {
			return nil, fmt.Errorf("get video uploader failed: %w", err)
		}
		existing.IsVisible = req.IsVisible
		existing.UpdatedAt = time.Now()

		if err := s.videoUploaderRepository.SetVideoUploaderWithoutFiles(ctx, existing); err != nil {
			return nil, fmt.Errorf("save video uploader failed: %w", err)
		}
		videoUploader = existing
	} else {
		newVideo := &model.VideoUploader{
			CreatedBy: currentUser.ID,
			IsVisible: req.IsVisible,
			LanguageConfig: []model.VideoLanguageConfig{
				{
					LanguageID: req.LanguageID,
					Title:      req.Title,
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.videoUploaderRepository.SetVideoUploaderWithoutFiles(ctx, newVideo); err != nil {
			return nil, fmt.Errorf("save video uploader failed: %w", err)
		}
		videoUploader = newVideo
	}

	// Step 2: tạo Redis key mới (có orgID)
	key := helper.BuildVideoUploaderRedisKey(videoUploader.ID.Hex())

	// Ghi trạng thái pending
	status := map[string]interface{}{
		"status":     "pending",
		"message":    "upload started",
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	}
	if err := s.redisService.SetUploaderStatus(ctx, key, status); err != nil {
		fmt.Printf("[UploadVideoUploader] failed to set redis pending status: %v\n", err)
	}

	// Step 3: chạy upload async (fire-and-forget)
	ctxUpload, cancel := context.WithTimeout(context.Background(), 15*time.Minute)

	if token, ok := ctx.Value(constants.Token).(string); ok {
		ctxUpload = context.WithValue(ctxUpload, constants.Token, token)
	}

	go func() {
		defer cancel()
		s.asyncUploadProcess(ctxUpload, videoUploader, req, key)
	}()

	// Step 4: trả về response ngay lập tức
	return videoUploader, nil
}

func (s *videoUploaderService) asyncUploadProcess(ctx context.Context, videoUploader *model.VideoUploader, req request.UploadVideoUploaderRequest, redisKey string) {
	defer func() {
		if r := recover(); r != nil {
			s.redisService.SetUploaderStatus(ctx, redisKey, map[string]interface{}{
				"status":     "failed",
				"message":    fmt.Sprintf("panic: %v", r),
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			})
		}
	}()

	// upload video
	if helper.IsValidFile(req.VideoFile) {
		// neu co file thi check neu co file cu thi xoa
		oldVideoKey := ""
		for _, lc := range videoUploader.LanguageConfig {
			if lc.LanguageID == req.LanguageID {
				oldVideoKey = lc.VideoKey
				break
			}
		}
		if oldVideoKey != "" {
			_ = s.s3Service.Delete(ctx, oldVideoKey)
		}
		// upload moi
		err := s.processVideoUpload(ctx, videoUploader.ID.Hex(), req)
		if err != nil {
			s.redisService.SetUploaderStatus(ctx, redisKey, map[string]interface{}{
				"status":     "failed",
				"message":    fmt.Sprintf("video upload failed: %v", err),
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}
	}

	// upload ảnh preview
	if helper.IsValidFile(req.ImagePreviewFile) {
		// neu co file thi check neu co file cu thi xoa
		oldImageKey := ""
		for _, lc := range videoUploader.LanguageConfig {
			if lc.LanguageID == req.LanguageID {
				oldImageKey = lc.ImagePreviewKey
				break
			}
		}
		if oldImageKey != "" {
			_ = s.s3Service.Delete(ctx, oldImageKey)
		}
		// upload moi
		err := s.processImagePreviewUpload(ctx, videoUploader.ID.Hex(), req)
		if err != nil {
			s.redisService.SetUploaderStatus(ctx, redisKey, map[string]interface{}{
				"status":     "failed",
				"message":    fmt.Sprintf("image upload failed: %v", err),
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			})
		}
	}

	// thành công
	s.redisService.SetUploaderStatus(ctx, redisKey, map[string]interface{}{
		"status":     "success",
		"message":    "upload completed successfully",
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	})
}

// ======================================================
// =============== PRIVATE HELPERS ======================
// ======================================================

// xử lý upload video và lưu metadata
func (s *videoUploaderService) processVideoUpload(ctx context.Context, videoUploaderID string, req request.UploadVideoUploaderRequest) error {
	if req.VideoFile == nil {
		return fmt.Errorf("video file is required")
	}

	key := helper.BuildObjectKeyS3("media_video_uploader", req.VideoFile.Filename, "video_"+req.Title)
	f, openErr := req.VideoFile.Open()
	if openErr != nil {
		return openErr
	}
	defer f.Close()
	ct := req.VideoFile.Header.Get("Content-Type")
	url, err := s.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPublic)
	if err != nil {
		return err
	}

	// cập nhật metadata vào DB
	lang := helper.GetAppLanguage(ctx, 1)
	if err := s.videoUploaderRepository.SetVideoMetadata(ctx, videoUploaderID, lang, key, deref(url)); err != nil {
		return fmt.Errorf("save video metadata failed: %w", err)
	}

	return nil
}

// xử lý upload ảnh preview và lưu metadata
func (s *videoUploaderService) processImagePreviewUpload(ctx context.Context, videoUploaderID string, req request.UploadVideoUploaderRequest) error {
	if req.ImagePreviewFile == nil {
		return fmt.Errorf("image preview file is required")
	}

	key := helper.BuildObjectKeyS3("media_video_uploader", req.ImagePreviewFile.Filename, "image_preview_"+req.Title)
	f, openErr := req.ImagePreviewFile.Open()
	if openErr != nil {
		return openErr
	}
	defer f.Close()
	ct := req.ImagePreviewFile.Header.Get("Content-Type")
	url, err := s.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPublic)
	if err != nil {
		return err
	}

	// cập nhật metadata vào DB
	lang := helper.GetAppLanguage(ctx, 1)
	if err := s.videoUploaderRepository.SetImagePreviewMetadata(ctx, videoUploaderID, lang, key, deref(url)); err != nil {
		return fmt.Errorf("save image metadata failed: %w", err)
	}

	return nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (s *videoUploaderService) GetUploaderStatus(ctx context.Context, videoUploaderID string) (response.GetUploaderStatusResponse, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if !currentUser.IsSuperAdmin {
		return response.GetUploaderStatusResponse{}, fmt.Errorf("access denied")
	}

	key := helper.BuildVideoUploaderRedisKey(videoUploaderID)
	status, err := s.redisService.GetUploaderStatus(ctx, key)
	if err != nil {
		return response.GetUploaderStatusResponse{}, err
	}
	if len(status) == 0 {
		return response.GetUploaderStatusResponse{}, nil
	}

	s.redisService.DeleteUploaderStatusKey(ctx, key)

	statusVal, _ := status["status"].(string)
	messageVal, _ := status["message"].(string)
	updatedAtVal, _ := status["updated_at"].(string)

	return response.GetUploaderStatusResponse{
		Status:    statusVal,
		Message:   messageVal,
		UpdatedAt: updatedAtVal,
	}, nil
}

func (s *videoUploaderService) GetVideosUploader4Web(ctx context.Context) ([]response.GetVideoUploaderResponse4Web, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser.IsSuperAdmin {
		videoUploaders, err := s.videoUploaderRepository.GetAllVideos(ctx)
		if err != nil {
			return nil, err
		}
		return mapper.ToGetVideosResponse4Web(videoUploaders, currentUser.Nickname), nil
	}

	videoUploaders, err := s.videoUploaderRepository.GetVideosIsVisible(ctx)
	if err != nil {
		return nil, err
	}

	user, _ := s.userGateway.GetUserByID(ctx, videoUploaders[0].CreatedBy)
	createdByName := ""
	if user != nil {
		createdByName = user.Name
	}
	return mapper.ToGetVideosResponse4Web(videoUploaders, createdByName), nil
}
