package usecase

import (
	"context"
	"media-service/helper"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	s3svc "media-service/internal/s3"
)

type GetTopicAppUseCase interface {
	GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error)
	GetTopics4App(ctx context.Context, organizationID string) ([]*response.GetTopic4StudentResponse4App, error)
}

type getTopicAppUseCase struct {
	topicRepo repository.TopicRepository
	s3Service s3svc.Service
}

func NewGetTopicAppUseCase(topicRepo repository.TopicRepository, s3Service s3svc.Service) GetTopicAppUseCase {
	return &getTopicAppUseCase{
		topicRepo: topicRepo,
		s3Service: s3Service,
	}
}

func (uc *getTopicAppUseCase) GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error) {
	topics, err := uc.topicRepo.GetAllTopicsIsPublished(ctx)
	if err != nil {
		return nil, err
	}
	appLanguage := helper.GetAppLanguage(ctx, 1)

	return mapper.ToTopic4StudentResponses4App(topics, appLanguage), nil
}

// Hien tai khong dung den organizationID
func (uc *getTopicAppUseCase) GetTopics4App(ctx context.Context, organizationID string) ([]*response.GetTopic4StudentResponse4App, error) {
	topics, err := uc.topicRepo.GetAllTopicsIsPublished(ctx)
	if err != nil {
		return nil, err
	}
	appLanguage := helper.GetAppLanguage(ctx, 1)

	for ti := range topics {
		// duyệt qua từng language config
		for li := range topics[ti].LanguageConfig {
			langCfg := &topics[ti].LanguageConfig[li]

			// images
			for ii := range langCfg.Images {
				img := &langCfg.Images[ii]
				if img.ImageKey != "" {
					url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
					if err == nil && url != nil {
						img.UploadedUrl = *url
					}
				}
			}
		}
	}

	return mapper.ToTopic4StudentResponses4App(topics, appLanguage), nil
}
