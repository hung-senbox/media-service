package usecase

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
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
	fileGateway       gateway.FileGateway
}

func NewGetTopicWebUseCase(
	topicRepo repository.TopicRepository,
	topicResourceRepo repository.TopicResourceRepository,
	cachedUserGw gateway.UserGateway,
	fileGateway gateway.FileGateway) GetTopicWebUseCase {
	return &getTopicWebUseCase{
		topicRepo:         topicRepo,
		topicResourceRepo: topicResourceRepo,
		cachedUserGw:      cachedUserGw,
		fileGateway:       fileGateway,
	}
}

func (uc *getTopicWebUseCase) GetTopics4Student4Web(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Web, error) {
	// get org by student
	student, err := uc.cachedUserGw.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, err
	}
	topics, err := uc.topicRepo.GetAllTopicByOrganizationIDAndIsPublished(ctx, student.OrganizationID)
	if err != nil {
		return nil, err
	}
	for ti := range topics {
		for li := range topics[ti].LanguageConfig {
			langCfg := &topics[ti].LanguageConfig[li]
			for ii := range langCfg.Images {
				img := &langCfg.Images[ii]
				if img.ImageKey != "" {
					url, err := uc.fileGateway.GetImageUrl(ctx, gw_request.GetFileUrlRequest{
						Key:  img.ImageKey,
						Mode: "private",
					})
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

	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}
	orgID := currentUser.OrganizationAdmin.ID
	if studentID != "" {
		return uc.getTopicsByStudentFromResource4Web(ctx, studentID)
	}

	topics, err := uc.topicRepo.GetAllTopicByOrganizationID(ctx, orgID)
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
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)

	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}

	topics, err := uc.topicRepo.GetAllTopicByOrganizationIDAndIsPublished(ctx, currentUser.OrganizationAdmin.ID)
	if err != nil {
		return nil, err
	}

	for ti := range topics {
		for li := range topics[ti].LanguageConfig {
			langCfg := &topics[ti].LanguageConfig[li]
			for ii := range langCfg.Images {
				img := &langCfg.Images[ii]
				if img.ImageKey != "" {
					url, err := uc.fileGateway.GetImageUrl(ctx, gw_request.GetFileUrlRequest{
						Key:  img.ImageKey,
						Mode: "private",
					})
					if err == nil && url != nil {
						img.UploadedUrl = *url
					}
				}
			}
		}
	}

	return mapper.ToTopic2Assign4Web(topics, 1), nil
}

func (uc *getTopicWebUseCase) getTopicsByStudentFromResource4Web(ctx context.Context, studentID string) ([]response.TopicResponse4Web, error) {

	topicResources, err := uc.topicResourceRepo.GetTopicResouresByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// loc ra topicResources co topicID khac nhau
	topicIDs := make([]string, 0)
	for _, topicResource := range topicResources {
		topicIDs = append(topicIDs, topicResource.TopicID)
	}
	topicIDs = helper.RemoveDuplicateString(topicIDs)

	var topics []model.Topic
	for _, topicID := range topicIDs {
		topic, _ := uc.topicRepo.GetTopicByID(ctx, topicID)
		if topic != nil {
			topics = append(topics, *topic)
		}
	}

	for ti := range topics {
		uc.populateMediaUrlsForTopic(ctx, &topics[ti])
	}

	return mapper.ToTopicResponses4Web(topics), nil
}

// populateMediaUrlsForTopic enriches a topic's language configs with signed media URLs when keys exist
func (uc *getTopicWebUseCase) populateMediaUrlsForTopic(ctx context.Context, topic *model.Topic) {
	for li := range topic.LanguageConfig {
		langCfg := &topic.LanguageConfig[li]

		// images
		for ii := range langCfg.Images {
			img := &langCfg.Images[ii]
			if img.ImageKey != "" {
				url, err := uc.fileGateway.GetImageUrl(ctx, gw_request.GetFileUrlRequest{
					Key:  img.ImageKey,
					Mode: "private",
				})
				if err == nil && url != nil {
					img.UploadedUrl = *url
				}
			}
		}

		// video
		if langCfg.Video.VideoKey != "" {
			url, err := uc.fileGateway.GetVideoUrl(ctx, gw_request.GetFileUrlRequest{
				Key:  langCfg.Video.VideoKey,
				Mode: "private",
			})
			if err == nil && url != nil {
				langCfg.Video.UploadedUrl = *url
			}
		}

		// audio
		if langCfg.Audio.AudioKey != "" {
			url, err := uc.fileGateway.GetAudioUrl(ctx, gw_request.GetFileUrlRequest{
				Key:  langCfg.Audio.AudioKey,
				Mode: "private",
			})
			if err == nil && url != nil {
				langCfg.Audio.UploadedUrl = *url
			}
		}
	}
}
