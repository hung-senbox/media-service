package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoUploader struct {
	ID             primitive.ObjectID    `bson:"_id" json:"id"`
	CreatedBy      string                `bson:"created_by" json:"created_by"`
	IsVisible      bool                  `bson:"is_visible" json:"is_visible"`
	LanguageConfig []VideoLanguageConfig `bson:"language_config" json:"language_config"`
	CreatedAt      time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time             `bson:"updated_at" json:"updated_at"`
}

type VideoLanguageConfig struct {
	LanguageID            uint   `json:"language_id" bson:"language_id"`
	Title                 string `json:"title" bson:"title"`
	VideoKey              string `bson:"video_key" json:"video_key"`
	VideoPublicUrl        string `bson:"video_public_url" json:"video_public_url"`
	ImagePreviewKey       string `bson:"image_preview_key" json:"image_preview_key"`
	ImagePreviewPublicUrl string `bson:"image_preview_public_url" json:"image_preview_public_url"`
}
