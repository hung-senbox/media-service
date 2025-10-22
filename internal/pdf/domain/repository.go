package domain

import (
	"context"
	"fmt"
	"media-service/internal/pdf/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PDFRepository interface {
	CreatePDF(ctx context.Context, pdf *model.StudentReportPDF) error
	GetPDFsByStudent(ctx context.Context, studentID string) ([]*model.StudentReportPDF, error)
}

type pdfRepository struct {
	PDFCollection *mongo.Collection
}

func NewPDFRepository(collection *mongo.Collection) PDFRepository {
	return &pdfRepository{
		PDFCollection: collection,
	}
}

func (p *pdfRepository) CreatePDF(ctx context.Context, pdf *model.StudentReportPDF) error {
	_, err := p.PDFCollection.InsertOne(ctx, pdf)
	return err
}

func (p *pdfRepository) GetPDFsByStudent(ctx context.Context, studentID string) ([]*model.StudentReportPDF, error) {

	var studentPdfs []*model.StudentReportPDF

	filter := bson.M{}
	fmt.Printf("studentID: %s\n", studentID)
	if studentID != "" {
		filter["student_id"] = studentID
	}

	cursor, err := p.PDFCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &studentPdfs)
	if err != nil {
		return nil, err
	}

	return studentPdfs, nil

}
