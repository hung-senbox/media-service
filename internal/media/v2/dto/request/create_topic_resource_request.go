package request

import "mime/multipart"

type CreateTopicResourceRequest struct {
	TopicID   string                `form:"topic_id" binding:"required"`
	StudentID string                `form:"student_id" binding:"required"`
	FileName  string                `form:"file_name" binding:"required"`
	File      *multipart.FileHeader `form:"file" binding:"required"`
}
