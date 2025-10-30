package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResource struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Organization string             `json:"organization" bson:"organization"`
	Type         string             `json:"type" bson:"type"`
	UploaderID   *Owner             `json:"uploader_id" bson:"uploader_id"`
	TargetID     *Owner             `json:"target_id" bson:"target_id"`
	ResourceType string             `json:"resource_type" bson:"resource_type"`
	FileName     *string            `json:"file_name" bson:"file_name"`
	Folder       string             `json:"folder" bson:"folder"`
	Color        string             `json:"color" bson:"color"`
	SignatureKey *string            `json:"signature_key" bson:"signature_key"`
	URL          *string            `json:"url" bson:"url"`
	PDFKey       *string            `json:"pdf_key" bson:"pdf_key"`
	CreatedBy    string             `json:"created_by" bson:"created_by"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type Owner struct {
	OwnerID   string `json:"owner_id" bson:"owner_id"`
	OwnerRole string `json:"owner_role" bson:"owner_role"`
}
