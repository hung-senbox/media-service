package response

import (
	"media-service/internal/media/model"
	"time"
)

type GetVideoUploaderResponse4Web struct {
	ID               string    `json:"id"`
	LanguageID       uint      `json:"language_id"`
	LanguageConfigID string    `json:"language_config_id"`
	IsVisible        bool      `json:"is_visible"`
	CreatedByName    string    `json:"created_by_name"`
	Title            string    `json:"title"`
	WikiCode         string    `json:"wiki_code"`
	VideoUrl         string    `json:"video_url"`
	ImagePreviewUrl  string    `json:"image_preview_url"`
	Note             string    `json:"note"`
	Transcript       string    `json:"transcript"`
	CreatedAt        time.Time `json:"created_at"`
}

type GetDetailVideo4WebResponse struct {
	ID            string                            `json:"id"`
	IsVisible     bool                              `json:"is_visible"`
	Title         string                            `json:"title"`
	WikiCode      string                            `json:"wiki_code"`
	CreatedByName string                            `json:"created_by_name"`
	MessageLangs  []DetailVideoMessageLanguageEntry `json:"message_languages"`
	CreatedAt     time.Time                         `json:"created_at"`
}

type DetailVideoMessageLanguageEntry struct {
	LanguageID int                         `json:"language_id"`
	Contents   DetailVideoLanguageContents `json:"contents"`
}

type DetailVideoLanguageContents struct {
	Note            string `json:"note"`
	Transcript      string `json:"transcript"`
	VideoUrl        string `json:"video_url"`
	ImagePreviewUrl string `json:"image_preview_url"`
}

type GetVideosByWikiCode4WebResponse struct {
	ID             string                              `json:"id"`
	Title          string                              `json:"title"`
	WikiCode       string                              `json:"wiki_code"`
	LanguageConfig []model.VideoUploaderLanguageConfig `json:"language_config"`
	CreatedAt      time.Time                           `json:"created_at"`
}
