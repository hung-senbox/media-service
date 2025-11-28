package service

import (
	"context"
	"fmt"
	"io"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/media/v2/usecase"
	"media-service/internal/s3"
	"media-service/pkg/uploader"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicResourceService interface {
	CreateTopicResource(ctx context.Context, req request.CreateTopicResourceRequest) (string, error)
	GetTopicResources(ctx context.Context, topicID, studentID, orgID string) ([]*response.GetTopicResourceResponse, error)
	GetTopicResource(ctx context.Context, topicResourceID, orgID string) (*response.GetTopicResourceResponse, error)
	UpdateTopicResource(ctx context.Context, topicResourceID string, req request.UpdateTopicResourceRequest) (string, error)
	DeleteTopicResource(ctx context.Context, topicResourceID string) error
	GetTopicResourcesByTopic4Web(ctx context.Context, topicID string) ([]*response.GetTopicResourcesResponse4Web, error)
	GetTopicResourcesByTopicAndStudent4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error)
	SetOutputTopicResource(ctx context.Context, req request.SetOutputTopicResourceRequest) error
	GetOutputResources4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error)
	GetOutputResources4App(ctx context.Context, studentID string, day, month, year int, topicID string) ([]*response.GetTopicResourcesResponse4App, error)
	OffOutputTopicResource(ctx context.Context, topicResourceID string) error
	GetTopicResourcesByStudent4Web(ctx context.Context, studentID string) ([]*response.GetTopicResourcesResponseByStudent4Web, error)
}

type topicResourceService struct {
	topicResourceRepository     repository.TopicResourceRepository
	topicRepository             repository.TopicRepository
	s3Service                   s3.Service
	userGw                      gateway.UserGateway
	getTopicResourcesWebUseCase usecase.GetTopicResourcesWebUseCase
	getTopicResourceAppUseCase  usecase.GetTopicResourceAppUseCase
}

func NewTopicResourceService(
	topicResourceRepository repository.TopicResourceRepository,
	topicRepository repository.TopicRepository,
	s3Service s3.Service,
	userGw gateway.UserGateway,
	getTopicResourcesWebUseCase usecase.GetTopicResourcesWebUseCase,
	getTopicResourceAppUseCase usecase.GetTopicResourceAppUseCase,
) TopicResourceService {
	return &topicResourceService{
		topicResourceRepository:     topicResourceRepository,
		topicRepository:             topicRepository,
		s3Service:                   s3Service,
		userGw:                      userGw,
		getTopicResourcesWebUseCase: getTopicResourcesWebUseCase,
		getTopicResourceAppUseCase:  getTopicResourceAppUseCase,
	}
}

func (s *topicResourceService) CreateTopicResource(ctx context.Context, req request.CreateTopicResourceRequest) (string, error) {

	if req.TopicID == "" {
		return "", fmt.Errorf("topic id is required")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student id is required")
	}

	if req.FileName == "" {
		return "", fmt.Errorf("file name is required")
	}

	if req.File == nil {
		return "", fmt.Errorf("file is required")
	}

	key := helper.BuildObjectKeyS3("topic_resource", req.File.Filename, req.FileName)
	file, err := req.File.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	_, err = s.s3Service.Save(ctx, bytes, key, uploader.UploadPrivate)
	if err != nil {
		return "", err
	}

	ID := primitive.NewObjectID()

	topicResource := &model.TopicResource{
		ID:        ID,
		TopicID:   req.TopicID,
		StudentID: req.StudentID,
		FileName:  req.FileName,
		ImageKey:  key,
		CreatedBy: helper.GetUserID(ctx),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.topicResourceRepository.CreateTopicResource(ctx, topicResource)
	if err != nil {
		return "", err
	}

	return ID.Hex(), nil
}

func (s *topicResourceService) GetTopicResources(ctx context.Context, topicID, studentID, orgID string) ([]*response.GetTopicResourceResponse, error) {

	if orgID == "" {
		return nil, fmt.Errorf("organization id is required")
	}

	topicResources, err := s.topicResourceRepository.GetTopicResources(ctx, topicID, studentID)
	if err != nil {
		return nil, err
	}

	result := mapper.ToGetTopicResourceResponses(ctx, orgID, topicResources, s.topicRepository, s.userGw, s.s3Service)

	return result, nil
}

func (s *topicResourceService) GetTopicResource(ctx context.Context, topicResourceID, orgID string) (*response.GetTopicResourceResponse, error) {

	if orgID == "" {
		return nil, fmt.Errorf("organization id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(topicResourceID)
	if err != nil {
		return nil, err
	}

	topicResource, err := s.topicResourceRepository.GetTopicResource(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if topicResource == nil {
		return nil, fmt.Errorf("topic resource not found")
	}

	return mapper.ToGetTopicResourceResponse(ctx, orgID, topicResource, s.topicRepository, s.userGw, s.s3Service), nil
}

func (s *topicResourceService) UpdateTopicResource(ctx context.Context, topicResourceID string, req request.UpdateTopicResourceRequest) (string, error) {

	objectID, err := primitive.ObjectIDFromHex(topicResourceID)
	if err != nil {
		return "", err
	}

	topicResource, err := s.topicResourceRepository.GetTopicResource(ctx, objectID)
	if err != nil {
		return "", err
	}

	if topicResource == nil {
		return "", fmt.Errorf("topic resource not found")
	}

	if req.FileName != "" {
		topicResource.FileName = req.FileName
	}

	if req.TopicID != "" {
		topicResource.TopicID = req.TopicID
	}

	if req.File != nil {
		if topicResource.ImageKey != "" {
			err = s.s3Service.Delete(ctx, topicResource.ImageKey)
			if err != nil {
				return "", err
			}
		}
		key := helper.BuildObjectKeyS3("topic_resource", req.File.Filename, req.FileName)
		f, err := req.File.Open()
		if err != nil {
			return "", err
		}
		defer f.Close()
		bs, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		_, err = s.s3Service.Save(ctx, bs, key, uploader.UploadPrivate)
		if err != nil {
			return "", err
		}
		topicResource.ImageKey = key
	}

	topicResource.UpdatedAt = time.Now()

	err = s.topicResourceRepository.UpdateTopicResource(ctx, objectID, topicResource)
	if err != nil {
		return "", err
	}

	return topicResource.ID.Hex(), nil
}

func (s *topicResourceService) DeleteTopicResource(ctx context.Context, topicResourceID string) error {
	objectID, err := primitive.ObjectIDFromHex(topicResourceID)
	if err != nil {
		return err
	}

	topicResource, err := s.topicResourceRepository.GetTopicResource(ctx, objectID)
	if err != nil {
		return err
	}

	if topicResource == nil {
		return fmt.Errorf("topic resource not found")
	}

	if topicResource.ImageKey != "" {
		err = s.s3Service.Delete(ctx, topicResource.ImageKey)
		if err != nil {
			return err
		}
	}

	err = s.topicResourceRepository.DeleteTopicResource(ctx, objectID)
	if err != nil {
		return err
	}
	return nil

}

func (s *topicResourceService) GetTopicResourcesByTopic4Web(ctx context.Context, topicID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	return s.getTopicResourcesWebUseCase.GetTopicResourcesByTopic4Web(ctx, topicID)
}

func (s *topicResourceService) GetTopicResourcesByTopicAndStudent4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	return s.getTopicResourcesWebUseCase.GetTopicResourcesByTopicAndStudent4Web(ctx, topicID, studentID)
}

func (s *topicResourceService) SetOutputTopicResource(ctx context.Context, req request.SetOutputTopicResourceRequest) error {

	objectID, err := primitive.ObjectIDFromHex(req.TopicResourceID)
	if err != nil {
		return err
	}

	topicResource, err := s.topicResourceRepository.GetTopicResource(ctx, objectID)
	if err != nil {
		return err
	}

	// neu target student va topic khong trung voi topicResource thi khong the set output
	if topicResource.StudentID != req.TargetStudentID || topicResource.TopicID != req.TargetTopicID {
		return fmt.Errorf("target student and topic do not match the topic resource")
	}

	if topicResource == nil {
		return fmt.Errorf("topic resource not found")
	}

	topicResource.IsOutput = true
	topicResource.UpdatedAt = time.Now()

	err = s.topicResourceRepository.UpdateTopicResource(ctx, objectID, topicResource)
	if err != nil {
		return err
	}

	return nil
}

func (s *topicResourceService) GetOutputResources4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	return s.getTopicResourcesWebUseCase.GetOutputResources4Web(ctx, topicID, studentID)
}

func (s *topicResourceService) GetOutputResources4App(ctx context.Context, studentID string, day, month, year int, topicID string) ([]*response.GetTopicResourcesResponse4App, error) {
	return s.getTopicResourceAppUseCase.GetOutputResources4App(ctx, studentID, day, month, year, topicID)
}

func (s *topicResourceService) OffOutputTopicResource(ctx context.Context, topicResourceID string) error {
	objectID, err := primitive.ObjectIDFromHex(topicResourceID)
	if err != nil {
		return err
	}

	topicResource, err := s.topicResourceRepository.GetTopicResource(ctx, objectID)
	if err != nil {
		return err
	}

	if topicResource == nil {
		return fmt.Errorf("topic resource not found")
	}

	topicResource.IsOutput = false
	topicResource.UpdatedAt = time.Now()

	err = s.topicResourceRepository.UpdateTopicResource(ctx, objectID, topicResource)
	if err != nil {
		return err
	}

	return nil
}

func (s *topicResourceService) GetTopicResourcesByStudent4Web(ctx context.Context, studentID string) ([]*response.GetTopicResourcesResponseByStudent4Web, error) {
	return s.getTopicResourcesWebUseCase.GetTopicResourcesByStudent4Web(ctx, studentID)
}
