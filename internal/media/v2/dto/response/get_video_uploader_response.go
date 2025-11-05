package response

import "time"

type GetVideoUploaderResponse4Web struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	IsVisible       bool      `json:"is_visible"`
	VideoURL        string    `json:"video_url"`
	ImagePreviewURL string    `json:"image_preview_url"`
	CreatedByName   string    `json:"created_by_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
