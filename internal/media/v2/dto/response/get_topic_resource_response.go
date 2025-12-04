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

type TopicResourceResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	PicID     string    `json:"pic_id"`
}

type GetTopicResourcesResponse4Web struct {
	ID        string                    `json:"id"`
	FileName  string                    `json:"file_name"`
	ImageUrl  string                    `json:"image_url"`
	CreatedAt time.Time                 `json:"created_at"`
	PicID     string                    `json:"pic_id"`
	Topic     *TopicResponse2Assign4Web `json:"topic"`
}

type GetTopicResourcesResponse4WebV2 struct {
	Date     string                     `json:"date"`
	Pictures []*TopicResourceResponseV2 `json:"pictures"`
}

type TopicResourceResponseV2 struct {
	ID        string                    `json:"id"`
	ImageKey  string                    `json:"image_key"`
	FileName  string                    `json:"file_name"`
	ImageUrl  string                    `json:"image_url"`
	CreatedAt time.Time                 `json:"created_at"`
	PicID     string                    `json:"pic_id"`
	TopicID   string                    `json:"topic_id"`
	Topic     *TopicResponse2Assign4Web `json:"topic"`
}

type GetTopicResourcesResponse4App struct {
	ID        string               `json:"id"`
	FileName  string               `json:"file_name"`
	ImageUrl  string               `json:"image_url"`
	CreatedAt time.Time            `json:"created_at"`
	PicID     string               `json:"pic_id"`
	Topic     GetTopicResponse4App `json:"topic"`
}

type GetTopicResourcesResponseByStudent4Web struct {
	Topic     *TopicResponse           `json:"topic"`
	Resources []*TopicResourceResponse `json:"resources"`
}
