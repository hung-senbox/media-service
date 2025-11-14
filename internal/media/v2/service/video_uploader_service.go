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
	"sort"
	"strings"
	"time"
)

type VideoUploaderService interface {
	UploadVideoUploader(ctx context.Context, req request.UploadVideoUploaderRequest) error
	GetUploaderStatus(ctx context.Context, videoUploaderID string) (response.GetUploaderStatusResponse, error)
	GetVideosUploader4Web(ctx context.Context, languageID string, title string, sortBy []request.GetVideoUploaderSortBy) ([]response.GetVideoUploaderResponse4Web, error)
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
func (s *videoUploaderService) UploadVideoUploader(ctx context.Context, req request.UploadVideoUploaderRequest) error {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if !currentUser.IsSuperAdmin {
		return fmt.Errorf("access denied")
	}

	var videoUploader *model.VideoUploader

	// Step 1: tạo record trong MongoDB (insert hoặc update)
	if req.VideoUploaderID != "" {
		existing, err := s.videoUploaderRepository.GetVideoUploaderByID(ctx, req.VideoUploaderID)
		if err != nil {
			return fmt.Errorf("get video uploader failed: %w", err)
		}
		existing.IsVisible = req.IsVisible
		existing.Title = req.Title
		existing.LanguageID = req.LanguageID
		existing.Note = req.Note
		existing.Transcript = req.Transcript
		existing.UpdatedAt = time.Now()
		if req.IsDeletedVideo {
			s.videoUploaderRepository.DeleteVideoMetadata(ctx, req.VideoUploaderID)
		}
		if req.IsDeletedImagePreview {
			s.videoUploaderRepository.DeleteImagePreviewMetadata(ctx, req.VideoUploaderID)
		}
		if err := s.videoUploaderRepository.SetVideoUploaderWithoutFiles(ctx, existing); err != nil {
			return fmt.Errorf("save video uploader failed: %w", err)
		}
		videoUploader = existing
	} else {
		newVideo := &model.VideoUploader{
			CreatedBy:  currentUser.ID,
			IsVisible:  req.IsVisible,
			LanguageID: req.LanguageID,
			Title:      req.Title,
			Note:       req.Note,
			Transcript: req.Transcript,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		if err := s.videoUploaderRepository.SetVideoUploaderWithoutFiles(ctx, newVideo); err != nil {
			return fmt.Errorf("save video uploader failed: %w", err)
		}
		videoUploader = newVideo
	}

	// Step 2: Upload đồng bộ, không dùng Redis
	// Upload video nếu có
	if helper.IsValidFile(req.VideoFile) {
		oldVideoKey := videoUploader.VideoKey
		if oldVideoKey != "" {
			_ = s.s3Service.Delete(ctx, oldVideoKey)
		}
		if err := s.processVideoUpload(ctx, videoUploader.ID.Hex(), req); err != nil {
			return fmt.Errorf("video upload failed: %w", err)
		}
	}
	// Upload ảnh preview nếu có
	if helper.IsValidFile(req.ImagePreviewFile) {
		oldImageKey := videoUploader.ImagePreviewKey
		if oldImageKey != "" {
			_ = s.s3Service.Delete(ctx, oldImageKey)
		}
		if err := s.processImagePreviewUpload(ctx, videoUploader.ID.Hex(), req); err != nil {
			return fmt.Errorf("image upload failed: %w", err)
		}
	}

	return nil
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
	if err := s.videoUploaderRepository.SetVideoMetadata(ctx, videoUploaderID, key, deref(url)); err != nil {
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
	if err := s.videoUploaderRepository.SetImagePreviewMetadata(ctx, videoUploaderID, key, deref(url)); err != nil {
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

	return response.GetUploaderStatusResponse{
		Status:    "done",
		Message:   "upload completed",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func (s *videoUploaderService) GetVideosUploader4Web(ctx context.Context, languageID, title string, sortBy []request.GetVideoUploaderSortBy) ([]response.GetVideoUploaderResponse4Web, error) {
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
			videoUploaders = filterVideosByTitle(videoUploaders, strings.TrimSpace(title))
			videoUploaders = sortVideos(videoUploaders, sortBy)
			return mapper.ToGetVideosResponse4Web(videoUploaders, currentUser.Nickname), nil
		}
	}

	// neu khong co languageID hoac languageID = 0 thi lay tat ca
	videoUploaders, err = s.videoUploaderRepository.GetAllVideos(ctx)
	if err != nil {
		return nil, err
	}
	videoUploaders = filterVideosByTitle(videoUploaders, strings.TrimSpace(title))
	videoUploaders = sortVideos(videoUploaders, sortBy)
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

	if videoUploader.VideoKey != "" {
		err = s.s3Service.Delete(ctx, videoUploader.VideoKey)
		if err != nil {
			return err
		}
	}
	if videoUploader.ImagePreviewKey != "" {
		err = s.s3Service.Delete(ctx, videoUploader.ImagePreviewKey)
		if err != nil {
			return err
		}
	}

	return s.videoUploaderRepository.DeleteVideoUploader(ctx, videoUploaderID)
}

func filterVideosByTitle(videoUploaders []model.VideoUploader, title string) []model.VideoUploader {
	var result []model.VideoUploader
	for _, videoUploader := range videoUploaders {
		if title == "" || strings.Contains(strings.ToLower(videoUploader.Title), strings.ToLower(title)) {
			result = append(result, videoUploader)
		}
	}
	return result
}

func sortVideos(vs []model.VideoUploader, sortBy []request.GetVideoUploaderSortBy) []model.VideoUploader {
	if len(sortBy) == 0 {
		return vs
	}

	sort.Slice(vs, func(i, j int) bool {
		for _, sb := range sortBy {
			field := strings.ToLower(sb.Field)
			asc := strings.ToLower(sb.Order) != request.GetVideoUploaderSortByOrderDesc

			switch field {
			case "title":
				if vs[i].Title != vs[j].Title {
					if asc {
						return strings.ToLower(vs[i].Title) < strings.ToLower(vs[j].Title)
					}
					return strings.ToLower(vs[i].Title) > strings.ToLower(vs[j].Title)
				}
			case "language_id":
				if vs[i].LanguageID != vs[j].LanguageID {
					if asc {
						return vs[i].LanguageID < vs[j].LanguageID
					}
					return vs[i].LanguageID > vs[j].LanguageID
				}
			case "updated_at":
				if !vs[i].UpdatedAt.Equal(vs[j].UpdatedAt) {
					if asc {
						return vs[i].UpdatedAt.Before(vs[j].UpdatedAt)
					}
					return vs[i].UpdatedAt.After(vs[j].UpdatedAt)
				}
			case "created_at":
				if !vs[i].CreatedAt.Equal(vs[j].CreatedAt) {
					if asc {
						return vs[i].CreatedAt.Before(vs[j].CreatedAt)
					}
					return vs[i].CreatedAt.After(vs[j].CreatedAt)
				}
			}
			// nếu bằng nhau → chuyển sang sort field tiếp theo
		}
		return false
	})

	return vs
}
