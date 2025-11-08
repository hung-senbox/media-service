package usecase

import (
	"context"
	"fmt"
	"media-service/internal/gateway"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	s3svc "media-service/internal/s3"
	"media-service/logger"
	"media-service/pkg/constants"
)

type GetTopicWebUseCase interface {
	GetTopics4Student4Web(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Web, error)
	GetTopics4Web(ctx context.Context, studentID string) ([]response.TopicResponse4Web, error)
	GetTopic4Web(ctx context.Context, topicID string) (*response.TopicResponse4Web, error)
	GetTopics2Assign4Web(ctx context.Context) ([]*response.TopicResponse2Assign4Web, error)
}

type getTopicWebUseCase struct {
	topicRepo         repository.TopicRepository
	topicResourceRepo repository.TopicResourceRepository
	cachedUserGw      gateway.UserGateway
	s3Service         s3svc.Service
}

func NewGetTopicWebUseCase(
	topicRepo repository.TopicRepository,
	topicResourceRepo repository.TopicResourceRepository,
	cachedUserGw gateway.UserGateway,
	s3Service s3svc.Service) GetTopicWebUseCase {
	return &getTopicWebUseCase{
		topicRepo:         topicRepo,
		topicResourceRepo: topicResourceRepo,
		cachedUserGw:      cachedUserGw,
		s3Service:         s3Service,
	}
}

func (uc *getTopicWebUseCase) GetTopics4Student4Web(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Web, error) {
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
	return mapper.ToTopic4StudentResponses4Web(topics, 1), nil
}

func (uc *getTopicWebUseCase) GetTopics4Web(ctx context.Context, studentID string) ([]response.TopicResponse4Web, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser == nil {
		return nil, fmt.Errorf("access denied")
	}

	if currentUser.IsSuperAdmin {
		topics, err := uc.topicRepo.GetAllTopics(ctx)
		if err != nil {
			return nil, err
		}
		for ti := range topics {
			uc.populateMediaUrlsForTopic(ctx, &topics[ti])
		}
		return mapper.ToTopicResponses4Web(topics), nil
	}

	if currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}

	topics, err := uc.topicRepo.GetAllTopicsIsPublished(ctx)
	if err != nil {
		return nil, err
	}

	for ti := range topics {
		uc.populateMediaUrlsForTopic(ctx, &topics[ti])
	}

	return mapper.ToTopicResponses4Web(topics), nil

}

func (uc *getTopicWebUseCase) GetTopic4Web(ctx context.Context, topicID string) (*response.TopicResponse4Web, error) {

	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("get topic failed: %w", err)
	}

	uc.populateMediaUrlsForTopic(ctx, topic)

	return mapper.ToTopicResponse4Web(topic), nil
}

func (uc *getTopicWebUseCase) GetTopics2Assign4Web(ctx context.Context) ([]*response.TopicResponse2Assign4Web, error) {
	// currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)

	// if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
	// 	return nil, fmt.Errorf("access denied")
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

	return mapper.ToTopic2Assign4Web(topics, 1), nil
}

// populateMediaUrlsForTopic enriches a topic's language configs with signed media URLs when keys exist
func (uc *getTopicWebUseCase) populateMediaUrlsForTopic(ctx context.Context, topic *model.Topic) {
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
