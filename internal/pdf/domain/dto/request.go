package dto

import (
	"media-service/internal/pdf/model"
	"mime/multipart"
)

type CreateResourceRequest struct {
	Role           string       `json:"role" bson:"role"`
	OrganizationID string       `json:"organization_id" bson:"organization_id"`
	UploaderID     *model.Owner `json:"uploader_id" bson:"uploader_id"`
	TargetID       *model.Owner `json:"target_id" bson:"target_id"`
	Type           string       `json:"type" bson:"type"`
	ResourceType   string       `json:"resource_type" bson:"resource_type"`
	Color          string       `json:"color" bson:"color"`
	Folder         string       `json:"folder" bson:"folder"`
	SignatureKey   string       `json:"signature_key" bson:"signature_key"`
}
type UpdateResourceRequest struct {
	FileName     *string               `form:"file_name"`
	File         *multipart.FileHeader `form:"file"`
	ResourceType string                `form:"resource_type"`
	Url          *string               `form:"url"`
}

type UploadSignatureRequest struct {
	SignatureKey string `json:"signature_key" bson:"signature_key"`
}
