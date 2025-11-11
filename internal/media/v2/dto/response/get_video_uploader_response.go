package response

import "time"

type GetVideoUploaderResponse4Web struct {
	ID              string    `json:"id"`
	LanguageID      uint      `json:"language_id"`
	IsVisible       bool      `json:"is_visible"`
	CreatedByName   string    `json:"created_by_name"`
	Title           string    `json:"title"`
	VideoUrl        string    `json:"video_url"`
	ImagePreviewUrl string    `json:"image_preview_url"`
	CreatedAt       time.Time `json:"created_at"`
}
