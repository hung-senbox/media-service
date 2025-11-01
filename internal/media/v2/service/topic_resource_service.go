package service

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicResourceService interface {
	CreateTopicResource(ctx context.Context, req request.CreateTopicResourceRequest) (string, error)
	GetTopicResources(ctx context.Context, topicID, studentID, orgID string) ([]*response.GetTopicResourceResponse, error)
	GetTopicResource(ctx context.Context, topicResourceID, orgID string) (*response.GetTopicResourceResponse, error)
	UpdateTopicResource(ctx context.Context, topicResourceID string, req request.UpdateTopicResourceRequest) (string, error)
	DeleteTopicResource(ctx context.Context, topicResourceID string) error
}

type topicResourceService struct {
	topicResourceRepository repository.TopicResourceRepository
	topicRepository         repository.TopicRepository
	fileGateway             gateway.FileGateway
	userGw                  gateway.UserGateway
}

func NewTopicResourceService(topicResourceRepository repository.TopicResourceRepository, topicRepository repository.TopicRepository, fileGateway gateway.FileGateway, userGw gateway.UserGateway) TopicResourceService {
	return &topicResourceService{topicResourceRepository: topicResourceRepository, topicRepository: topicRepository, fileGateway: fileGateway, userGw: userGw}
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

	resp, err := s.fileGateway.UploadImage(ctx, gw_request.UploadFileRequest{
		File:     req.File,
		Folder:   "topic_resource",
		FileName: req.FileName,
		Mode:     "private",
	})
	if err != nil {
		return "", err
	}

	ID := primitive.NewObjectID()

	topicResource := &model.TopicResource{
		ID:        ID,
		TopicID:   req.TopicID,
		StudentID: req.StudentID,
		FileName:  req.FileName,
		ImageKey:  resp.Key,
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

	result := mapper.ToGetTopicResourceResponses(ctx, orgID, topicResources, s.topicRepository, s.userGw, s.fileGateway)

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

	return mapper.ToGetTopicResourceResponse(ctx, orgID, topicResource, s.topicRepository, s.userGw, s.fileGateway), nil
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
			err = s.fileGateway.DeleteImage(ctx, topicResource.ImageKey)
			if err != nil {
				return "", err
			}
		}
		resp, err := s.fileGateway.UploadImage(ctx, gw_request.UploadFileRequest{
			File:     req.File,
			Folder:   "topic_resource",
			FileName: req.FileName,
			Mode:     "private",
		})
		if err != nil {
			return "", err
		}
		topicResource.ImageKey = resp.Key
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
		err = s.fileGateway.DeleteImage(ctx, topicResource.ImageKey)
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
