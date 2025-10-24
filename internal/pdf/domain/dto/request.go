package dto

import "mime/multipart"

type CreatePDFRequest struct {
	FileName  string `json:"file_name" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
	Color     string `json:"color" binding:"required"`
}

type UpdatePDFRequest struct {
	FileName string                `form:"file_name"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
	Color    string                `form:"color"`
}
