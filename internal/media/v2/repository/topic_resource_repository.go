package repository

import (
	"context"
	"media-service/internal/media/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TopicResourceRepository interface {
	CreateTopicResource(ctx context.Context, topicResource *model.TopicResource) error
	GetTopicResources(ctx context.Context, topicID, studentID string) ([]*model.TopicResource, error)
	GetTopicResource(ctx context.Context, topicResourceID primitive.ObjectID) (*model.TopicResource, error)
	UpdateTopicResource(ctx context.Context, topicResourceID primitive.ObjectID, topicResource *model.TopicResource) error
	DeleteTopicResource(ctx context.Context, topicResourceID primitive.ObjectID) error
	GetTopicResouresByOrganizationAndTopicID(ctx context.Context, topicID, organizationID string) ([]*model.TopicResource, error)
}

type topicResourceRepository struct {
	topicResourceCollection *mongo.Collection
}

func NewTopicResourceRepository(topicResourceCollection *mongo.Collection) TopicResourceRepository {
	return &topicResourceRepository{topicResourceCollection: topicResourceCollection}
}

func (r *topicResourceRepository) CreateTopicResource(ctx context.Context, topicResource *model.TopicResource) error {
	_, err := r.topicResourceCollection.InsertOne(ctx, topicResource)
	return err
}

func (r *topicResourceRepository) GetTopicResources(ctx context.Context, topicID, studentID string) ([]*model.TopicResource, error) {

	var result []*model.TopicResource

	filter := bson.M{}

	if topicID != "" {
		filter["topic_id"] = topicID
	}

	if studentID != "" {
		filter["student_id"] = studentID
	}

	cursor, err := r.topicResourceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *topicResourceRepository) GetTopicResource(ctx context.Context, topicResourceID primitive.ObjectID) (*model.TopicResource, error) {
	var result model.TopicResource
	err := r.topicResourceCollection.FindOne(ctx, bson.M{"_id": topicResourceID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *topicResourceRepository) UpdateTopicResource(ctx context.Context, topicResourceID primitive.ObjectID, topicResource *model.TopicResource) error {
	_, err := r.topicResourceCollection.UpdateOne(ctx, bson.M{"_id": topicResourceID}, bson.M{"$set": topicResource})
	return err
}

func (r *topicResourceRepository) DeleteTopicResource(ctx context.Context, topicResourceID primitive.ObjectID) error {
	_, err := r.topicResourceCollection.DeleteOne(ctx, bson.M{"_id": topicResourceID})
	return err
}

func (r *topicResourceRepository) GetTopicResouresByOrganizationAndTopicID(ctx context.Context, topicID, organizationID string) ([]*model.TopicResource, error) {
	var result []*model.TopicResource
	filter := bson.M{
		"topic_id":        topicID,
		"organization_id": organizationID,
	}
	cursor, err := r.topicResourceCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
