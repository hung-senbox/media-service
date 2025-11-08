package dto

type UploadRequest struct {
	Folder    string  `form:"folder"`
	Mode      string  `form:"mode"`       // private | public
	MediaType *string `form:"media_type"` // optional override: image|audio|video|pdf
	File      any     `form:"file" binding:"required"`
}
