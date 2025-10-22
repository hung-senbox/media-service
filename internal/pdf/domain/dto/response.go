package dto

import (
	"time"
)

type StudentReportPDFResponse struct {
	ID        string    `json:"id" bson:"_id"`
	StudentID string    `json:"student_id" bson:"student_id"`
	PDFName   string    `json:"pdf_name" bson:"pdf_name"`
	Folder    string    `json:"folder" bson:"folder"`
	PDFUrl    *string   `json:"pdf_url" bson:"pdf_url"`
	CreatedBy string    `json:"created_by" bson:"created_by"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
