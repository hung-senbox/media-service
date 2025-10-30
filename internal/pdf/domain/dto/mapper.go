package dto

import (
	"context"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/pdf/model"
)

func ToResourceResponses(
	ctx context.Context,
	organizationID string,
	data []*model.UserResource,
	userGw gateway.UserGateway,
	fileGw gateway.FileGateway,
) []*ResourceResponse {
	responses := make([]*ResourceResponse, 0, len(data))
	for _, r := range data {
		resp := &ResourceResponse{
			ID:             r.ID.Hex(),
			OrganizationID: r.Organization,
			ResourceType:   r.ResourceType,
			Folder:         r.Folder,
			Color:          r.Color,
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
			url, err := fileGw.GetImageUrl(ctx, gw_request.GetFileUrlRequest{
				Key:  *r.SignatureKey,
				Mode: "private",
			})
			if err == nil {
				resp.SignatureUrl = url
			}
		}

		if r.PDFKey != nil && *r.PDFKey != "" {
			url, err := fileGw.GetPDFUrl(ctx, gw_request.GetFileUrlRequest{
				Key:  *r.PDFKey,
				Mode: "private",
			})
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
