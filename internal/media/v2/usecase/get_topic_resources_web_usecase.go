package usecase

import (
	"context"
	"fmt"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
	"media-service/pkg/constants"
)

type GetTopicResourcesWebUseCase interface {
	GetTopicResourcesByTopicAndStudent4Web(ctx context.Context, topicID string, studentID string) ([]*response.GetTopicResourcesResponse4Web, error)
	GetTopicResourcesByTopic4Web(ctx context.Context, topicID string) ([]*response.GetTopicResourcesResponse4Web, error)
	GetOutputResources4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error)
}

type getTopicResourcesWebUseCase struct {
	topicResourceRepository repository.TopicResourceRepository
	s3Service               s3.Service
}

func NewGetTopicResourcesWebUseCase(topicResourceRepository repository.TopicResourceRepository, s3Service s3.Service) GetTopicResourcesWebUseCase {
	return &getTopicResourcesWebUseCase{topicResourceRepository: topicResourceRepository, s3Service: s3Service}
}

func (uc *getTopicResourcesWebUseCase) GetTopicResourcesByTopicAndStudent4Web(ctx context.Context, topicID string, studentID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}
	topicResources, err := uc.topicResourceRepository.GetTopicResouresByTopicAndStudent(ctx, topicID, studentID)
	if err != nil {
		return nil, err
	}

	result := make([]*response.GetTopicResourcesResponse4Web, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		var imageUrl string
		if tr.ImageKey != "" {
			if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
				imageUrl = *url
			}
		}
		result = append(result, mapper.ToGetTopicResourcesResponse4Web(ctx, tr, imageUrl))
	}
	return result, nil
}

func (uc *getTopicResourcesWebUseCase) GetTopicResourcesByTopic4Web(ctx context.Context, topicID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	// currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	// if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
	// 	return nil, fmt.Errorf("access denied")
	// }
	topicResources, err := uc.topicResourceRepository.GetTopicResouresByTopic(ctx, topicID)
	if err != nil {
		return nil, err
	}
	result := make([]*response.GetTopicResourcesResponse4Web, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		var imageUrl string
		if tr.ImageKey != "" {
			if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
				imageUrl = *url
			}
		}
		result = append(result, mapper.ToGetTopicResourcesResponse4Web(ctx, tr, imageUrl))
	}
	return result, nil
}

func (uc *getTopicResourcesWebUseCase) GetOutputResources4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	topicResources, err := uc.topicResourceRepository.GetTopicResources(ctx, topicID, studentID)
	if err != nil {
		return nil, err
	}
	result := make([]*response.GetTopicResourcesResponse4Web, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		var imageUrl string
		if tr.ImageKey != "" {
			if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
				imageUrl = *url
			}
		}
		if tr.IsOutput {
			result = append(result, mapper.ToGetTopicResourcesResponse4Web(ctx, tr, imageUrl))
		}
	}
	return result, nil
}
