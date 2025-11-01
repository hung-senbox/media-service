package repository

import (
	"context"
	"errors"
	"fmt"
	"media-service/internal/media/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TopicRepository interface {
	CreateTopic(ctx context.Context, topic *model.Topic) (*model.Topic, error)
	UpdateTopic(ctx context.Context, topic *model.Topic) (*model.Topic, error)
	SetLanguageConfig(ctx context.Context, topicID string, lang model.TopicLanguageConfig) error
	SetImage(ctx context.Context, topicID string, languageID uint, img model.TopicImageConfig) error
	SetAudio(ctx context.Context, topicID string, languageID uint, aud model.TopicAudioConfig) error
	SetVideo(ctx context.Context, topicID string, languageID uint, vid model.TopicVideoConfig) error
	GetAllTopicByOrganizationID(ctx context.Context, orgID string) ([]model.Topic, error)
	GetByID(ctx context.Context, id string) (*model.Topic, error)
	GetAllTopicByOrganizationIDAndIsPublished(ctx context.Context, orgID string) ([]model.Topic, error)
}

type topicRepository struct {
	topicCollection *mongo.Collection
}

func NewTopicRepository(topicCollection *mongo.Collection) TopicRepository {
	return &topicRepository{topicCollection: topicCollection}
}

func (r *topicRepository) CreateTopic(ctx context.Context, topic *model.Topic) (*model.Topic, error) {
	result, err := r.topicCollection.InsertOne(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("insert topic failed: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		topic.ID = oid
	}

	return topic, nil
}

func (r *topicRepository) SetImage(ctx context.Context, topicID string, languageID uint, img model.TopicImageConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[SetImageForTopic] invalid topicID=%s: %w", topicID, err)
	}

	// 1️⃣ Thử update nếu image_type đã tồn tại
	filter := bson.M{
		"_id": objID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$[lang].images.$[imgElem].image_key":  img.ImageKey,
			"language_config.$[lang].images.$[imgElem].link_url":   img.LinkUrl,
			"language_config.$[lang].images.$[imgElem].image_type": img.ImageType,
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"lang.language_id": languageID},
			bson.M{"imgElem.image_type": img.ImageType},
		},
	})

	res, err := r.topicCollection.UpdateOne(ctx, filter, update, arrayFilters)
	if err != nil {
		return fmt.Errorf("[SetImageForTopic] update failed: %w", err)
	}

	// 2️⃣ Nếu chưa có image_type → push thêm vào mảng
	if res.ModifiedCount == 0 {
		pushUpdate := bson.M{
			"$push": bson.M{
				"language_config.$[lang].images": bson.M{
					"image_type": img.ImageType,
					"image_key":  img.ImageKey,
					"link_url":   img.LinkUrl,
				},
			},
		}

		pushFilter := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"lang.language_id": languageID},
			},
		})

		_, err = r.topicCollection.UpdateOne(ctx, filter, pushUpdate, pushFilter)
		if err != nil {
			return fmt.Errorf("[SetImageForTopic] push failed: %w", err)
		}
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

func (r *topicRepository) SetVideo(ctx context.Context, topicID string, languageID uint, vid model.TopicVideoConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[SetVideo4Topic] invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$.video": bson.M{
				"video_key":  vid.VideoKey,
				"link_url":   vid.LinkUrl,
				"start_time": vid.StartTime,
				"end_time":   vid.EndTime,
			},
		},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[SetVideoForTopic] update failed: %w", err)
	}

	return nil
}

func (r *topicRepository) SetAudio(ctx context.Context, topicID string, languageID uint, aud model.TopicAudioConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[SetAudio4Topic] invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$.audio": bson.M{
				"audio_key":  aud.AudioKey,
				"link_url":   aud.LinkUrl,
				"start_time": aud.StartTime,
				"end_time":   aud.EndTime,
			},
		},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[SetAudioForTopic] update failed: %w", err)
	}

	return nil
}

func (r *topicRepository) GetAllTopicByOrganizationID(ctx context.Context, orgID string) ([]model.Topic, error) {
	var topics []model.Topic
	filter := bson.M{
		"organization_id": orgID,
		"$or": []bson.M{
			{"parent_id": ""},
			{"parent_id": bson.M{"$exists": false}},
			{"parent_id": nil},
		},
	}

	cursor, err := r.topicCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

func (r *topicRepository) UpdateTopic(ctx context.Context, topic *model.Topic) (*model.Topic, error) {
	filter := bson.M{"_id": topic.ID}

	updateFields := bson.M{
		"is_published": topic.IsPublished,
		"updated_at":   time.Now(),
	}

	// xử lý LanguageConfig
	if len(topic.LanguageConfig) > 0 {
		var langs []bson.M
		for _, lc := range topic.LanguageConfig {
			langUpdate := bson.M{
				"language_id": lc.LanguageID,
				"file_name":   lc.FileName,
				"title":       lc.Title,
				"note":        lc.Note,
				"description": lc.Description,
			}
			langs = append(langs, langUpdate)
		}
		updateFields["language_config"] = langs
	}

	update := bson.M{"$set": updateFields}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedTopic model.Topic
	err := r.topicCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedTopic)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("topic not found")
		}
		return nil, fmt.Errorf("update topic failed: %w", err)
	}

	return &updatedTopic, nil
}

func (r *topicRepository) GetByID(ctx context.Context, topicID string) (*model.Topic, error) {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return nil, fmt.Errorf("invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{"_id": objID}
	var topic model.Topic
	err = r.topicCollection.FindOne(ctx, filter).Decode(&topic)
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *topicRepository) GetAllTopicByOrganizationIDAndIsPublished(ctx context.Context, orgID string) ([]model.Topic, error) {
	var topics []model.Topic
	filter := bson.M{
		"organization_id": orgID,
		"is_published":    true,
		"$or": []bson.M{
			{"parent_id": ""},
			{"parent_id": bson.M{"$exists": false}},
			{"parent_id": nil},
		},
	}

	cursor, err := r.topicCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}
