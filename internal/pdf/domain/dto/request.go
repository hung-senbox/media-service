package dto

import (
	"media-service/internal/pdf/model"
	"mime/multipart"
)

type CreateResource struct {
	UploaderID   *model.Owner `json:"uploader_id" bson:"uploader_id"`
	TargetID     *model.Owner `json:"target_id" bson:"target_id"`
	ResourceType string       `json:"resource_type" bson:"resource_type"`
	Color        string       `json:"color" bson:"color"`
	Folder       string       `json:"folder" bson:"folder"`
	SignatureKey string       `json:"signature_key" bson:"signature_key"`
}
type UpdatePDFRequest struct {
	FileName string                `form:"file_name"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
	Color    string                `form:"color"`
}
