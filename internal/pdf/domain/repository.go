package domain

import (
	"context"
	"media-service/internal/pdf/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserResourceRepository interface {
	CreateResource(ctx context.Context, pdf *model.UserResource) error
	GetResourceByID(ctx context.Context, id primitive.ObjectID) (*model.UserResource, error)
	GetSelfResources(ctx context.Context, ownerID string) ([]*model.UserResource, error)
	GetRelatedResources(ctx context.Context, ownerID string) ([]*model.UserResource, error)
	GetStudentResources(ctx context.Context, studentIDs []string) ([]*model.UserResource, error)
	UpdateResourceByID(ctx context.Context, id primitive.ObjectID, pdf *model.UserResource) error
	UpdateResourceFields(ctx context.Context, id primitive.ObjectID, updateFields bson.M) error
	DeleteResourceByID(ctx context.Context, id primitive.ObjectID) error
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

func (p *userResourceRepository) GetResourceByID(ctx context.Context, id primitive.ObjectID) (*model.UserResource, error) {

	var pdf *model.UserResource

	if err := p.UserResourceCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&pdf); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return pdf, nil

}

func (p *userResourceRepository) GetSelfResources(ctx context.Context, ownerID string) ([]*model.UserResource, error) {

	var resources []*model.UserResource

	filter := bson.M{
		"uploader_id.owner_id": ownerID,
		"target_id":            nil,
	}

	cursor, err := p.UserResourceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (p *userResourceRepository) GetRelatedResources(ctx context.Context, ownerID string) ([]*model.UserResource, error) {

	var resources []*model.UserResource

	filter := bson.M{
		"uploader_id.owner_id": ownerID,
		"target_id.owner_id":   bson.M{"$exists": true},
	}

	cursor, err := p.UserResourceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &resources)
	if err != nil {
		return nil, err
	}

	return resources, nil

}

func (p *userResourceRepository) GetStudentResources(ctx context.Context, studentIDs []string) ([]*model.UserResource, error) {
	var resources []*model.UserResource
	filter := bson.M{
		"target_id.owner_id": bson.M{"$in": studentIDs},
	}
	cursor, err := p.UserResourceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (p *userResourceRepository) UpdateResourceByID(ctx context.Context, id primitive.ObjectID, pdf *model.UserResource) error {
	_, err := p.UserResourceCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": pdf})
	return err
}

func (p *userResourceRepository) UpdateResourceFields(ctx context.Context, id primitive.ObjectID, updateFields bson.M) error {
	_, err := p.UserResourceCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateFields})
	return err
}

func (p *userResourceRepository) DeleteResourceByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := p.UserResourceCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
