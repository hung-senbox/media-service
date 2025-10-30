package usecase

import (
	"context"
	"fmt"
	"math"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/repository"
	"media-service/internal/redis"
	"media-service/pkg/constants"
)

type GetUploadProgressUseCase interface {
	GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error)
}

type getUploadProgressUseCase struct {
	topicRepo    repository.TopicRepository
	redisService *redis.RedisService
}

func NewGetUploadProgressUseCase(redisService *redis.RedisService) GetUploadProgressUseCase {
	return &getUploadProgressUseCase{
		redisService: redisService,
	}
}

func (uc *getUploadProgressUseCase) GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}

	orgID := currentUser.OrganizationAdmin.ID
	total, err := uc.redisService.GetTotalUploadTask(ctx, orgID, topicID)
	if err != nil {
		return nil, err
	}
	remaining, err := uc.redisService.GetUploadProgress(ctx, orgID, topicID)
	if err != nil {
		return nil, err
	}
	rawErrors, err := uc.redisService.GetUploadErrors(ctx, orgID, topicID)
	if err != nil {
		return nil, err
	}
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	// Nếu chưa từng tạo task upload nào => chưa upload gì cả (hoac case da dat 100 progress thi da xoa het cache)
	if total == 0 {
		// goi delete cache
		_ = uc.redisService.DeleteUploadProgress(ctx, orgID, topicID)

		return &response.GetUploadProgressResponse{
			Progress: -1,
			FileName: topic.LanguageConfig[0].FileName,
			UploadErrors: map[string]any{
				"audio_error": "",
				"video_error": "",
				"image_error": map[string]string{
					"full_background":  "",
					"clear_background": "",
					"clip_part":        "",
					"drawing":          "",
					"icon":             "",
					"bm":               "",
				},
			},
		}, nil
	}

	// Nếu có task upload
	progress := int((total - remaining) * 100 / total)
	progress = int(math.Min(float64(progress), 100.0))

	imageErr := map[string]string{
		"full_background":  rawErrors["image_full_background"],
		"clear_background": rawErrors["image_clear_background"],
		"clip_part":        rawErrors["image_clip_part"],
		"drawing":          rawErrors["image_drawing"],
		"icon":             rawErrors["image_icon"],
		"bm":               rawErrors["image_bm"],
	}

	// goi delete cache
	_ = uc.redisService.DeleteUploadProgress(ctx, orgID, topicID)

	return &response.GetUploadProgressResponse{
		Progress: progress,
		FileName: topic.LanguageConfig[0].FileName,
		UploadErrors: map[string]any{
			"audio_error": rawErrors["audio_error"],
			"video_error": rawErrors["video_error"],
			"image_error": imageErr,
		},
	}, nil
}
