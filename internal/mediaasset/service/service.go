package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"media-service/internal/mediaasset/model"
	"media-service/internal/mediaasset/repository"
	"media-service/internal/s3"
	"media-service/pkg/uploader"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MediaService interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader, folder, mode string, mediaType *string) (*model.MediaAsset, *string, error)
	GetURL(ctx context.Context, id string, duration *time.Duration) (*string, error)
	GetMeta(ctx context.Context, id string) (*model.MediaAsset, error)
	GetURLByKey(ctx context.Context, key string, duration *time.Duration) (*string, error)
	Delete(ctx context.Context, id string) error
}

type mediaService struct {
	repo repository.MediaRepository
	s3   s3.Service
}

func NewMediaService(repo repository.MediaRepository) MediaService {
	return &mediaService{
		repo: repo,
		s3:   s3.NewFromConfig(),
	}
}

func (s *mediaService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, folder, mode string, mediaType *string) (*model.MediaAsset, *string, error) {
	if fileHeader == nil {
		return nil, nil, fmt.Errorf("file is required")
	}
	if folder == "" {
		folder = "uploads"
	}
	upMode, err := uploader.UploadModeFromString(mode)
	if err != nil {
		upMode = uploader.UploadPrivate
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	key := s.buildObjectKey(folder, fileHeader.Filename)
	url, err := s.s3.Save(ctx, data, key, upMode)
	if err != nil {
		return nil, nil, err
	}

	ct := http.DetectContentType(data)
	mt := detectMediaType(ct, mediaType)

	now := time.Now()
	doc := &model.MediaAsset{
		ID:          primitive.NewObjectID(),
		Type:        mt,
		Key:         key,
		FileName:    fileHeader.Filename,
		ContentType: ct,
		Size:        fileHeader.Size,
		Mode:        strings.ToLower(mode),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err = s.repo.Create(ctx, doc)
	if err != nil {
		return nil, nil, err
	}
	return doc, url, nil
}

func (s *mediaService) GetURL(ctx context.Context, id string, duration *time.Duration) (*string, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	doc, err := s.repo.GetByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("media not found")
	}
	return s.s3.Get(ctx, doc.Key, duration)
}

func (s *mediaService) GetMeta(ctx context.Context, id string) (*model.MediaAsset, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, oid)
}

func (s *mediaService) GetURLByKey(ctx context.Context, key string, duration *time.Duration) (*string, error) {
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}
	return s.s3.Get(ctx, key, duration)
}

func (s *mediaService) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	doc, err := s.repo.GetByID(ctx, oid)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("media not found")
	}
	if err := s.s3.Delete(ctx, doc.Key); err != nil {
		return err
	}
	return s.repo.DeleteByID(ctx, oid)
}

func (s *mediaService) buildObjectKey(folder, filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	if base == "" {
		base = "file"
	}
	randBytes := make([]byte, 4)
	_, _ = rand.Read(randBytes)
	randomHex := hex.EncodeToString(randBytes)
	datePath := time.Now().Format("2006/01/02")
	safeFolder := strings.Trim(folder, "/")
	safeBase := sanitizePath(base)
	safeExt := sanitizePath(ext)
	return fmt.Sprintf("%s/%s/%s-%d-%s%s", safeFolder, datePath, safeBase, time.Now().UnixNano(), randomHex, safeExt)
}

func detectMediaType(contentType string, override *string) model.MediaType {
	if override != nil && *override != "" {
		switch strings.ToLower(*override) {
		case "image":
			return model.MediaImage
		case "audio":
			return model.MediaAudio
		case "video":
			return model.MediaVideo
		case "pdf":
			return model.MediaPDF
		}
	}
	if strings.HasPrefix(contentType, "image/") {
		return model.MediaImage
	}
	if strings.HasPrefix(contentType, "audio/") {
		return model.MediaAudio
	}
	if strings.HasPrefix(contentType, "video/") {
		return model.MediaVideo
	}
	if contentType == "application/pdf" {
		return model.MediaPDF
	}
	return model.MediaOther
}

func sanitizePath(s string) string {
	s = strings.ReplaceAll(s, "..", "")
	s = strings.ReplaceAll(s, "\\", "/")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
