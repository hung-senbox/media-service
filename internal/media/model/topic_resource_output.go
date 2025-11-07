package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicResourceOutput struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	StudentID      string             `json:"student_id" bson:"student_id"`
	TopicResoureID string             `json:"topic_resoure_id" bson:"topic_resoure_id"`
	TopicResource  TopicResource      `json:"topic_resource" bson:"topic_resource"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
