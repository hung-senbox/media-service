package repository

import (
	"context"
	"media-service/internal/department/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegionRepository interface {
	CreateRegion(ctx context.Context, section *model.Region) error
	GetByID(ctx context.Context, id string) (*model.Region, error)
	GetAllByOrgID(ctx context.Context, organizationID string) ([]*model.Region, error)
	UpdateRegion(ctx context.Context, region *model.Region) error
	UpdateRegionName(ctx context.Context, id string, name string) error
}

type regionRepository struct {
	collection *mongo.Collection
}

func NewRegionRepository(collection *mongo.Collection) RegionRepository {
	return &regionRepository{collection}
}

func (r *regionRepository) CreateRegion(ctx context.Context, region *model.Region) error {
	_, err := r.collection.InsertOne(ctx, region)
	return err
}

func (r *regionRepository) GetByID(ctx context.Context, id string) (*model.Region, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var region model.Region
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&region)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *regionRepository) GetAllByOrgID(ctx context.Context, organizationID string) ([]*model.Region, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}}) // ASC

	// filter theo organization_id
	filter := bson.M{"organization_id": organizationID}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var regions []*model.Region
	if err := cursor.All(ctx, &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

func (r *regionRepository) UpdateRegion(ctx context.Context, region *model.Region) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": region.ID}, bson.M{"$set": region})
	return err
}

func (r *regionRepository) UpdateRegionName(ctx context.Context, id string, name string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err // id không hợp lệ
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"name": name}},
	)
	return err
}
