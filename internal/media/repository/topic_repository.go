package repository

import (
	"context"
	"fmt"
	"media-service/internal/media/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err != nil {
		fmt.Printf("[UploadTopic] InsertOne failed: %v\n", err)
	}
	return err
}

func (r *topicRepository) AddImageToTopic(ctx context.Context, topicID string, img model.TopicImageConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[AddImageToTopic] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{"language_config.0.images": img},
	}

	res, err := r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[AddImageToTopic] UpdateOne failed for topicID=%s: %v\n", topicID, err)
		return err
	}

	if res.MatchedCount == 0 {
		fmt.Printf("[AddImageToTopic] no document matched for topicID=%s\n", topicID)
	} else if res.ModifiedCount == 0 {
		fmt.Printf("[AddImageToTopic] document matched but no modification for topicID=%s\n", topicID)
	}

	return nil
}

func (r *topicRepository) AddVideoToTopic(ctx context.Context, topicID string, vid model.TopicVideoConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[AddVideoToTopic] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{"language_config.0.videos": vid},
	}

	res, err := r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[AddVideoToTopic] UpdateOne failed for topicID=%s: %v\n", topicID, err)
		return err
	}

	if res.MatchedCount == 0 {
		fmt.Printf("[AddVideoToTopic] no document matched for topicID=%s\n", topicID)
	} else if res.ModifiedCount == 0 {
		fmt.Printf("[AddVideoToTopic] document matched but no modification for topicID=%s\n", topicID)
	}

	return nil
}

func (r *topicRepository) AddAudioToTopic(ctx context.Context, topicID string, aud model.TopicAudioConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[AddAudioToTopic] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{"language_config.0.audios": aud},
	}

	res, err := r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[AddAudioToTopic] UpdateOne failed for topicID=%s: %v\n", topicID, err)
		return err
	}

	if res.MatchedCount == 0 {
		fmt.Printf("[AddAudioToTopic] no document matched for topicID=%s\n", topicID)
	} else if res.ModifiedCount == 0 {
		fmt.Printf("[AddAudioToTopic] document matched but no modification for topicID=%s\n", topicID)
	}

	return nil
}
