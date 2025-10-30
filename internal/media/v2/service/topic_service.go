package service

import (
	"context"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/usecase"
)

type TopicService interface {
	UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error)
	GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error)
	GetTopics4Web(ctx context.Context) ([]response.TopicResponse4Web, error)
	GetTopic4Web(ctx context.Context, topicID string) (*response.TopicResponse4Web, error)
	GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error)
	GetTopics4Student4Web(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Web, error)
	GetTopics4Student4Gw(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Gw, error)
	GetTopic4Gw(ctx context.Context, topicID string) (*response.TopicResponse4GW, error)
	GetAllTopicsByOrganization4Gw(ctx context.Context, organizationID string) ([]*response.TopicResponse4GW, error)
	GetTopics2Assign4Web(ctx context.Context) ([]*response.TopicResponse2Assign4Web, error)
}

type topicService struct {
	uploadTopicUseCase       usecase.UploadTopicUseCase
	getUploadProgressUseCase usecase.GetUploadProgressUseCase
	getTopicAppUseCase       usecase.GetTopicAppUseCase
	getTopicWebUseCase       usecase.GetTopicWebUseCase
	getTopicGatewayUseCase   usecase.GetTopicGatewayUseCase
}

func NewTopicService(
	uploadTopicUseCase usecase.UploadTopicUseCase,
	getUploadProgressUseCase usecase.GetUploadProgressUseCase,
	getTopicAppUseCase usecase.GetTopicAppUseCase,
	getTopicWebUseCase usecase.GetTopicWebUseCase,
	getTopicGatewayUseCase usecase.GetTopicGatewayUseCase,
) TopicService {
	return &topicService{
		uploadTopicUseCase:       uploadTopicUseCase,
		getUploadProgressUseCase: getUploadProgressUseCase,
		getTopicAppUseCase:       getTopicAppUseCase,
		getTopicWebUseCase:       getTopicWebUseCase,
		getTopicGatewayUseCase:   getTopicGatewayUseCase,
	}
}

// ------------------- Upload Topic -------------------
func (s *topicService) UploadTopic(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	return s.uploadTopicUseCase.UploadTopic(ctx, req)
}

// ------------------- Get upload progress -------------------
func (s *topicService) GetUploadProgress(ctx context.Context, topicID string) (*response.GetUploadProgressResponse, error) {
	return s.getUploadProgressUseCase.GetUploadProgress(ctx, topicID)
}

// =============== Get Topic 4 App ================
func (s *topicService) GetTopics4Student4App(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4App, error) {
	return s.getTopicAppUseCase.GetTopics4Student4App(ctx, studentID)
}

// =============== Get Topic 4 Web ================
func (s *topicService) GetTopics4Web(ctx context.Context) ([]response.TopicResponse4Web, error) {
	return s.getTopicWebUseCase.GetTopics4Web(ctx)
}
func (s *topicService) GetTopic4Web(ctx context.Context, topicID string) (*response.TopicResponse4Web, error) {
	return s.getTopicWebUseCase.GetTopic4Web(ctx, topicID)
}
func (s *topicService) GetTopics4Student4Web(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Web, error) {
	return s.getTopicWebUseCase.GetTopics4Student4Web(ctx, studentID)
}
func (s *topicService) GetTopics2Assign4Web(ctx context.Context) ([]*response.TopicResponse2Assign4Web, error) {
	return s.getTopicWebUseCase.GetTopics2Assign4Web(ctx)
}

// =============== Get Topic 4 Gateway ================
func (s *topicService) GetTopics4Student4Gw(ctx context.Context, studentID string) ([]*response.GetTopic4StudentResponse4Gw, error) {
	return s.getTopicGatewayUseCase.GetTopics4Student4Gw(ctx, studentID)
}
func (s *topicService) GetTopic4Gw(ctx context.Context, topicID string) (*response.TopicResponse4GW, error) {
	return s.getTopicGatewayUseCase.GetTopic4Gw(ctx, topicID)
}
func (s *topicService) GetAllTopicsByOrganization4Gw(ctx context.Context, organizationID string) ([]*response.TopicResponse4GW, error) {
	return s.getTopicGatewayUseCase.GetAllTopicsByOrganization4Gw(ctx, organizationID)
}
