package usecase

import (
	"context"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
)

type GetTopicAppUseCase interface {
	GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error)
}

type getTopicAppUseCase struct {
	topicRepo    repository.TopicRepository
	cachedUserGw gateway.UserGateway
}

func NewGetTopicAppUseCase(topicRepo repository.TopicRepository, cachedUserGw gateway.UserGateway) GetTopicAppUseCase {
	return &getTopicAppUseCase{
		topicRepo:    topicRepo,
		cachedUserGw: cachedUserGw,
	}
}

func (s *getTopicAppUseCase) GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error) {
	// get org by student
	student, err := s.cachedUserGw.GetStudentInfo(ctx, studentID)
	if err != nil {
		return nil, err
	}
	topics, err := s.topicRepo.GetAllTopicByOrganizationIDAndIsPublished(ctx, student.OrganizationID)
	if err != nil {
		return nil, err
	}
	appLanguage := helper.GetAppLanguage(ctx, 1)

	return mapper.ToTopic4StudentResponses4App(topics, appLanguage), nil
}
