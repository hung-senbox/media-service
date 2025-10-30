package request

import "mime/multipart"

type UploadImageRequest struct {
	LanguageID          uint                  `json:"language_id" binding:"required"`
	FullBackgroundFile  *multipart.FileHeader `form:"full_background_file"`
	FullBackgroundLink  string                `form:"full_background_link_url"`
	ClearBackgroundFile *multipart.FileHeader `form:"clear_background_file"`
	ClearBackgroundLink string                `form:"clear_background_link_url"`
	ClipPartFile        *multipart.FileHeader `form:"clip_part_file"`
	ClipPartLink        string                `form:"clip_part_link_url"`
	DrawingFile         *multipart.FileHeader `form:"drawing_file"`
	DrawingLink         string                `form:"drawing_link_url"`
	IconFile            *multipart.FileHeader `form:"icon_file"`
	IconLink            string                `form:"icon_link_url"`
	BMFile              *multipart.FileHeader `form:"bm_file"`
	BMLink              string                `form:"bm_link_url"`
}
