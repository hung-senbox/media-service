package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TopicResource struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	TopicID   string             `json:"topic_id" bson:"topic_id"`
	StudentID string             `json:"student_id" bson:"student_id"`
	FileName  string             `json:"file_name" bson:"file_name"`
	ImageKey  string             `json:"image_key" bson:"image_key"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
