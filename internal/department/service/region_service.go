package service

import (
	"context"
	"fmt"
	"media-service/internal/department/dto/request"
	"media-service/internal/department/dto/response"
	"media-service/internal/department/mapper"
	"media-service/internal/department/model"
	"media-service/internal/department/repository"
	"media-service/internal/gateway"
	"time"
)

type RegionService interface {
	CreateRegion(ctx context.Context, req request.CreateRegionRequest) error
	UpdateRegionName(ctx context.Context, req request.UpdateRegionRequest) error
	GetAll4Web(ctx context.Context) ([]*response.RegionResponseDTO, error)
}

type regionService struct {
	repo   repository.RegionRepository
	userGW gateway.UserGateway
}

func NewRegionService(repo repository.RegionRepository, userGw gateway.UserGateway) RegionService {
	return &regionService{
		repo:   repo,
		userGW: userGw,
	}
}

func (s *regionService) CreateRegion(ctx context.Context, req request.CreateRegionRequest) error {
	currentUser, err := s.userGW.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("get current user info failed")
	}

	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil
	}

	region := model.Region{
		Name:           req.Name,
		OrganizationID: currentUser.OrganizationAdmin.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.repo.CreateRegion(ctx, &region)
}

func (s *regionService) UpdateRegionName(ctx context.Context, req request.UpdateRegionRequest) error {
	return s.repo.UpdateRegionName(ctx, req.ID, req.Name)
}

func (s *regionService) GetAll4Web(ctx context.Context) ([]*response.RegionResponseDTO, error) {
	currentUser, err := s.userGW.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current user info failed")
	}

	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, nil
	}

	regions, err := s.repo.GetAllByOrgID(ctx, currentUser.OrganizationAdmin.ID)
	if err != nil {
		return nil, err
	}
	return mapper.MapRegionsToResponse(regions), nil
}
