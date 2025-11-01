package request

import "mime/multipart"

type UploadFileRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	ImageName string                `form:"image_name"`
	Mode      string                `form:"mode" binding:"required"`
}

type UploadVideoRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	VideoName string                `form:"video_name"`
	Mode      string                `form:"mode" binding:"required"`
}

type UploadAudioRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	Folder    string                `form:"folder" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	AudioName string                `form:"audio_name"`
	Mode      string                `form:"mode" binding:"required"`
}
