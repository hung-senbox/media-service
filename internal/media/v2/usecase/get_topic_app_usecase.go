package usecase

import (
	"context"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
)

type GetTopicAppUseCase interface {
	GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error)
	GetTopics4App(ctx context.Context, organizationID string) ([]*response.GetTopic4StudentResponse4App, error)
}

type getTopicAppUseCase struct {
	topicRepo    repository.TopicRepository
	cachedUserGw gateway.UserGateway
	fileGateway  gateway.FileGateway
}

func NewGetTopicAppUseCase(topicRepo repository.TopicRepository, cachedUserGw gateway.UserGateway, fileGateway gateway.FileGateway) GetTopicAppUseCase {
	return &getTopicAppUseCase{
		topicRepo:    topicRepo,
		cachedUserGw: cachedUserGw,
		fileGateway:  fileGateway,
	}
}

func (uc *getTopicAppUseCase) GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error) {
	// get org by student
	// student, err := uc.cachedUserGw.GetStudentInfo(ctx, studentID)
	// if err != nil {
	// 	return nil, err
	// }
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

	return mapper.ToTopic4StudentResponses4App(topics, appLanguage), nil
}
