package request

import "mime/multipart"

type UploadTopicRequest struct {
	LanguageID  string `form:"language_id"`
	IsPublished bool   `form:"is_published"`
	FileName    string `form:"file_name"`
	Title       string `form:"title"`
	Note        string `form:"note"`
	Description string `form:"description"`

	// audio
	AudioFile    *multipart.FileHeader `form:"audio_file"`
	AudioLinkUrl string                `form:"audio_link_url"`
	AudioStart   string                `form:"audio_start_time"`
	AudioEnd     string                `form:"audio_end_time"`

	// video
	VideoFile    *multipart.FileHeader `form:"video_file"`
	VideoLinkUrl string                `form:"video_link_url"`
	VideoStart   string                `form:"video_start_time"`
	VideoEnd     string                `form:"video_end_time"`

	// images
	FullBackgroundFile *multipart.FileHeader `form:"full_background_file"`
	FullBackgroundLink string                `form:"full_background_link_url"`

	ClearBackgroundFile *multipart.FileHeader `form:"clear_background_file"`
	ClearBackgroundLink string                `form:"clear_background_link_url"`

	ClipPartFile *multipart.FileHeader `form:"clip_part_file"`
	ClipPartLink string                `form:"clip_part_link_url"`

	DrawingFile *multipart.FileHeader `form:"drawing_file"`
	DrawingLink string                `form:"drawing_link_url"`

	IconFile *multipart.FileHeader `form:"icon_file"`
	IconLink string                `form:"icon_link_url"`

	BMFile *multipart.FileHeader `form:"bm_file"`
	BMLink string                `form:"bm_link_url"`
}
