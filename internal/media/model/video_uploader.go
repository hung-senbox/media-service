package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoUploader struct {
	ID             primitive.ObjectID            `bson:"_id" json:"id"`
	CreatedBy      string                        `bson:"created_by" json:"created_by"`
	IsVisible      bool                          `bson:"is_visible" json:"is_visible"`
	Title          string                        `bson:"title" json:"title"`
	LanguageConfig []VideoUploaderLanguageConfig `bson:"language_config" json:"language_config"`
	CreatedAt      time.Time                     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time                     `bson:"updated_at" json:"updated_at"`
}

type VideoUploaderLanguageConfig struct {
	ID                    primitive.ObjectID `bson:"_id" json:"id"`
	LanguageID            uint               `json:"language_id" bson:"language_id"`
	VideoKey              string             `json:"video_key" bson:"video_key"`
	VideoPublicUrl        string             `json:"video_public_url" bson:"video_public_url"`
	ImagePreviewKey       string             `json:"image_preview_key" bson:"image_preview_key"`
	ImagePreviewPublicUrl string             `json:"image_preview_public_url" bson:"image_preview_public_url"`
	Transcript            string             `json:"transcript" bson:"transcript"`
	Note                  string             `json:"note" bson:"note"`
}
