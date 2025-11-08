package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaType string

const (
	MediaImage MediaType = "image"
	MediaAudio MediaType = "audio"
	MediaVideo MediaType = "video"
	MediaPDF   MediaType = "pdf"
	MediaOther MediaType = "other"
)

type MediaAsset struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type        MediaType          `bson:"type" json:"type"`
	Key         string             `bson:"key" json:"key"`
	FileName    string             `bson:"file_name" json:"file_name"`
	ContentType string             `bson:"content_type" json:"content_type"`
	Size        int64              `bson:"size" json:"size"`
	Folder      string             `bson:"folder" json:"folder"`
	Mode        string             `bson:"mode" json:"mode"` // public | private
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy   *string            `bson:"created_by,omitempty" json:"created_by,omitempty"`
}
