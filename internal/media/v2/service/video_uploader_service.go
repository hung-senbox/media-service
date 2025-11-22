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

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoUploaderService interface {
	UploadVideoUploader(ctx context.Context, req request.UploadVideoUploaderRequest) (*model.VideoUploader, error)
	GetVideosUploader4Web(ctx context.Context, languageID string, title string, sortBy []request.GetVideoUploaderSortBy) ([]response.GetVideoUploaderResponse4Web, error)
	DeleteVideoUploader(ctx context.Context, videoUploaderID string) error
	GetVideo4Web(ctx context.Context, videoUploaderID string) (*response.GetDetailVideo4WebResponse, error)
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

	// Step 1: tạo / lấy record trong MongoDB (insert hoặc update, chưa xử lý file)
	if req.VideoFolderID != "" {
		existing, err := s.videoUploaderRepository.GetVideoUploaderByID(ctx, req.VideoFolderID)
		if err != nil {
			return nil, fmt.Errorf("get video uploader failed: %w", err)
		}
		existing.IsVisible = req.IsVisible
		existing.Title = req.Title
		existing.UpdatedAt = time.Now()
		videoUploader = existing
	} else {
		newVideo := &model.VideoUploader{
			CreatedBy:      currentUser.ID,
			IsVisible:      req.IsVisible,
			Title:          req.Title,
			LanguageConfig: make([]model.VideoUploaderLanguageConfig, 0),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		videoUploader = newVideo
	}

	// Step 2: xử lý language config tương ứng với LanguageID
	langCfgIdx := -1
	for i := range videoUploader.LanguageConfig {
		if videoUploader.LanguageConfig[i].LanguageID == req.LanguageID {
			langCfgIdx = i
			break
		}
	}
	if langCfgIdx == -1 {
		videoUploader.LanguageConfig = append(videoUploader.LanguageConfig, model.VideoUploaderLanguageConfig{
			ID:         primitive.NewObjectID(),
			LanguageID: req.LanguageID,
		})
		langCfgIdx = len(videoUploader.LanguageConfig) - 1
	}
	cfg := &videoUploader.LanguageConfig[langCfgIdx]

	// Apply note & transcript (giá trị mới nhất từ request)
	cfg.Note = req.Note
	cfg.Transcript = req.Transcript

	// Xử lý xoá trước khi upload mới
	if req.IsDeletedVideo {
		if cfg.VideoKey != "" {
			_ = s.s3Service.Delete(ctx, cfg.VideoKey)
		}
		cfg.VideoKey = ""
		cfg.VideoPublicUrl = ""
	}
	if req.IsDeletedImagePreview {
		if cfg.ImagePreviewKey != "" {
			_ = s.s3Service.Delete(ctx, cfg.ImagePreviewKey)
		}
		cfg.ImagePreviewKey = ""
		cfg.ImagePreviewPublicUrl = ""
	}

	// Step 3: Upload đồng bộ video & image cho language config này
	// Upload video nếu có
	if helper.IsValidFile(req.VideoFile) {
		if cfg.VideoKey != "" {
			_ = s.s3Service.Delete(ctx, cfg.VideoKey)
		}
		videoKey, videoUrl, err := s.processVideoUpload(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("video upload failed: %w", err)
		}
		cfg.VideoKey = videoKey
		cfg.VideoPublicUrl = videoUrl
	}
	// Upload ảnh preview nếu có
	if helper.IsValidFile(req.ImagePreviewFile) {
		if cfg.ImagePreviewKey != "" {
			_ = s.s3Service.Delete(ctx, cfg.ImagePreviewKey)
		}
		imageKey, imageUrl, err := s.processImagePreviewUpload(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("image upload failed: %w", err)
		}
		cfg.ImagePreviewKey = imageKey
		cfg.ImagePreviewPublicUrl = imageUrl
	}

	// Step 4: Lưu toàn bộ document (bao gồm language_config) vào MongoDB
	if err := s.videoUploaderRepository.SetVideoUploader(ctx, videoUploader); err != nil {
		return nil, fmt.Errorf("save video uploader failed: %w", err)
	}

	return videoUploader, nil
}

// ======================================================
// =============== PRIVATE HELPERS ======================
// ======================================================

// xử lý upload video và trả về key + public URL
func (s *videoUploaderService) processVideoUpload(ctx context.Context, req request.UploadVideoUploaderRequest) (string, string, error) {
	if req.VideoFile == nil {
		return "", "", fmt.Errorf("video file is required")
	}

	key := helper.BuildObjectKeyS3("media_video_uploader", req.VideoFile.Filename, "video_"+req.Title)
	f, openErr := req.VideoFile.Open()
	if openErr != nil {
		return "", "", openErr
	}
	defer f.Close()
	ct := req.VideoFile.Header.Get("Content-Type")
	url, err := s.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPublic)
	if err != nil {
		return "", "", err
	}
	return key, deref(url), nil
}

// xử lý upload ảnh preview và trả về key + public URL
func (s *videoUploaderService) processImagePreviewUpload(ctx context.Context, req request.UploadVideoUploaderRequest) (string, string, error) {
	if req.ImagePreviewFile == nil {
		return "", "", fmt.Errorf("image preview file is required")
	}

	key := helper.BuildObjectKeyS3("media_video_uploader", req.ImagePreviewFile.Filename, "image_preview_"+req.Title)
	f, openErr := req.ImagePreviewFile.Open()
	if openErr != nil {
		return "", "", openErr
	}
	defer f.Close()
	ct := req.ImagePreviewFile.Header.Get("Content-Type")
	url, err := s.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPublic)
	if err != nil {
		return "", "", err
	}
	return key, deref(url), nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
			return mapper.ToGetVideosResponse4Web(videoUploaders, currentUser.Nickname, langID), nil
		}
	}

	// neu khong co languageID hoac languageID = 0 thi lay tat ca
	videoUploaders, err = s.videoUploaderRepository.GetAllVideos(ctx)
	if err != nil {
		return nil, err
	}
	videoUploaders = filterVideosByTitle(videoUploaders, strings.TrimSpace(title))
	videoUploaders = sortVideos(videoUploaders, sortBy)
	return mapper.ToGetVideosResponse4Web(videoUploaders, currentUser.Nickname, 0), nil

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

	// Xoá toàn bộ file video & image preview của tất cả language config
	for _, cfg := range videoUploader.LanguageConfig {
		if cfg.VideoKey != "" {
			if err := s.s3Service.Delete(ctx, cfg.VideoKey); err != nil {
				return err
			}
		}
		if cfg.ImagePreviewKey != "" {
			if err := s.s3Service.Delete(ctx, cfg.ImagePreviewKey); err != nil {
				return err
			}
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
				// sort theo language_id đầu tiên trong language_config (nếu có)
				var li, lj uint
				if len(vs[i].LanguageConfig) > 0 {
					li = vs[i].LanguageConfig[0].LanguageID
				}
				if len(vs[j].LanguageConfig) > 0 {
					lj = vs[j].LanguageConfig[0].LanguageID
				}
				if li != lj {
					if asc {
						return li < lj
					}
					return li > lj
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

func (s *videoUploaderService) GetVideo4Web(ctx context.Context, videoUploaderID string) (*response.GetDetailVideo4WebResponse, error) {

	videoUploader, err := s.videoUploaderRepository.GetVideoUploaderByID(ctx, videoUploaderID)
	if err != nil {
		return nil, err
	}
	return mapper.ToGetDetailVideo4WebResponse(videoUploader), nil
}
