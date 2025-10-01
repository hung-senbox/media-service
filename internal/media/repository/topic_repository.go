package repository

import (
	"context"
	"fmt"
	"media-service/internal/media/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TopicRepository interface {
	UploadTopic(ctx context.Context, topic *model.Topic) error
	AddImageToTopic(ctx context.Context, topicID, languageID, imgType, oldKey string, img model.TopicImageConfig) error
	AddVideoToTopic(ctx context.Context, topicID, languageID, oldKey string, vid model.TopicVideoConfig) error
	AddAudioToTopic(ctx context.Context, topicID, languageID, oldKey string, aud model.TopicAudioConfig) error
	AddLanguageConfigToTopic(ctx context.Context, topicID string, lang model.TopicLanguageConfig) error
}

type topicRepository struct {
	topicCollection *mongo.Collection
}

func NewTopicRepository(topicCollection *mongo.Collection) TopicRepository {
	return &topicRepository{topicCollection: topicCollection}
}

func (r *topicRepository) UploadTopic(ctx context.Context, topic *model.Topic) error {
	filter := bson.M{"_id": topic.ID}
	update := bson.M{
		"$set": topic,
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.topicCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		fmt.Printf("[UploadTopic] Upsert failed: %v\n", err)
	}
	return err
}

func (r *topicRepository) AddImageToTopic(ctx context.Context, topicID, languageID, imgType, oldKey string, img model.TopicImageConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[UpsertImage] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	img.ImageType = imgType // đảm bảo lưu đúng type

	if oldKey != "" {
		// Step 1: Update nếu oldKey tồn tại trong language_config tương ứng
		filter := bson.M{"_id": objID, "language_config.language_id": languageID}
		update := bson.M{"$set": bson.M{"language_config.$[lang].images.$[img]": img}}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"lang.language_id": languageID},
				bson.M{"img.image_key": oldKey, "img.image_type": imgType},
			},
		}
		opts := options.Update().SetArrayFilters(arrayFilters)

		res, err := r.topicCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			fmt.Printf("[UpsertImage] UpdateOne failed: %v\n", err)
			return err
		}

		if res.MatchedCount > 0 {
			// Đã update thành công
			return nil
		}
	}

	// Step 2: Nếu oldKey rỗng hoặc không tìm thấy -> push mới
	filter := bson.M{"_id": objID, "language_config.language_id": languageID}
	update := bson.M{"$push": bson.M{"language_config.$.images": img}}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[UpsertImage] Push new image failed: %v\n", err)
		return err
	}

	return nil
}

func (r *topicRepository) AddVideoToTopic(ctx context.Context, topicID, languageID, oldKey string, vid model.TopicVideoConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[UpsertVideo] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	if oldKey != "" {
		// --- Step 1: Update nếu oldKey tồn tại trong language_config tương ứng ---
		filter := bson.M{"_id": objID, "language_config.language_id": languageID}
		update := bson.M{"$set": bson.M{"language_config.$[lang].videos.$[vid]": vid}}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"lang.language_id": languageID},
				bson.M{"vid.video_key": oldKey},
			},
		}
		opts := options.Update().SetArrayFilters(arrayFilters)

		res, err := r.topicCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			fmt.Printf("[UpsertVideo] UpdateOne failed: %v\n", err)
			return err
		}

		if res.MatchedCount > 0 {
			// Đã update thành công
			return nil
		}
	}

	// --- Step 2: Nếu oldKey rỗng hoặc không tìm thấy -> push mới ---
	filter := bson.M{"_id": objID, "language_config.language_id": languageID}
	update := bson.M{"$push": bson.M{"language_config.$.videos": vid}}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[UpsertVideo] Push new video failed: %v\n", err)
		return err
	}

	return nil
}

func (r *topicRepository) AddAudioToTopic(ctx context.Context, topicID, languageID, oldKey string, aud model.TopicAudioConfig) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		fmt.Printf("[UpsertAudio] invalid topicID=%s: %v\n", topicID, err)
		return err
	}

	if oldKey != "" {
		// --- Step 1: Update nếu oldKey tồn tại trong language_config tương ứng ---
		filter := bson.M{"_id": objID, "language_config.language_id": languageID}
		update := bson.M{"$set": bson.M{"language_config.$[lang].audios.$[aud]": aud}}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"lang.language_id": languageID},
				bson.M{"aud.audio_key": oldKey},
			},
		}
		opts := options.Update().SetArrayFilters(arrayFilters)

		res, err := r.topicCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			fmt.Printf("[UpsertAudio] UpdateOne failed: %v\n", err)
			return err
		}

		if res.MatchedCount > 0 {
			// Đã update thành công
			return nil
		}
	}

	// --- Step 2: Nếu oldKey rỗng hoặc không tìm thấy -> push mới ---
	filter := bson.M{"_id": objID, "language_config.language_id": languageID}
	update := bson.M{"$push": bson.M{"language_config.$.audios": aud}}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("[UpsertAudio] Push new audio failed: %v\n", err)
		return err
	}

	return nil
}

func (r *topicRepository) AddLanguageConfigToTopic(ctx context.Context, topicID string, lang model.TopicLanguageConfig) error {
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
