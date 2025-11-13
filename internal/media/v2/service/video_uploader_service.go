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
	"media-service/internal/s3"
	"media-service/pkg/constants"
	"media-service/pkg/uploader"
	"strings"
	"time"
)

type VideoUploaderService interface {
	UploadVideoUploader(ctx context.Context, req request.UploadVideoUploaderRequest) (*model.VideoUploader, error)
	GetVideosUploader4Web(ctx context.Context, languageID string) ([]response.GetVideoUploaderResponse4Web, error)
	DeleteVideoUploader(ctx context.Context, videoUploaderID string) error
}

type videoUploaderService struct {
	videoUploaderRepository repository.VideoUploaderRepository
	s3Service               s3.Service
	userGateway             gateway.UserGateway
}

func NewVideoUploaderService(videoUploaderRepository repository.VideoUploaderRepository, s3Service s3.Service, userGateway gateway.UserGateway) VideoUploaderService {
	return &videoUploaderService{videoUploaderRepository: videoUploaderRepository, s3Service: s3Service, userGateway: userGateway}
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
		if req.IsDeletedVideo {
			s.videoUploaderRepository.DeleteVideoMetadata(ctx, req.VideoUploaderID)
		}
		if req.IsDeletedImagePreview {
			s.videoUploaderRepository.DeleteImagePreviewMetadata(ctx, req.VideoUploaderID)
		}
		if err := s.videoUploaderRepository.SetVideoUploaderWithoutFiles(ctx, existing); err != nil {
			return nil, fmt.Errorf("save video uploader failed: %w", err)
		}
		videoUploader = existing
	} else {
		newVideo := &model.VideoUploader{
			CreatedBy: currentUser.ID,
			IsVisible: req.IsVisible,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.videoUploaderRepository.SetVideoUploaderWithoutFiles(ctx, newVideo); err != nil {
			return nil, fmt.Errorf("save video uploader failed: %w", err)
		}
		videoUploader = newVideo
	}

	// Step 2: xử lý upload đồng bộ
	if helper.IsValidFile(req.VideoFile) {
		if err := s.processVideoUpload(ctx, videoUploader.ID.Hex(), req); err != nil {
			return nil, fmt.Errorf("video upload failed: %w", err)
		}
	}
	if helper.IsValidFile(req.ImagePreviewFile) {
		if err := s.processImagePreviewUpload(ctx, videoUploader.ID.Hex(), req); err != nil {
			return nil, fmt.Errorf("image upload failed: %w", err)
		}
	}

	// Step 3: trả về sau khi hoàn tất
	return videoUploader, nil
}

// ======================================================
// =============== PRIVATE HELPERS ======================
// ======================================================

// xử lý upload video và lưu metadata
func (s *videoUploaderService) processVideoUpload(ctx context.Context, videoUploaderID string, req request.UploadVideoUploaderRequest) error {

	if helper.IsValidFile(req.VideoFile) {
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
		if err := s.videoUploaderRepository.SetVideoMetadata(ctx, videoUploaderID, key, deref(url)); err != nil {
			return fmt.Errorf("save video metadata failed: %w", err)
		}
	}

	return nil
}

// xử lý upload ảnh preview và lưu metadata
func (s *videoUploaderService) processImagePreviewUpload(ctx context.Context, videoUploaderID string, req request.UploadVideoUploaderRequest) error {

	if helper.IsValidFile(req.ImagePreviewFile) {
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
		if err := s.videoUploaderRepository.SetImagePreviewMetadata(ctx, videoUploaderID, key, deref(url)); err != nil {
			return fmt.Errorf("save image metadata failed: %w", err)
		}
	}

	return nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// removed GetUploaderStatus (no Redis tracking)

func (s *videoUploaderService) GetVideosUploader4Web(ctx context.Context, languageID string) ([]response.GetVideoUploaderResponse4Web, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if !currentUser.IsSuperAdmin {
		return nil, fmt.Errorf("access denied")
	}

	// Default: fetch all
	var (
		videoUploaders []model.VideoUploader
		err            error
	)

	// case co languageID thi lay theo languageID
	if langStr := strings.TrimSpace(languageID); langStr != "" {
		var langID uint
		if _, scanErr := fmt.Sscan(langStr, &langID); scanErr == nil && langID > 0 {
			videoUploaders, err = s.videoUploaderRepository.GetVideosByLanguageID(ctx, langID)
			if err != nil {
				return nil, err
			}
			return mapper.ToGetVideosResponse4Web(videoUploaders, currentUser.Nickname), nil
		}
	}

	// neu khong co languageID hoac languageID = 0 thi lay tat ca
	videoUploaders, err = s.videoUploaderRepository.GetAllVideos(ctx)
	if err != nil {
		return nil, err
	}
	return mapper.ToGetVideosResponse4Web(videoUploaders, currentUser.Nickname), nil

}

func (s *videoUploaderService) DeleteVideoUploader(ctx context.Context, videoUploaderID string) error {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if !currentUser.IsSuperAdmin {
		return fmt.Errorf("access denied")
	}

	videoUploader, err := s.videoUploaderRepository.GetVideoUploaderByID(ctx, videoUploaderID)
	if err != nil {
		return err
	}
	if videoUploader == nil {
		return fmt.Errorf("video uploader not found")
	}

	return s.videoUploaderRepository.DeleteVideoUploader(ctx, videoUploaderID)
}
