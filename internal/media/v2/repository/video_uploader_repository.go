package repository

import (
	"context"
	"fmt"
	"media-service/internal/media/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoUploaderRepository interface {
	SetVideoUploader(ctx context.Context, videoUploader *model.VideoUploader) error
	SetVideoMetadata(ctx context.Context, videoUploaderID string, videoKey, videoPublicUrl string) error
	SetImagePreviewMetadata(ctx context.Context, videoUploaderID string, imageKey, imagePublicUrl string) error
	GetVideoUploaderByID(ctx context.Context, videoUploaderID string) (*model.VideoUploader, error)
	GetVideosByCreatedBy(ctx context.Context, createdBy string) ([]model.VideoUploader, error)
	GetAllVideos(ctx context.Context) ([]model.VideoUploader, error)
	GetVideosIsVisible(ctx context.Context) ([]model.VideoUploader, error)
	GetVideosByLanguageID(ctx context.Context, languageID uint) ([]model.VideoUploader, error)
	DeleteVideoUploader(ctx context.Context, videoUploaderID string) error
	DeleteVideoMetadata(ctx context.Context, videoUploaderID string, languageID uint) error
	DeleteImagePreviewMetadata(ctx context.Context, videoUploaderID string, languageID uint) error
	GetVideosByWikiCode(ctx context.Context, wikiCode string) ([]model.VideoUploader, error)
}

type videoUploaderRepository struct {
	videoUploaderCollection *mongo.Collection
}

func NewVideoUploaderRepository(videoUploaderCollection *mongo.Collection) VideoUploaderRepository {
	return &videoUploaderRepository{videoUploaderCollection: videoUploaderCollection}
}

func (r *videoUploaderRepository) SetVideoUploader(ctx context.Context, videoUploader *model.VideoUploader) error {
	filter := bson.M{"_id": videoUploader.ID}

	// kiểm tra document có tồn tại không
	count, err := r.videoUploaderCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to count video uploader: %w", err)
	}

	if count == 0 {
		// chưa có → tạo mới
		videoUploader.CreatedAt = time.Now()
		videoUploader.UpdatedAt = time.Now()
		if videoUploader.ID.IsZero() {
			videoUploader.ID = primitive.NewObjectID()
		}
		_, err := r.videoUploaderCollection.InsertOne(ctx, videoUploader)
		if err != nil {
			return fmt.Errorf("failed to insert video uploader: %w", err)
		}
		return nil
	}

	// có rồi → update basic fields + language_config
	update := bson.M{
		"$set": bson.M{
			"is_visible":      videoUploader.IsVisible,
			"title":           videoUploader.Title,
			"wiki_code":       videoUploader.WikiCode,
			"language_config": videoUploader.LanguageConfig,
			"updated_at":      time.Now(),
		},
	}

	_, err = r.videoUploaderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update video uploader: %w", err)
	}

	return nil
}

func (r *videoUploaderRepository) SetVideoMetadata(ctx context.Context, videoUploaderID string, videoKey, videoPublicUrl string) error {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return fmt.Errorf("invalid videoUploaderID: %w", err)
	}

	// Giữ lại function để tương thích, nhưng hiện tại metadata được cập nhật
	// thông qua SetVideoUploaderWithoutFiles, nên không làm gì thêm ở đây.
	_ = objID
	_ = videoKey
	_ = videoPublicUrl
	return nil
}

func (r *videoUploaderRepository) SetImagePreviewMetadata(ctx context.Context, videoUploaderID string, imageKey, imagePublicUrl string) error {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return fmt.Errorf("invalid videoUploaderID: %w", err)
	}

	// Giữ lại function để tương thích, nhưng hiện tại metadata được cập nhật
	// thông qua SetVideoUploaderWithoutFiles, nên không làm gì thêm ở đây.
	_ = objID
	_ = imageKey
	_ = imagePublicUrl
	return nil
}

func (r *videoUploaderRepository) GetVideoUploaderByID(ctx context.Context, videoUploaderID string) (*model.VideoUploader, error) {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return nil, fmt.Errorf("invalid videoUploaderID: %w", err)
	}

	var videoUploader model.VideoUploader
	err = r.videoUploaderCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&videoUploader)
	if err != nil {
		return nil, fmt.Errorf("failed to get video uploader: %w", err)
	}

	return &videoUploader, nil
}

func (r *videoUploaderRepository) GetVideosByCreatedBy(ctx context.Context, createdBy string) ([]model.VideoUploader, error) {
	var videoUploaders []model.VideoUploader
	cursor, err := r.videoUploaderCollection.Find(ctx, bson.M{"created_by": createdBy})
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by created by: %w", err)
	}
	err = cursor.All(ctx, &videoUploaders)
	if err != nil {
		return nil, fmt.Errorf("failed to decode videos by created by: %w", err)
	}
	return videoUploaders, nil
}

func (r *videoUploaderRepository) GetAllVideos(ctx context.Context) ([]model.VideoUploader, error) {
	var videoUploaders []model.VideoUploader
	cursor, err := r.videoUploaderCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all videos: %w", err)
	}
	err = cursor.All(ctx, &videoUploaders)
	if err != nil {
		return nil, fmt.Errorf("failed to decode all videos: %w", err)
	}
	return videoUploaders, nil
}

func (r *videoUploaderRepository) GetVideosIsVisible(ctx context.Context) ([]model.VideoUploader, error) {
	var videoUploaders []model.VideoUploader
	cursor, err := r.videoUploaderCollection.Find(ctx, bson.M{"is_visible": true})
	if err != nil {
		return nil, fmt.Errorf("failed to get videos is visible: %w", err)
	}
	err = cursor.All(ctx, &videoUploaders)
	if err != nil {
		return nil, fmt.Errorf("failed to decode videos is visible: %w", err)
	}
	return videoUploaders, nil
}

func (r *videoUploaderRepository) GetVideosByLanguageID(ctx context.Context, languageID uint) ([]model.VideoUploader, error) {
	var videoUploaders []model.VideoUploader
	// lọc theo language_config.language_id
	cursor, err := r.videoUploaderCollection.Find(ctx, bson.M{"language_config.language_id": languageID})
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by language id: %w", err)
	}
	err = cursor.All(ctx, &videoUploaders)
	if err != nil {
		return nil, fmt.Errorf("failed to decode videos by language id: %w", err)
	}
	return videoUploaders, nil
}

func (r *videoUploaderRepository) DeleteVideoUploader(ctx context.Context, videoUploaderID string) error {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return fmt.Errorf("invalid videoUploaderID: %w", err)
	}
	_, err = r.videoUploaderCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete video uploader: %w", err)
	}
	return nil
}

func (r *videoUploaderRepository) DeleteVideoMetadata(ctx context.Context, videoUploaderID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return fmt.Errorf("invalid videoUploaderID: %w", err)
	}
	filter := bson.M{"_id": objID, "language_config.language_id": languageID}
	update := bson.M{"$set": bson.M{"language_config.$.video_key": "", "language_config.$.video_public_url": "", "updated_at": time.Now()}}
	_, err = r.videoUploaderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete video metadata: %w", err)
	}
	return nil

}

func (r *videoUploaderRepository) DeleteImagePreviewMetadata(ctx context.Context, videoUploaderID string, languageID uint) error {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return fmt.Errorf("invalid videoUploaderID: %w", err)
	}
	filter := bson.M{"_id": objID, "language_config.language_id": languageID}
	update := bson.M{"$set": bson.M{"language_config.$.image_preview_key": "", "language_config.$.image_preview_public_url": "", "updated_at": time.Now()}}
	_, err = r.videoUploaderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete image preview metadata: %w", err)
	}
	return nil
}

func (r *videoUploaderRepository) GetVideosByWikiCode(ctx context.Context, wikiCode string) ([]model.VideoUploader, error) {
	var videoUploaders []model.VideoUploader
	cursor, err := r.videoUploaderCollection.Find(ctx, bson.M{"wiki_code": wikiCode})
	if err != nil {
		return nil, fmt.Errorf("failed to get videos by wiki code: %w", err)
	}
	err = cursor.All(ctx, &videoUploaders)
	if err != nil {
		return nil, fmt.Errorf("failed to decode videos by wiki code: %w", err)
	}
	return videoUploaders, nil
}
