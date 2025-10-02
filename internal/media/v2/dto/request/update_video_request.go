package request

import "mime/multipart"

type UpdateVideoRequest struct {
	TopicID      string                `json:"topic_id" binding:"required"`
	LanguageID   uint                  `json:"language_id" binding:"required"`
	VideoFile    *multipart.FileHeader `form:"video_file"`
	VideoLinkUrl string                `form:"video_link_url"`
	VideoStart   string                `form:"video_start_time"`
	VideoEnd     string                `form:"video_end_time"`
	VideoOldKey  string                `json:"video_old_key"`
}
