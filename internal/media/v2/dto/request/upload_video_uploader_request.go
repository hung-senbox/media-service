package request

import "mime/multipart"

type UploadVideoUploaderRequest struct {
	VideoUploaderID  string                `form:"video_uploader_id"`
	Title            string                `form:"title"`
	VideoFile        *multipart.FileHeader `form:"video_file"`
	ImagePreviewFile *multipart.FileHeader `form:"image_preview_file"`
	IsVisible        bool                  `form:"is_visible"`
}
