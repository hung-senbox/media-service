package domain

import (
	"context"
	"media-service/internal/pdf/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserResourceRepository interface {
	CreateResource(ctx context.Context, pdf *model.UserResource) error
	// GetPDFsByStudent(ctx context.Context, studentID string) ([]*model.StudentReportPDF, error)
	// GetPDFByID(ctx context.Context, id primitive.ObjectID) (*model.StudentReportPDF, error)
	// DeletePDFByID(ctx context.Context, id primitive.ObjectID) error
	// UpdatePDFByID(ctx context.Context, id primitive.ObjectID, pdf *model.StudentReportPDF) error
}

type userResourceRepository struct {
	UserResourceCollection *mongo.Collection
}

func NewUserResourceRepository(collection *mongo.Collection) UserResourceRepository {
	return &userResourceRepository{
		UserResourceCollection: collection,
	}
}

func (p *userResourceRepository) CreateResource(ctx context.Context, pdf *model.UserResource) error {
	_, err := p.UserResourceCollection.InsertOne(ctx, pdf)
	return err
}

// func (p *pdfRepository) GetPDFsByStudent(ctx context.Context, studentID string) ([]*model.StudentReportPDF, error) {

// 	var studentPdfs []*model.StudentReportPDF

// 	filter := bson.M{}
// 	fmt.Printf("studentID: %s\n", studentID)
// 	if studentID != "" {
// 		filter["student_id"] = studentID
// 	}

// 	cursor, err := p.PDFCollection.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = cursor.All(ctx, &studentPdfs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return studentPdfs, nil

// }

// func (p *pdfRepository) GetPDFByID(ctx context.Context, id primitive.ObjectID) (*model.StudentReportPDF, error) {

// 	var pdf *model.StudentReportPDF

// 	if err := p.PDFCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&pdf); err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return pdf, nil
// }

// func (p *pdfRepository) DeletePDFByID(ctx context.Context, id primitive.ObjectID) error {
// 	_, err := p.PDFCollection.DeleteOne(ctx, bson.M{"_id": id})
// 	return err
// }

// func (p *pdfRepository) UpdatePDFByID(ctx context.Context, id primitive.ObjectID, pdf *model.StudentReportPDF) error {
// 	_, err := p.PDFCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": pdf})
// 	return err
// }
