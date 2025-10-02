package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicImageConfig struct {
	ImageType string `json:"image_type" bson:"image_type,omitempty"`
	ImageKey  string `json:"image_key" bson:"image_key,omitempty"`
	LinkUrl   string `json:"link_url" bson:"link_url,omitempty"`
}

type TopicVideoConfig struct {
	VideoKey  string `json:"video_key" bson:"video_key,omitempty"`
	LinkUrl   string `json:"link_url" bson:"link_url,omitempty"`
	StartTime string `json:"start_time" bson:"start_time,omitempty"`
	EndTime   string `json:"end_time" bson:"end_time,omitempty"`
}

type TopicAudioConfig struct {
	AudioKey  string `json:"audio_key" bson:"audio_key,omitempty"`
	LinkUrl   string `json:"link_url" bson:"link_url,omitempty"`
	StartTime string `json:"start_time" bson:"start_time,omitempty"`
	EndTime   string `json:"end_time" bson:"end_time,omitempty"`
}

type TopicLanguageConfig struct {
	LanguageID  uint               `json:"language_id" bson:"language_id"`
	FileName    string             `json:"file_name" bson:"file_name"`
	Title       string             `json:"title" bson:"title"`
	Note        string             `json:"note" bson:"note"`
	Description string             `json:"description" bson:"description"`
	Images      []TopicImageConfig `json:"images" bson:"images"`
	Video       TopicVideoConfig   `json:"videos" bson:"video"`
	Audio       TopicAudioConfig   `json:"audios" bson:"audio"`
}

type Topic struct {
	ID             primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	ParentID       string                `json:"parent_id" bson:"parent_id,omitempty"`
	OrganizationID string                `json:"organization_id" bson:"organization_id,omitempty"`
	IsPublished    bool                  `json:"is_published" bson:"is_published,omitempty"`
	LanguageConfig []TopicLanguageConfig `json:"language_config" bson:"language_config,omitempty"`
	CreatedAt      time.Time             `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time             `json:"updated_at" bson:"updated_at,omitempty"`
}
