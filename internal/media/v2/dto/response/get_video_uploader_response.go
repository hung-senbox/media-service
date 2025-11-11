package response

type GetVideoUploaderResponse4Web struct {
	ID            string                 `json:"id"`
	IsVisible     bool                   `json:"is_visible"`
	CreatedByName string                 `json:"created_by_name"`
	MessageLangs  []MessageLanguageVideo `json:"message_languages"`
}

type MessageLanguageVideo struct {
	LanguageID int                   `json:"language_id"`
	Contents   VideoLanguageContents `json:"contents"`
}

type VideoLanguageContents struct {
	Title           string `json:"title"`
	VideoURL        string `json:"video_url"`
	ImagePreviewURL string `json:"image_preview_url"`
}
