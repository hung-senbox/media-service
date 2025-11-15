package dto

import (
	"time"
)

type ResourceResponse struct {
	ID             string     `json:"id" bson:"_id"`
	OrganizationID string     `json:"organization_id" bson:"organization_id"`
	UploaderInfor  *UserInfor `json:"uploader_infor" bson:"uploader_infor"`
	TargetInfor    *UserInfor `json:"target_infor" bson:"target_infor"`
	ResourceType   string     `json:"resource_type" bson:"resource_type"`
	FileName       *string    `json:"file_name" bson:"file_name"`
	Folder         string     `json:"folder" bson:"folder"`
	Color          string     `json:"color" bson:"color"`
	Status         int        `json:"status" bson:"status"`               // 0 waiting, 1 viewed, 2 rejected, 3 signed, 4 need to helps
	IsDownloaded   int        `json:"is_downloaded" bson:"is_downloaded"` // 0 not downloaded, 1 downloaded
	SignatureUrl   *string    `json:"signature_url" bson:"signature_url"`
	URL            *string    `json:"url" bson:"url"`
	PDFUrl         *string    `json:"pdf_url" bson:"pdf_url"`
	CreatedBy      string     `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
}
type UserInfor struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
}

type GroupedResourceResponse struct {
	SelfResources    []*ResourceResponse `json:"self_resources" bson:"self_resources"`
	RelatedResources []*ResourceResponse `json:"related_resources" bson:"related_resources"`
	StudentResources []*ResourceResponse `json:"student_resources,omitempty" bson:"student_resources,omitempty"`
}
