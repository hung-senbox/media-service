package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicImageConfig struct {
	ImageType   string `json:"image_type" bson:"image_type"`
	ImageKey    string `json:"image_key" bson:"image_key"`
	LinkUrl     string `json:"link_url" bson:"link_url"`
	UploadedUrl string `json:"uploaded_url" bson:"uploaded_url,omitempty"`
}

type TopicVideoConfig struct {
	VideoKey    string `json:"video_key" bson:"video_key"`
	LinkUrl     string `json:"link_url" bson:"link_url"`
	StartTime   string `json:"start_time" bson:"start_time"`
	EndTime     string `json:"end_time" bson:"end_time"`
	UploadedUrl string `json:"uploaded_url" bson:"uploaded_url,omitempty"`
}

type TopicAudioConfig struct {
	AudioKey    string `json:"audio_key" bson:"audio_key"`
	LinkUrl     string `json:"link_url" bson:"link_url"`
	StartTime   string `json:"start_time" bson:"start_time"`
	EndTime     string `json:"end_time" bson:"end_time"`
	UploadedUrl string `json:"uploaded_url" bson:"uploaded_url,omitempty"`
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
	ID             primitive.ObjectID    `json:"id" bson:"_id"`
	IsAllPic       bool                  `json:"is_all_pic" bson:"is_all_pic"`
	ParentID       string                `json:"parent_id" bson:"parent_id"`
	OrganizationID string                `json:"organization_id" bson:"organization_id"`
	IsPublished    bool                  `json:"is_published" bson:"is_published"`
	LanguageConfig []TopicLanguageConfig `json:"language_config" bson:"language_config"`
	CreatedAt      time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at" bson:"updated_at"`
}
