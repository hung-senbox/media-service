package domain

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/pdf/domain/dto"
	"media-service/internal/pdf/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResourceService interface {
	CreateResource(ctx context.Context, req dto.CreateResourceRequest) (string, error)
	GetResources(ctx context.Context, role, organizationID string) (*dto.GroupedResourceResponse, error)
	UploadDocumentToResource(ctx context.Context, id string, req dto.UpdateResourceRequest) (string, error)
	UploadSignatureToResource(ctx context.Context, id string, req dto.UploadSignatureRequest) (string, error)
	DeleteResource(ctx context.Context, id string) error
}

type userResourceService struct {
	UserResourceRepository UserResourceRepository
	fileGateway            gateway.FileGateway
	userGateway            gateway.UserGateway
}

func NewUserResourceService(userResourceRepository UserResourceRepository,
	fileGateway gateway.FileGateway,
	userGateway gateway.UserGateway) UserResourceService {
	return &userResourceService{
		UserResourceRepository: userResourceRepository,
		fileGateway:            fileGateway,
		userGateway:            userGateway,
	}
}

func (s *userResourceService) CreateResource(ctx context.Context, req dto.CreateResourceRequest) (string, error) {

	if req.Folder == "" {
		return "", fmt.Errorf("folder cannot be empty")
	}

	if req.Role == "" {
		return "", fmt.Errorf("role cannot be empty")
	}

	if req.Type == "" {
		return "", fmt.Errorf("type cannot be empty")
	}

	var uploaderData *model.Owner
	var targetData *model.Owner

	switch req.Role {
	case "teacher":
		data, err := s.userGateway.GetTeacherByUserAndOrganization(ctx, helper.GetUserID(ctx), req.OrganizationID)
		if err != nil {
			return "", err
		}
		req.UploaderID = &model.Owner{OwnerID: data.ID, OwnerRole: "teacher"}
		uploaderData = req.UploaderID

		if req.TargetID != nil {
			targetData = req.TargetID
		} else {
			targetData = nil
		}

	case "staff":
		data, err := s.userGateway.GetStaffByUserAndOrganization(ctx, helper.GetUserID(ctx), req.OrganizationID)
		if err != nil {
			return "", err
		}
		req.UploaderID = &model.Owner{OwnerID: data.ID, OwnerRole: "staff"}
		uploaderData = req.UploaderID

		if req.TargetID != nil {
			targetData = req.TargetID
		} else {
			targetData = nil
		}

	case "parent":
		data, err := s.userGateway.GetParentByUser(ctx, helper.GetUserID(ctx))
		if err != nil {
			return "", err
		}
		req.UploaderID = &model.Owner{OwnerID: data.ID, OwnerRole: "parent"}
		uploaderData = req.UploaderID

		if req.TargetID != nil {
			targetData = req.TargetID
		} else {
			targetData = nil
		}

	default:
		return "", fmt.Errorf("unsupported role: %s", req.Role)
	}

	ID := primitive.NewObjectID()

	err := s.UserResourceRepository.CreateResource(ctx, &model.UserResource{
		ID:           ID,
		UploaderID:   uploaderData,
		TargetID:     targetData,
		Type:         req.Type,
		ResourceType: "",
		FileName:     nil,
		Folder:       req.Folder,
		Color:        req.Color,
		SignatureKey: nil,
		URL:          nil,
		PDFKey:       nil,
		CreatedBy:    helper.GetUserID(ctx),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		return "", err
	}

	return ID.Hex(), nil

}

func (s *userResourceService) GetResources(ctx context.Context, role, organizationID string) (*dto.GroupedResourceResponse, error) {

	if role == "" {
		return nil, fmt.Errorf("role cannot be empty")
	}

	var (
		selfResources    []*model.UserResource
		relatedResources []*model.UserResource
	)

	switch role {
	case "teacher":
		teacher, err := s.userGateway.GetTeacherByUserAndOrganization(ctx, helper.GetUserID(ctx), organizationID)
		if err != nil {
			return nil, err
		}
		selfResources, _ = s.UserResourceRepository.GetSelfResources(ctx, teacher.ID)
		relatedResources, _ = s.UserResourceRepository.GetRelatedResources(ctx, teacher.ID)
	case "staff":
		staff, err := s.userGateway.GetStaffByUserAndOrganization(ctx, helper.GetUserID(ctx), organizationID)
		if err != nil {
			return nil, err
		}
		selfResources, _ = s.UserResourceRepository.GetSelfResources(ctx, staff.ID)
		relatedResources, _ = s.UserResourceRepository.GetRelatedResources(ctx, staff.ID)
	case "parent":
		parent, err := s.userGateway.GetParentByUser(ctx, helper.GetUserID(ctx))
		if err != nil {
			return nil, err
		}
		selfResources, _ = s.UserResourceRepository.GetSelfResources(ctx, parent.ID)
		relatedResources, _ = s.UserResourceRepository.GetRelatedResources(ctx, parent.ID)
	default:
		return nil, fmt.Errorf("unsupported role: %s", role)
	}
	

	return &dto.GroupedResourceResponse{
		SelfResources:    dto.ToResourceResponses(ctx, organizationID, selfResources, s.userGateway, s.fileGateway),
		RelatedResources: dto.ToResourceResponses(ctx, organizationID, relatedResources, s.userGateway, s.fileGateway),
	}, nil

}

func (s *userResourceService) UploadDocumentToResource(ctx context.Context, id string, req dto.UpdateResourceRequest) (string, error) {

	if req.ResourceType == "" {
		return "", fmt.Errorf("resource type cannot be empty")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	pdfData, err := s.UserResourceRepository.GetResourceByID(ctx, objectID)
	if err != nil {
		return "", err
	}

	if pdfData.PDFKey != nil {
		err = s.fileGateway.DeletePDF(ctx, *pdfData.PDFKey)
		if err != nil {
			return "", err
		}
	}

	if pdfData == nil {
		return "", fmt.Errorf("pdf not found")
	}

	if req.ResourceType == "pdf" && req.File != nil {
		resource, err := s.UserResourceRepository.GetResourceByID(ctx, objectID)
		if err != nil {
			return "", err
		}
		if resource == nil {
			return "", fmt.Errorf("resource not found")
		}

		resp, err := s.fileGateway.UploadPDF(ctx, gw_request.UploadFileRequest{
			File:     req.File,
			Folder:   "pdf_media",
			FileName: *req.FileName,
			Mode:     "private",
		})
		if err != nil {
			return "", err
		}

		resource.FileName = req.FileName
		resource.ResourceType = req.ResourceType
		resource.PDFKey = &resp.Key
		resource.URL = nil
		resource.UpdatedAt = time.Now()

		err = s.UserResourceRepository.UpdateResourceByID(ctx, objectID, resource)
		if err != nil {
			return "", err
		}

		return resp.Key, nil

	} else if req.ResourceType == "url" && req.Url != nil {
		resource, err := s.UserResourceRepository.GetResourceByID(ctx, objectID)
		if err != nil {
			return "", err
		}
		if resource == nil {
			return "", fmt.Errorf("resource not found")
		}

		resource.ResourceType = req.ResourceType
		resource.URL = req.Url
		resource.PDFKey = nil
		resource.FileName = nil
		resource.UpdatedAt = time.Now()

		err = s.UserResourceRepository.UpdateResourceByID(ctx, objectID, resource)
		if err != nil {
			return "", err
		}

		return *req.Url, nil
	} else {
		return "", fmt.Errorf("resource type not supported")
	}
}

func (s *userResourceService) UploadSignatureToResource(ctx context.Context, id string, req dto.UploadSignatureRequest) (string, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	pdfData, err := s.UserResourceRepository.GetResourceByID(ctx, objectID)
	if err != nil {
		return "", err
	}

	if pdfData == nil {
		return "", fmt.Errorf("pdf not found")
	}

	if pdfData.SignatureKey != nil {
		err = s.fileGateway.DeleteImage(ctx, *pdfData.SignatureKey)
		if err != nil {
			return "", err
		}
	}

	pdfData.SignatureKey = &req.SignatureKey
	pdfData.UpdatedAt = time.Now()

	err = s.UserResourceRepository.UpdateResourceByID(ctx, objectID, pdfData)
	if err != nil {
		return "", err
	}

	return req.SignatureKey, nil

}

func (s *userResourceService) DeleteResource(ctx context.Context, id string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	resource, err := s.UserResourceRepository.GetResourceByID(ctx, objectID)
	if err != nil {
		return err
	}

	if resource == nil {
		return fmt.Errorf("resource not found")
	}

	if resource.PDFKey != nil {
		err = s.fileGateway.DeletePDF(ctx, *resource.PDFKey)
		if err != nil {
			return err
		}
	}

	if resource.SignatureKey != nil {
		err = s.fileGateway.DeleteImage(ctx, *resource.SignatureKey)
		if err != nil {
			return err
		}
	}

	err = s.UserResourceRepository.DeleteResourceByID(ctx, objectID)
	if err != nil {
		return err
	}

	return nil

}
