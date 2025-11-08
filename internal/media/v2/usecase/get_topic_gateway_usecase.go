package usecase

import (
	"context"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	s3svc "media-service/internal/s3"
)

type GetTopicGatewayUseCase interface {
	GetTopics4Student4Gw(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Gw, error)
	GetTopic4Gw(ctx context.Context, topicID string) (*response.TopicResponse4GW, error)
	GetAllTopicsByOrganization4Gw(ctx context.Context, organizationID string) ([]*response.TopicResponse4GW, error)
}

type getTopicGatewayUseCase struct {
	topicRepo    repository.TopicRepository
	cachedUserGw gateway.UserGateway
	s3Service    s3svc.Service
}

func NewGetTopicGatewayUseCase(topicRepo repository.TopicRepository, cachedUserGw gateway.UserGateway, s3Service s3svc.Service) GetTopicGatewayUseCase {
	return &getTopicGatewayUseCase{
		topicRepo:    topicRepo,
		cachedUserGw: cachedUserGw,
		s3Service:    s3Service,
	}
}

func (uc *getTopicGatewayUseCase) GetTopics4Student4Gw(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Gw, error) {
	// get org by student
	// student, err := uc.cachedUserGw.GetStudentInfo(ctx, studentID)
	// if err != nil {
	// 	return nil, err
	// }
	topics, err := uc.topicRepo.GetAllTopicsIsPublished(ctx)
	if err != nil {
		return nil, err
	}
	for ti := range topics {
		for li := range topics[ti].LanguageConfig {
			langCfg := &topics[ti].LanguageConfig[li]
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
	return mapper.ToTopic4StudentResponses4Gw(topics, 1), nil
}
func (uc *getTopicGatewayUseCase) GetTopic4Gw(ctx context.Context, topicID string) (*response.TopicResponse4GW, error) {

	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	appLang := helper.GetAppLanguage(ctx, 1)

	// images
	for ii := range topic.LanguageConfig[0].Images {
		img := &topic.LanguageConfig[0].Images[ii]
		if img.ImageKey != "" {
			url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
			if err == nil && url != nil {
				img.UploadedUrl = *url
			}
		}
	}

	return mapper.ToTopicResponses4GW(topic, appLang), nil
}

// Hien tai khong dung den organizationID
func (uc *getTopicGatewayUseCase) GetAllTopicsByOrganization4Gw(ctx context.Context, organizationID string) ([]*response.TopicResponse4GW, error) {

	topics, err := uc.topicRepo.GetAllTopicsIsPublished(ctx)
	if err != nil {
		return nil, err
	}

	appLang := helper.GetAppLanguage(ctx, 1)

	var result []*response.TopicResponse4GW

	for _, topic := range topics {
		// xử lý ảnh từng topic
		for ii := range topic.LanguageConfig[0].Images {
			img := &topic.LanguageConfig[0].Images[ii]
			if img.ImageKey != "" {
				url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
				if err == nil && url != nil {
					img.UploadedUrl = *url
				}
			}
		}

		topicRes := mapper.ToTopicResponses4GW(&topic, appLang)
		if topicRes != nil {
			result = append(result, topicRes)
		}
	}

	return result, nil
}
