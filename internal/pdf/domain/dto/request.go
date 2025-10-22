package dto

import "mime/multipart"

type CreatePDFRequest struct {
	File      *multipart.FileHeader `form:"file" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	StudentID string                `form:"student_id" binding:"required"`
}
