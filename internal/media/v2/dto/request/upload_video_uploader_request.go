package request

import "mime/multipart"

type UploadVideoUploaderRequest struct {
	VideoFolderID         string                `form:"video_folder_id"`
	Title                 string                `form:"title"`
	WikiCode              string                `form:"wiki_code"`
	VideoFile             *multipart.FileHeader `form:"video_file"`
	ImagePreviewFile      *multipart.FileHeader `form:"image_preview_file"`
	IsVisible             bool                  `form:"is_visible"`
	LanguageID            uint                  `form:"language_id"`
	IsDeletedVideo        bool                  `form:"is_deleted_video"`
	IsDeletedImagePreview bool                  `form:"is_deleted_image"`
	Note                  string                `form:"note"`
	Transcript            string                `form:"transcript"`
}
