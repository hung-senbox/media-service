package repository

import (
	"context"
	"media-service/internal/mediaasset/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MediaRepository interface {
	Create(ctx context.Context, media *model.MediaAsset) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.MediaAsset, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	UpdateFields(ctx context.Context, id primitive.ObjectID, update bson.M) error
	FindByKey(ctx context.Context, key string) (*model.MediaAsset, error)
}

type mediaRepository struct {
	col *mongo.Collection
}

func NewMediaRepository(col *mongo.Collection) MediaRepository {
	return &mediaRepository{col: col}
}

func (r *mediaRepository) Create(ctx context.Context, media *model.MediaAsset) (primitive.ObjectID, error) {
	res, err := r.col.InsertOne(ctx, media)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}
	return primitive.NilObjectID, mongo.ErrNilDocument
}

func (r *mediaRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*model.MediaAsset, error) {
	var out model.MediaAsset
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *mediaRepository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mediaRepository) UpdateFields(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *mediaRepository) FindByKey(ctx context.Context, key string) (*model.MediaAsset, error) {
	var out model.MediaAsset
	if err := r.col.FindOne(ctx, bson.M{"key": key}).Decode(&out); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}
