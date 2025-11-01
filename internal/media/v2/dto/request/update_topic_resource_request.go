package request

import "mime/multipart"

type UpdateTopicResourceRequest struct {
	FileName string                `form:"file_name"`
	File     *multipart.FileHeader `form:"file"`
	TopicID  string                `form:"topic_id"`
}
