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

type VocabularyRepository interface {
	CreateVocabulary(ctx context.Context, vocabulary *model.Vocabulary) (*model.Vocabulary, error)
	SetLanguageConfig(ctx context.Context, vocabularyID string, lang model.VocabularyLanguageConfig) error
	InitImages(ctx context.Context, vocabularyID string, languageID uint) error
	GetByID(ctx context.Context, vocabularyID string) (*model.Vocabulary, error)
	UpdateVocabulary(ctx context.Context, vocabulary *model.Vocabulary) (*model.Vocabulary, error)
	DeleteAudioKey(ctx context.Context, vocabularyID string, languageID uint) error
	SetAudio(ctx context.Context, vocabularyID string, languageID uint, aud model.VocabularyAudioConfig) error
	DeleteVideoKey(ctx context.Context, vocabularyID string, languageID uint) error
	SetVideo(ctx context.Context, vocabularyID string, languageID uint, vid model.VocabularyVideoConfig) error
	DeleteImageKey(ctx context.Context, vocabularyID string, languageID uint, imageType string) error
	SetImage(ctx context.Context, vocabularyID string, languageID uint, img model.VocabularyImageConfig) error
	GetAllVocabulariesByTopicID(ctx context.Context, topicID string) ([]model.Vocabulary, error)
}

type vocabularyRepository struct {
	vocabularyCollection *mongo.Collection
}

func NewVocabularyRepository(vocabularyCollection *mongo.Collection) VocabularyRepository {
	return &vocabularyRepository{vocabularyCollection: vocabularyCollection}
}

func (r *vocabularyRepository) CreateVocabulary(ctx context.Context, vocabulary *model.Vocabulary) (*model.Vocabulary, error) {
	result, err := r.vocabularyCollection.InsertOne(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("insert vocabulary failed: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		vocabulary.ID = oid
	}

	return vocabulary, nil
}

func (r *vocabularyRepository) SetLanguageConfig(ctx context.Context, vocabularyID string, lang model.VocabularyLanguageConfig) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("invalid vocabularyID: %w", err)
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

	res, err := r.vocabularyCollection.UpdateOne(ctx, filterUpdate, update)
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

	_, err = r.vocabularyCollection.UpdateOne(ctx, filterPush, updatePush)
	if err != nil {
		return fmt.Errorf("insert language config failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) InitImages(ctx context.Context, vocabularyID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[InitImages] invalid vocabularyID=%s: %w", vocabularyID, err)
	}

	// Danh sách 9 loại hình mặc định
	defaultImageTypes := []string{
		"full_background",
		"clear_background",
		"clip_part",
		"drawing",
		"icon",
		"bm",
		"sign_lang",
		"gif",
		"order",
	}

	// Lấy language_config hiện tại cho languageID
	var vocabulary model.Vocabulary
	if err := r.vocabularyCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&vocabulary); err != nil {
		return fmt.Errorf("[InitImages] get vocabulary failed: %w", err)
	}

	var (
		existsLang    = false
		currentImages []model.VocabularyImageConfig
	)
	for _, lc := range vocabulary.LanguageConfig {
		if lc.LanguageID == languageID {
			existsLang = true
			currentImages = lc.Images
			break
		}
	}

	// Nếu chưa có language_config -> push mới với full defaults
	if !existsLang {
		var images []model.VocabularyImageConfig
		for _, t := range defaultImageTypes {
			images = append(images, model.VocabularyImageConfig{
				ImageType: t,
				ImageKey:  "",
				LinkUrl:   "",
			})
		}
		newLang := model.VocabularyLanguageConfig{
			LanguageID: languageID,
			Images:     images,
		}
		_, err = r.vocabularyCollection.UpdateOne(ctx,
			bson.M{"_id": objID},
			bson.M{"$push": bson.M{"language_config": newLang}},
		)
		if err != nil {
			return fmt.Errorf("[InitImages] push new language_config failed: %w", err)
		}
		return nil
	}

	// Nếu đã có images:
	// - length == 0 => set default list
	// - length < 8 => đảm bảo thêm 3 type mới: sign_lang, gif, order
	// - else => không thay đổi
	var updatedImages []model.VocabularyImageConfig
	if len(currentImages) == 0 {
		for _, t := range defaultImageTypes {
			updatedImages = append(updatedImages, model.VocabularyImageConfig{
				ImageType: t,
				ImageKey:  "",
				LinkUrl:   "",
			})
		}
	} else {
		// copy hiện có
		updatedImages = append(updatedImages, currentImages...)
		if len(currentImages) < 8 {
			// map để check tồn tại
			exists := map[string]bool{}
			for _, img := range currentImages {
				exists[strings.ToLower(img.ImageType)] = true
			}
			for _, t := range []string{"sign_lang", "gif", "order"} {
				if !exists[t] {
					updatedImages = append(updatedImages, model.VocabularyImageConfig{
						ImageType: t,
						ImageKey:  "",
						LinkUrl:   "",
					})
				}
			}
		} else {
			// Không cần cập nhật nếu đã đủ (>=8)
			return nil
		}
	}

	filter := bson.M{
		"_id":                         objID,
		"language_config.language_id": languageID,
	}
	update := bson.M{
		"$set": bson.M{
			"language_config.$.images": updatedImages,
			"updated_at":               time.Now(),
		},
	}

	if _, err := r.vocabularyCollection.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("[InitImages] update images failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) GetByID(ctx context.Context, vocabularyID string) (*model.Vocabulary, error) {
	var vocabulary model.Vocabulary
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return nil, fmt.Errorf("[GetByID] invalid vocabularyID=%s: %w", vocabularyID, err)
	}
	if err := r.vocabularyCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&vocabulary); err != nil {
		return nil, fmt.Errorf("[GetByID] get vocabulary failed: %w", err)
	}
	return &vocabulary, nil
}

func (r *vocabularyRepository) UpdateVocabulary(ctx context.Context, vocabulary *model.Vocabulary) (*model.Vocabulary, error) {
	filter := bson.M{"_id": vocabulary.ID}

	// 1) Update or insert language-level metadata WITHOUT replacing entire language_config
	if len(vocabulary.LanguageConfig) > 0 {
		for _, lc := range vocabulary.LanguageConfig {
			langFilter := bson.M{
				"_id":                         vocabulary.ID,
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

			res, err := r.vocabularyCollection.UpdateOne(ctx, langFilter, langUpdate)
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
							"images":      []model.VocabularyImageConfig{},
							"audio":       model.TopicAudioConfig{},
							"video":       model.TopicVideoConfig{},
						},
					},
				}
				if _, err := r.vocabularyCollection.UpdateOne(ctx, bson.M{"_id": vocabulary.ID}, pushUpdate); err != nil {
					return nil, fmt.Errorf("insert language config failed: %w", err)
				}
			}
		}
	}

	// 2) Update top-level fields and return the updated document
	update := bson.M{"$set": bson.M{
		"is_published": vocabulary.IsPublished,
		"updated_at":   time.Now(),
	}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedVocabulary model.Vocabulary
	err := r.vocabularyCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedVocabulary)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("vocabulary not found")
		}
		return nil, fmt.Errorf("update vocabulary failed: %w", err)
	}

	return &updatedVocabulary, nil
}

func (r *vocabularyRepository) DeleteAudioKey(ctx context.Context, vocabularyID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[DeleteAudioKey] invalid vocabularyID=%s: %w", vocabularyID, err)
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

	_, err = r.vocabularyCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[DeleteAudioKey] update failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) SetAudio(ctx context.Context, vocabularyID string, languageID uint, aud model.VocabularyAudioConfig) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[SetAudio4Vocabulary] invalid vocabularyID=%s: %w", vocabularyID, err)
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

	_, err = r.vocabularyCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[SetAudioForVocabulary] update failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) SetVideo(ctx context.Context, vocabularyID string, languageID uint, vid model.VocabularyVideoConfig) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[SetVideo4Vocabulary] invalid vocabularyID=%s: %w", vocabularyID, err)
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

	_, err = r.vocabularyCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[SetVideoForVocabulary] update failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) DeleteVideoKey(ctx context.Context, vocabularyID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[DeleteVideoKey] invalid vocabularyID=%s: %w", vocabularyID, err)
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

	_, err = r.vocabularyCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("[DeleteVideoKey] update failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) DeleteImageKey(ctx context.Context, vocabularyID string, languageID uint, imageType string) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[DeleteImageKey] invalid vocabularyID=%s: %w", vocabularyID, err)
	}

	// Chuẩn hóa imageType và hỗ trợ nhiều kiểu số cho language_id
	imgType := strings.TrimSpace(strings.TrimSuffix(strings.ToLower(imageType), ","))
	langVariants := []interface{}{languageID, int32(languageID), int64(languageID)}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"language_config.$[lang].images.$[img].image_key": "",
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"lang.language_id": bson.M{"$in": langVariants}},
			bson.M{"img.image_type": imgType},
		},
	})

	_, err = r.vocabularyCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("[DeleteImageKey] update failed: %w", err)
	}

	return nil
}

func (r *vocabularyRepository) SetImage(ctx context.Context, vocabularyID string, languageID uint, img model.VocabularyImageConfig) error {
	objID, err := primitive.ObjectIDFromHex(vocabularyID)
	if err != nil {
		return fmt.Errorf("[SetImageForVocabulary] invalid vocabularyID=%s: %w", vocabularyID, err)
	}

	// Chuẩn hóa image_type
	img.ImageType = strings.TrimSpace(strings.TrimSuffix(img.ImageType, ","))
	img.ImageType = strings.ToLower(img.ImageType)

	langVariants := []interface{}{languageID, int32(languageID), int64(languageID)}

	// Đảm bảo tồn tại language_config
	langFilter := bson.M{
		"_id":                         objID,
		"language_config.language_id": bson.M{"$in": langVariants},
	}
	count, _ := r.vocabularyCollection.CountDocuments(ctx, langFilter)
	if count == 0 {
		newLang := bson.M{
			"language_id": languageID,
			"images":      []interface{}{},
		}
		_, err := r.vocabularyCollection.UpdateOne(ctx,
			bson.M{"_id": objID},
			bson.M{"$push": bson.M{"language_config": newLang}},
		)
		if err != nil {
			return fmt.Errorf("[SetImageForVocabulary] create new language_config failed: %w", err)
		}
	}

	// Cập nhật ảnh nếu image_type đã có
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

	res, err := r.vocabularyCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("[SetImageForVocabulary] update failed: %w", err)
	}

	//  Nếu chưa có image_type → push mới
	if res.MatchedCount == 0 {
		pushUpdate := bson.M{
			"$push": bson.M{
				"language_config.$[lang].images": bson.M{
					"image_type": img.ImageType,
					"image_key":  img.ImageKey,
					"link_url":   img.LinkUrl,
				},
			},
		}
		pushOpts := options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"lang.language_id": bson.M{"$in": langVariants}},
			},
		})
		_, err = r.vocabularyCollection.UpdateOne(ctx, filter, pushUpdate, pushOpts)
		if err != nil {
			return fmt.Errorf("[SetImageForTopic] push failed: %w", err)
		}
	}

	return nil
}

func (r *vocabularyRepository) GetAllVocabulariesByTopicID(ctx context.Context, topicID string) ([]model.Vocabulary, error) {
	var vocabularies []model.Vocabulary
	filter := bson.M{"topic_id": topicID}
	cursor, err := r.vocabularyCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get all vocabularies by topic id failed: %w", err)
	}
	if err := cursor.All(ctx, &vocabularies); err != nil {
		return nil, fmt.Errorf("decode vocabularies failed: %w", err)
	}
	return vocabularies, nil
}
