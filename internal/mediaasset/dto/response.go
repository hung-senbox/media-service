package dto

type UploadResponse struct {
	ID  string  `json:"id"`
	Key string  `json:"key"`
	URL *string `json:"url,omitempty"`
}
