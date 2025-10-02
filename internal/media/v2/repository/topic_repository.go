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
	CreateTopic(ctx context.Context, topic *model.Topic) error
	SetLanguageConfig(ctx context.Context, topicID string, lang model.TopicLanguageConfig) error
	SetImage(ctx context.Context, topicID, languageID string, img model.TopicImageConfig) error
	SetAudio(ctx context.Context, topicID, languageID string, aud model.TopicAudioConfig) error
	SetVideo(ctx context.Context, topicID, languageID string, vid model.TopicVideoConfig) error
}

type topicRepository struct {
	topicCollection *mongo.Collection
}

func NewTopicRepository(topicCollection *mongo.Collection) TopicRepository {
	return &topicRepository{topicCollection: topicCollection}
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *model.Topic) error {
	_, err := r.topicCollection.InsertOne(ctx, topic)
	if err != nil {
		fmt.Printf("[CreateTopic] Insert failed: %v\n", err)
		return err
	}
	return nil
}

func (r *topicRepository) SetImage(ctx context.Context, topicID, languageID string, img model.TopicImageConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[AddImageToTopic] invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{"_id": objID, "language_config.language_id": languageID}
	update := bson.M{"$push": bson.M{"language_config.$.images": img}}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[AddImageToTopic] push failed: %w", err)
	}
	return nil
}

func (r *topicRepository) SetLanguageConfig(ctx context.Context, topicID string, lang model.TopicLanguageConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("invalid topicID: %w", err)
	}

	// --- Bước 1: Thử update nếu LanguageID đã tồn tại ---
	filterUpdate := bson.M{
		"_id":                         objID,
		"language_config.language_id": lang.LanguageID,
	}
	update := bson.M{
		"$set": bson.M{
			"language_config.$.file_name":   lang.FileName,
			"language_config.$.title":       lang.Title,
			"language_config.$.note":        lang.Note,
			"language_config.$.description": lang.Description,
		},
	}

	res, err := r.topicCollection.UpdateOne(ctx, filterUpdate, update)
	if err != nil {
		return fmt.Errorf("update language config failed: %w", err)
	}

	if res.MatchedCount > 0 {
		// Đã update thành công, không cần push
		return nil
	}

	// --- Bước 2: Nếu chưa tồn tại, push vào mảng ---
	filterPush := bson.M{"_id": objID}
	updatePush := bson.M{"$push": bson.M{"language_config": lang}}

	_, err = r.topicCollection.UpdateOne(ctx, filterPush, updatePush)
	if err != nil {
		return fmt.Errorf("insert language config failed: %w", err)
	}

	return nil
}

func (r *topicRepository) SetVideo(ctx context.Context, topicID, languageID string, vid model.TopicVideoConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[SetVideo4Topic] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}
	update := bson.M{
		"$set": bson.M{"language_config.$.video": vid},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[SetVideoForTopic] Update failed: %v\n", err)
		return err
	}

	return nil
}

func (r *topicRepository) SetAudio(ctx context.Context, topicID, languageID string, aud model.TopicAudioConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[SetAudio4Topic] invalid topicID=%s: %w", topicID, err)
	}

	// filter theo topicID và languageID
	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	// set audio mới, ghi đè audio cũ
	update := bson.M{
		"$set": bson.M{"language_config.$.audio": aud},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[SetAudioForTopic] update failed: %w", err)
	}

	return nil
}
