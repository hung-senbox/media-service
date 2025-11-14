package request

import "mime/multipart"

type UploadTopicRequest struct {
	TopicID     string `form:"topic_id"`
	LanguageID  uint   `form:"language_id"`
	IsPublished bool   `form:"is_published"`
	FileName    string `form:"file_name"`
	Title       string `form:"title"`
	Note        string `form:"note"`
	Description string `form:"description"`

	// audio
	AudioFile      *multipart.FileHeader `form:"audio_file"`
	AudioLinkUrl   string                `form:"audio_link_url"`
	AudioStart     string                `form:"audio_start_time"`
	AudioEnd       string                `form:"audio_end_time"`
	IsDeletedAudio bool                  `form:"is_deleted_audio"`

	// video
	VideoFile      *multipart.FileHeader `form:"video_file"`
	VideoLinkUrl   string                `form:"video_link_url"`
	VideoStart     string                `form:"video_start_time"`
	VideoEnd       string                `form:"video_end_time"`
	IsDeletedVideo bool                  `form:"is_deleted_video"`

	// images
	FullBackgroundFile      *multipart.FileHeader `form:"full_background_file"`
	FullBackgroundLink      string                `form:"full_background_link_url"`
	IsDeletedFullBackground bool                  `form:"is_deleted_full_background"`

	ClearBackgroundFile      *multipart.FileHeader `form:"clear_background_file"`
	ClearBackgroundLink      string                `form:"clear_background_link_url"`
	IsDeletedClearBackground bool                  `form:"is_deleted_clear_background"`

	ClipPartFile      *multipart.FileHeader `form:"clip_part_file"`
	ClipPartLink      string                `form:"clip_part_link_url"`
	IsDeletedClipPart bool                  `form:"is_deleted_clip_part"`

	DrawingFile      *multipart.FileHeader `form:"drawing_file"`
	DrawingLink      string                `form:"drawing_link_url"`
	IsDeletedDrawing bool                  `form:"is_deleted_drawing"`

	IconFile      *multipart.FileHeader `form:"icon_file"`
	IconLink      string                `form:"icon_link_url"`
	IsDeletedIcon bool                  `form:"is_deleted_icon"`

	BMFile      *multipart.FileHeader `form:"bm_file"`
	BMLink      string                `form:"bm_link_url"`
	IsDeletedBM bool                  `form:"is_deleted_bm"`

	SignLangFile      *multipart.FileHeader `form:"sign_lang_file"`
	SignLangLink      string                `form:"sign_lang_link_url"`
	IsDeletedSignLang bool                  `form:"is_deleted_sign_lang"`

	GifFile      *multipart.FileHeader `form:"gif_file"`
	GifLink      string                `form:"gif_link_url"`
	IsDeletedGif bool                  `form:"is_deleted_gif"`

	OrderFile      *multipart.FileHeader `form:"order_file"`
	OrderLink      string                `form:"order_link_url"`
	IsDeletedOrder bool                  `form:"is_deleted_order"`
}
