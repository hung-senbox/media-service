package dto

import (
	"context"
	"media-service/internal/gateway"
	"media-service/internal/pdf/model"
	"media-service/internal/s3"
)

func ToResourceResponses(
	ctx context.Context,
	organizationID string,
	data []*model.UserResource,
	userGw gateway.UserGateway,
	s3Svc s3.Service,
) []*ResourceResponse {
	responses := make([]*ResourceResponse, 0, len(data))
	for _, r := range data {
		resp := &ResourceResponse{
			ID:             r.ID.Hex(),
			OrganizationID: r.Organization,
			ResourceType:   r.ResourceType,
			Folder:         r.Folder,
			Color:          r.Color,
			Status:         r.Status,
			IsDownloaded:   r.IsDownloaded,
			CreatedBy:      r.CreatedBy,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		}

		if r.UploaderID != nil {
			resp.UploaderInfor = getUserInfoByRole(ctx, userGw, r.UploaderID, organizationID)
		}

		if r.TargetID != nil {
			resp.TargetInfor = getUserInfoByRole(ctx, userGw, r.TargetID, organizationID)
		}

		if r.SignatureKey != nil {
			url, err := s3Svc.Get(ctx, *r.SignatureKey, nil)
			if err == nil {
				resp.SignatureUrl = url
			}
		}

		if r.PDFKey != nil && *r.PDFKey != "" {
			url, err := s3Svc.Get(ctx, *r.PDFKey, nil)
			if err == nil {
				resp.PDFUrl = url
			}
		}

		if r.URL != nil {
			resp.URL = r.URL
		}

		if r.FileName != nil {
			resp.FileName = r.FileName
		}

		responses = append(responses, resp)
	}

	return responses
}

func getUserInfoByRole(ctx context.Context, userGw gateway.UserGateway, owner *model.Owner, organizationID string) *UserInfor {
	if owner == nil {
		return nil
	}
	switch owner.OwnerRole {
	case "teacher":
		data, err := userGw.GetTeacherInfo(ctx, owner.OwnerID)
		if err != nil {
			return &UserInfor{ID: owner.OwnerID}
		}
		return &UserInfor{
			ID:             data.ID,
			Name:           data.Name,
			Code:           data.Code,
			OrganizationID: data.OrganizationID,
		}

	case "student":
		data, err := userGw.GetStudentInfo(ctx, owner.OwnerID)
		if err != nil {
			return &UserInfor{ID: owner.OwnerID}
		}
		return &UserInfor{
			ID:             data.ID,
			Name:           data.Name,
			Code:           data.Code,
			OrganizationID: data.OrganizationID,
		}

	case "parent":
		data, err := userGw.GetParentByUser(ctx, owner.OwnerID)
		if err != nil {
			return &UserInfor{ID: owner.OwnerID}
		}
		return &UserInfor{
			ID:             data.ID,
			Name:           data.Name,
			Code:           data.Code,
			OrganizationID: data.OrganizationID,
		}

	case "staff":
		data, err := userGw.GetStaffByUserAndOrganization(ctx, owner.OwnerID, organizationID)
		if err != nil {
			return &UserInfor{ID: owner.OwnerID}
		}
		return &UserInfor{
			ID:             data.ID,
			Name:           data.Name,
			Code:           data.Code,
			OrganizationID: data.OrganizationID,
		}

	default:
		return &UserInfor{ID: owner.OwnerID}
	}
}
