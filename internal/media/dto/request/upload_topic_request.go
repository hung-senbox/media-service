package request

import "mime/multipart"

type UploadTopicRequest struct {
	TopicID     string `form:"topic_id"`
	LanguageID  string `form:"language_id"`
	ParentID    string `form:"parent_id"`
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
	AudioOldKey  string                `form:"audio_old_key"`

	// video
	VideoFile    *multipart.FileHeader `form:"video_file"`
	VideoLinkUrl string                `form:"video_link_url"`
	VideoStart   string                `form:"video_start_time"`
	VideoEnd     string                `form:"video_end_time"`
	VideoOldKey  string                `form:"video_old_key"`

	// images
	FullBackgroundFile   *multipart.FileHeader `form:"full_background_file"`
	FullBackgroundLink   string                `form:"full_background_link_url"`
	FullBackgroundOldKey string                `form:"full_background_old_key"`

	ClearBackgroundFile   *multipart.FileHeader `form:"clear_background_file"`
	ClearBackgroundLink   string                `form:"clear_background_link_url"`
	ClearBackgroundOldKey string                `form:"clear_background_old_key"`

	ClipPartFile   *multipart.FileHeader `form:"clip_part_file"`
	ClipPartLink   string                `form:"clip_part_link_url"`
	ClipPartOldKey string                `form:"clip_part_old_key"`

	DrawingFile   *multipart.FileHeader `form:"drawing_file"`
	DrawingLink   string                `form:"drawing_link_url"`
	DrawingOldKey string                `form:"drawing_old_key"`

	IconFile   *multipart.FileHeader `form:"icon_file"`
	IconLink   string                `form:"icon_link_url"`
	IconOldKey string                `form:"icon_old_key"`

	BMFile   *multipart.FileHeader `form:"bm_file"`
	BMLink   string                `form:"bm_link_url"`
	BMOldKey string                `form:"bm_old_key"`
}
