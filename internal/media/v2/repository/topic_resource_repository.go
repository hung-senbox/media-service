package repository

import (
	"context"
	"fmt"
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
	GetTopicResouresByTopic(ctx context.Context, topicID string) ([]*model.TopicResource, error)
	GetTopicResouresByStudentID(ctx context.Context, studentID string) ([]*model.TopicResource, error)
	GetTopicResouresByTopicAndStudent(ctx context.Context, topicID, studentID string) ([]*model.TopicResource, error)
	SetOutputTopicResource(ctx context.Context, topicResourceID string) error
	GetTopicResourcesByStudent(ctx context.Context, studentID string) ([]*model.TopicResource, error)
	GetTopicResouresByStudentIDAndTopicID(ctx context.Context, studentID, topicID string) ([]*model.TopicResource, error)
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

func (r *topicResourceRepository) GetTopicResouresByTopic(ctx context.Context, topicID string) ([]*model.TopicResource, error) {
	var result []*model.TopicResource
	filter := bson.M{
		"topic_id": topicID,
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

func (r *topicResourceRepository) GetTopicResouresByStudentID(ctx context.Context, studentID string) ([]*model.TopicResource, error) {
	var result []*model.TopicResource
	filter := bson.M{
		"student_id": studentID,
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

func (r *topicResourceRepository) GetTopicResouresByTopicAndStudent(ctx context.Context, topicID, studentID string) ([]*model.TopicResource, error) {
	var result []*model.TopicResource
	filter := bson.M{
		"topic_id":   topicID,
		"student_id": studentID,
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

func (r *topicResourceRepository) SetOutputTopicResource(ctx context.Context, topicResourceID string) error {
	objectID, err := primitive.ObjectIDFromHex(topicResourceID)
	if err != nil {
		return fmt.Errorf("invalid topic resource id: %w", err)
	}
	_, err = r.topicResourceCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"is_output": true}})
	return err
}

func (r *topicResourceRepository) GetTopicResourcesByStudent(ctx context.Context, studentID string) ([]*model.TopicResource, error) {
	var result []*model.TopicResource
	filter := bson.M{
		"student_id": studentID,
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

func (r *topicResourceRepository) GetTopicResouresByStudentIDAndTopicID(ctx context.Context, studentID, topicID string) ([]*model.TopicResource, error) {
	var result []*model.TopicResource
	filter := bson.M{
		"student_id": studentID,
		"topic_id":   topicID,
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
