package repository

import (
	"context"
	"media-service/internal/media/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TopicRepository interface {
	UploadTopic(ctx context.Context, topic *model.Topic) error
	AddImageToTopic(ctx context.Context, topicID string, img model.TopicImageConfig) error
	AddVideoToTopic(ctx context.Context, topicID string, vid model.TopicVideoConfig) error
	AddAudioToTopic(ctx context.Context, topicID string, aud model.TopicAudioConfig) error
}

type topicRepository struct {
	topicCollection *mongo.Collection
}

func NewTopicRepository(topicCollection *mongo.Collection) TopicRepository {
	return &topicRepository{topicCollection: topicCollection}
}

func (r *topicRepository) UploadTopic(ctx context.Context, topic *model.Topic) error {
	_, err := r.topicCollection.InsertOne(ctx, topic)
	return err
}

func (r *topicRepository) AddImageToTopic(ctx context.Context, topicID string, img model.TopicImageConfig) error {
	filter := bson.M{"_id": topicID}
	update := bson.M{
		"$push": bson.M{"language_config.0.images": img},
	}
	_, err := r.topicCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *topicRepository) AddVideoToTopic(ctx context.Context, topicID string, vid model.TopicVideoConfig) error {
	filter := bson.M{"_id": topicID}
	update := bson.M{
		"$push": bson.M{"language_config.0.videos": vid},
	}
	_, err := r.topicCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *topicRepository) AddAudioToTopic(ctx context.Context, topicID string, aud model.TopicAudioConfig) error {
	filter := bson.M{"_id": topicID}
	update := bson.M{
		"$push": bson.M{"language_config.0.audios": aud},
	}
	_, err := r.topicCollection.UpdateOne(ctx, filter, update)
	return err
}
