package request

import "mime/multipart"

type UploadAudioRequest struct {
	TopicID      string                `form:"topic_id" binding:"required"`
	LanguageID   uint                  `form:"language_id" binding:"required"`
	AudioFile    *multipart.FileHeader `form:"audio_file"`
	AudioLinkUrl string                `form:"audio_link_url"`
	AudioStart   string                `form:"audio_start_time"`
	AudioEnd     string                `form:"audio_end_time"`
}
