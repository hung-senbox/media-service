package usecase

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	s3svc "media-service/internal/s3"
	"media-service/logger"
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

	uc.populateMediaUrlsForTopic(ctx, topic)

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
		uc.populateMediaUrlsForTopic(ctx, &topic)
		topicRes := mapper.ToTopicResponses4GW(&topic, appLang)
		if topicRes != nil {
			result = append(result, topicRes)
		}
	}

	return result, nil
}

func (uc *getTopicGatewayUseCase) populateMediaUrlsForTopic(ctx context.Context, topic *model.Topic) {
	for li := range topic.LanguageConfig {
		langCfg := &topic.LanguageConfig[li]

		// images
		for ii := range langCfg.Images {
			img := &langCfg.Images[ii]
			if img.ImageKey != "" {
				url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
				if err == nil && url != nil {
					img.UploadedUrl = *url
				} else {
					logger.WriteLogEx("get_topic_web_usecase", "populateMediaUrlsForTopic_images", fmt.Sprintf("error getting image url: %v", err))
				}
			}
		}

		// video
		if langCfg.Video.VideoKey != "" {
			url, err := uc.s3Service.Get(ctx, langCfg.Video.VideoKey, nil)
			if err == nil && url != nil {
				langCfg.Video.UploadedUrl = *url
			} else {
				logger.WriteLogEx("get_topic_web_usecase", "populateMediaUrlsForTopic_video", fmt.Sprintf("error getting image url: %v", err))
			}
		}

		// audio
		if langCfg.Audio.AudioKey != "" {
			url, err := uc.s3Service.Get(ctx, langCfg.Audio.AudioKey, nil)
			if err == nil && url != nil {
				langCfg.Audio.UploadedUrl = *url
			} else {
				logger.WriteLogEx("get_topic_web_usecase", "populateMediaUrlsForTopic_audio", fmt.Sprintf("error getting image url: %v", err))
			}
		}
	}
}
