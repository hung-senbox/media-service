package usecase

import (
	"context"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"time"
)

type GetTopicResourceAppUseCase interface {
	GetOutputResources4App(ctx context.Context, studentID string, month, year int) ([]*response.GetTopicResourcesResponse4App, error)
}

type getTopicResourceAppUseCase struct {
	topicResourceRepository repository.TopicResourceRepository
	fileGateway             gateway.FileGateway
}

func NewGetTopicResourceAppUseCase(topicResourceRepository repository.TopicResourceRepository) GetTopicResourceAppUseCase {
	return &getTopicResourceAppUseCase{topicResourceRepository: topicResourceRepository}
}

func (uc *getTopicResourceAppUseCase) GetOutputResources4App(ctx context.Context, studentID string, month, year int) ([]*response.GetTopicResourcesResponse4App, error) {
	topicResources, err := uc.topicResourceRepository.GetTopicResouresByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	// filter topic resources by month
	topicResources = filterTopicResourcesByMonth(topicResources, month, year)
	result := make([]*response.GetTopicResourcesResponse4App, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		var imageUrl string
		if tr.ImageKey != "" {
			if url, err := uc.fileGateway.GetImageUrl(ctx, gw_request.GetFileUrlRequest{Key: tr.ImageKey, Mode: "private"}); err == nil && url != nil {
				imageUrl = *url
			}
		}
		if tr.IsOutput {
			result = append(result, mapper.ToGetTopicResourcesResponse4App(ctx, tr, imageUrl))
		}
	}
	return result, nil
}

func filterTopicResourcesByMonth(topicResources []*model.TopicResource, month, year int) []*model.TopicResource {
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
