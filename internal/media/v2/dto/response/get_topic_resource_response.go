package response

import (
	"media-service/internal/gateway/dto/response"
	"time"
)

type GetTopicResourceResponse struct {
	ID        string                    `json:"id"`
	Topic     *TopicResponse2Assign4Web `json:"topic"`
	Student   *response.StudentResponse `json:"student"`
	FileName  string                    `json:"file_name"`
	ImageUrl  string                    `json:"image_url"`
	CreatedBy *response.TeacherResponse `json:"created_by"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

type GetTopicResourcesResponse4Web struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	PicID     string    `json:"pic_id"`
}
