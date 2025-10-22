package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentReportPDF struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	StudentID string             `json:"student_id" bson:"student_id"`
	PDFName   string             `json:"pdf_name" bson:"pdf_name"`
	Folder    string             `json:"folder" bson:"folder"`
	PDFKey    string             `json:"pdf_key" bson:"pdf_key"`
	CreatedBy string             `json:"created_by" bson:"created_by"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
