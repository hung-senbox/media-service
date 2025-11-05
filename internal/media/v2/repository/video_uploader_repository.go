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
	SetVideoUploaderWithoutFiles(ctx context.Context, videoUploader *model.VideoUploader) error
	SetVideoMetadata(ctx context.Context, videoUploaderID string, videoKey, videoPublicUrl string) error
	SetImagePreviewMetadata(ctx context.Context, videoUploaderID string, imageKey, imagePublicUrl string) error
	GetVideoUploaderByID(ctx context.Context, videoUploaderID string) (*model.VideoUploader, error)
	GetVideosByCreatedBy(ctx context.Context, createdBy string) ([]model.VideoUploader, error)
}

type videoUploaderRepository struct {
	videoUploaderCollection *mongo.Collection
}

func NewVideoUploaderRepository(videoUploaderCollection *mongo.Collection) VideoUploaderRepository {
	return &videoUploaderRepository{videoUploaderCollection: videoUploaderCollection}
}

func (r *videoUploaderRepository) SetVideoUploaderWithoutFiles(ctx context.Context, videoUploader *model.VideoUploader) error {
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

	// có rồi → chỉ update Title và IsVisible
	update := bson.M{
		"$set": bson.M{
			"title":      videoUploader.Title,
			"is_visible": videoUploader.IsVisible,
			"updated_at": time.Now(),
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

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"video_key":        videoKey,
			"video_public_url": videoPublicUrl,
			"updated_at":       time.Now(),
		},
	}

	_, err = r.videoUploaderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to set video metadata: %w", err)
	}
	return nil
}

func (r *videoUploaderRepository) SetImagePreviewMetadata(ctx context.Context, videoUploaderID string, imageKey, imagePublicUrl string) error {
	objID, err := primitive.ObjectIDFromHex(videoUploaderID)
	if err != nil {
		return fmt.Errorf("invalid videoUploaderID: %w", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"image_preview_key":        imageKey,
			"image_preview_public_url": imagePublicUrl,
			"updated_at":               time.Now(),
		},
	}

	_, err = r.videoUploaderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to set image preview metadata: %w", err)
	}
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
