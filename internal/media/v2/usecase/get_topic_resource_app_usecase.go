package usecase

import (
	"context"
	"media-service/helper"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
	"time"
)

type GetTopicResourceAppUseCase interface {
	GetOutputResources4App(ctx context.Context, studentID string, day, month, year int, topicID string) ([]*response.GetTopicResourcesResponse4App, error)
}

type getTopicResourceAppUseCase struct {
	topicRepo               repository.TopicRepository
	topicResourceRepository repository.TopicResourceRepository
	s3Service               s3.Service
}

func NewGetTopicResourceAppUseCase(topicRepo repository.TopicRepository, topicResourceRepository repository.TopicResourceRepository, s3Service s3.Service) GetTopicResourceAppUseCase {
	return &getTopicResourceAppUseCase{topicRepo: topicRepo, topicResourceRepository: topicResourceRepository, s3Service: s3Service}
}

func (uc *getTopicResourceAppUseCase) GetOutputResources4App(ctx context.Context, studentID string, day, month, year int, topicID string) ([]*response.GetTopicResourcesResponse4App, error) {
	topicResources, err := uc.topicResourceRepository.GetTopicResouresByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	// filter topic resources by month
	if day != 0 && month != 0 && year != 0 {
		topicResources = filterTopicResourcesByDayAndMonthAndYear(topicResources, day, month, year)
	} else if month != 0 && year != 0 {
		topicResources = filterTopicResourcesByMonthAndYear(topicResources, month, year)
	} else if year != 0 {
		topicResources = filterTopicResourcesByYear(topicResources, year)
	}

	if topicID != "" {
		topicResources = filterTopicResourcesByTopicID(topicResources, topicID)
	}

	result := make([]*response.GetTopicResourcesResponse4App, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		resourceImageUrl := ""
		if tr.ImageKey != "" {
			imageUrl, _ := uc.s3Service.Get(ctx, tr.ImageKey, nil)
			if imageUrl != nil {
				resourceImageUrl = *imageUrl
			}
		}
		if tr.IsOutput {
			topic, err := uc.topicRepo.GetByID(ctx, tr.TopicID)
			if err != nil {
				return nil, err
			}
			appLanguage := helper.GetAppLanguage(ctx, 1)
			// Select language config by LanguageID instead of using it as slice index
			var langCfg *model.TopicLanguageConfig
			for i := range topic.LanguageConfig {
				if topic.LanguageConfig[i].LanguageID == appLanguage {
					langCfg = &topic.LanguageConfig[i]
					break
				}
			}
			// fallback to first language if matching language not found
			if langCfg == nil && len(topic.LanguageConfig) > 0 {
				langCfg = &topic.LanguageConfig[0]
			}
			if langCfg != nil {
				for i := range langCfg.Images {
					img := &langCfg.Images[i]
					if img.ImageKey != "" {
						url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
						if err == nil && url != nil {
							img.UploadedUrl = *url
						}
					}
				}
			}

			topicResp := mapper.ToTopicResponse4App(topic, appLanguage)
			result = append(result, mapper.ToGetTopicResourcesResponse4App(ctx, tr, resourceImageUrl, *topicResp))
		}
	}
	return result, nil
}

func filterTopicResourcesByMonthAndYear(topicResources []*model.TopicResource, month, year int) []*model.TopicResource {
	result := make([]*model.TopicResource, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.CreatedAt.Month() == time.Month(month) && tr.CreatedAt.Year() == year {
			result = append(result, tr)
		}
	}
	return result
}

func filterTopicResourcesByTopicID(topicResources []*model.TopicResource, topicID string) []*model.TopicResource {
	result := make([]*model.TopicResource, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.TopicID == topicID {
			result = append(result, tr)
		}
	}
	return result
}

func filterTopicResourcesByDayAndMonthAndYear(topicResources []*model.TopicResource, day, month, year int) []*model.TopicResource {
	result := make([]*model.TopicResource, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.CreatedAt.Day() == day && tr.CreatedAt.Month() == time.Month(month) && tr.CreatedAt.Year() == year {
			result = append(result, tr)
		}
	}
	return result
}

func filterTopicResourcesByYear(topicResources []*model.TopicResource, year int) []*model.TopicResource {
	result := make([]*model.TopicResource, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.CreatedAt.Year() == year {
			result = append(result, tr)
		}
	}
	return result
}
