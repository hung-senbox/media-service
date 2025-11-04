package repository

import (
	"context"
	"errors"
	"fmt"
	"media-service/internal/media/model"
	"strings"
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
	InitImages(ctx context.Context, topicID string, languageID uint) error
	DeleteAudioKey(ctx context.Context, topicID string, languageID uint) error
	DeleteVideoKey(ctx context.Context, topicID string, languageID uint) error
	DeleteImageKey(ctx context.Context, topicID string, languageID uint, imageType string) error
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

	// Chuẩn hóa image_type
	img.ImageType = strings.TrimSpace(strings.TrimSuffix(img.ImageType, ","))
	img.ImageType = strings.ToLower(img.ImageType)

	langVariants := []interface{}{languageID, int32(languageID), int64(languageID)}

	// 1️⃣ Đảm bảo tồn tại language_config
	langFilter := bson.M{
		"_id":                         objID,
		"language_config.language_id": bson.M{"$in": langVariants},
	}
	count, _ := r.topicCollection.CountDocuments(ctx, langFilter)
	if count == 0 {
		newLang := bson.M{
			"language_id": languageID,
			"images":      []interface{}{},
		}
		_, err := r.topicCollection.UpdateOne(ctx,
			bson.M{"_id": objID},
			bson.M{"$push": bson.M{"language_config": newLang}},
		)
		if err != nil {
			return fmt.Errorf("[SetImageForTopic] create new language_config failed: %w", err)
		}
	}

	// 2️⃣ Cập nhật ảnh nếu image_type đã có
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"language_config.$[lang].images.$[img].image_key":    img.ImageKey,
			"language_config.$[lang].images.$[img].link_url":     img.LinkUrl,
			"language_config.$[lang].images.$[img].uploaded_url": img.UploadedUrl,
			"language_config.$[lang].images.$[img].image_type":   img.ImageType,
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"lang.language_id": bson.M{"$in": langVariants}},
			bson.M{"img.image_type": img.ImageType},
		},
	})

	res, err := r.topicCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("[SetImageForTopic] update failed: %w", err)
	}

	// 3️⃣ Nếu chưa có image_type → push mới
	if res.MatchedCount == 0 {
		pushUpdate := bson.M{
			"$push": bson.M{
				"language_config.$[lang].images": bson.M{
					"image_type":   img.ImageType,
					"image_key":    img.ImageKey,
					"link_url":     img.LinkUrl,
					"uploaded_url": img.UploadedUrl,
				},
			},
		}
		pushOpts := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"lang.language_id": bson.M{"$in": langVariants}},
			},
		})
		_, err = r.topicCollection.UpdateOne(ctx, filter, pushUpdate, pushOpts)
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

	// 1) Update or insert language-level metadata WITHOUT replacing entire language_config
	if len(topic.LanguageConfig) > 0 {
		for _, lc := range topic.LanguageConfig {
			langFilter := bson.M{
				"_id":                         topic.ID,
				"language_config.language_id": lc.LanguageID,
			}

			langUpdate := bson.M{
				"$set": bson.M{
					"language_config.$.file_name":   lc.FileName,
					"language_config.$.title":       lc.Title,
					"language_config.$.note":        lc.Note,
					"language_config.$.description": lc.Description,
					"updated_at":                    time.Now(),
				},
			}

			res, err := r.topicCollection.UpdateOne(ctx, langFilter, langUpdate)
			if err != nil {
				return nil, fmt.Errorf("update language config failed: %w", err)
			}

			if res.MatchedCount == 0 {
				// Insert new language config entry preserving media subdocuments by initializing empty values
				pushUpdate := bson.M{
					"$push": bson.M{
						"language_config": bson.M{
							"language_id": lc.LanguageID,
							"file_name":   lc.FileName,
							"title":       lc.Title,
							"note":        lc.Note,
							"description": lc.Description,
							"images":      []model.TopicImageConfig{},
							"audio":       model.TopicAudioConfig{},
							"video":       model.TopicVideoConfig{},
						},
					},
				}
				if _, err := r.topicCollection.UpdateOne(ctx, bson.M{"_id": topic.ID}, pushUpdate); err != nil {
					return nil, fmt.Errorf("insert language config failed: %w", err)
				}
			}
		}
	}

	// 2) Update top-level fields and return the updated document
	update := bson.M{"$set": bson.M{
		"is_published": topic.IsPublished,
		"updated_at":   time.Now(),
	}}

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

func (r *topicRepository) InitImages(ctx context.Context, topicID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[InitImages] invalid topicID=%s: %w", topicID, err)
	}

	// Danh sách 6 loại hình mặc định
	defaultImageTypes := []string{
		"full_background",
		"clear_background",
		"clip_part",
		"drawing",
		"icon",
		"bm",
	}

	// Tạo danh sách hình mặc định
	var images []model.TopicImageConfig
	for _, t := range defaultImageTypes {
		images = append(images, model.TopicImageConfig{
			ImageType: t,
			ImageKey:  "",
			LinkUrl:   "",
		})
	}

	// Cập nhật images vào language_config có language_id tương ứng
	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$.images": images,
			"updated_at":               time.Now(),
		},
	}

	res, err := r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[InitImages] update failed: %w", err)
	}

	// Nếu chưa tồn tại language_id thì thêm mới
	if res.MatchedCount == 0 {
		newLang := model.TopicLanguageConfig{
			LanguageID: languageID,
			Images:     images,
		}
		_, err = r.topicCollection.UpdateOne(ctx,
			bson.M{"_id": objID},
			bson.M{"$push": bson.M{"language_config": newLang}},
		)
		if err != nil {
			return fmt.Errorf("[InitImages] push new language_config failed: %w", err)
		}
	}

	return nil
}

func (r *topicRepository) DeleteAudioKey(ctx context.Context, topicID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[DeleteAudioKey] invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$.audio": bson.M{
				"audio_key": "",
			},
		},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[DeleteAudioKey] update failed: %w", err)
	}

	return nil
}

func (r *topicRepository) DeleteVideoKey(ctx context.Context, topicID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[DeleteVideoKey] invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$.video": bson.M{
				"video_key": "",
			},
		},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[DeleteVideoKey] update failed: %w", err)
	}

	return nil
}

func (r *topicRepository) DeleteImageKey(ctx context.Context, topicID string, languageID uint, imageType string) error {
	objID, err := primitive.ObjectIDFromHex(topicID)
	if err != nil {
		return fmt.Errorf("[DeleteImageKey] invalid topicID=%s: %w", topicID, err)
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}

	update := bson.M{
		"$set": bson.M{
			"language_config.$.images": bson.M{
				"image_key": "",
			},
		},
	}

	_, err = r.topicCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[DeleteImageKey] update failed: %w", err)
	}

	return nil
}
