package s3

import (
	"context"
	"io"
	"time"

	"media-service/pkg/config"
	"media-service/pkg/uploader"
)

type Service interface {
	Save(ctx context.Context, data []byte, key string, mode uploader.UploadMode) (*string, error)
	SaveReader(ctx context.Context, r io.Reader, key string, contentType string, mode uploader.UploadMode) (*string, error)
	Get(ctx context.Context, key string, duration *time.Duration) (*string, error)
	Delete(ctx context.Context, key string) error
}

type service struct {
	provider uploader.UploadProvider
}

func NewFromConfig() Service {
	s3Cfg := config.AppConfig.S3.SenboxFormSubmitBucket
	provider := uploader.NewS3Provider(
		s3Cfg.AccessKey,
		s3Cfg.SecretKey,
		s3Cfg.BucketName,
		s3Cfg.Region,
		s3Cfg.Domain,
		s3Cfg.CloudfrontKeyGroupID,
		s3Cfg.CloudfrontKeyPath,
	)
	return &service{provider: provider}
}

func (s *service) Save(ctx context.Context, data []byte, key string, mode uploader.UploadMode) (*string, error) {
	return s.provider.SaveFileUploaded(ctx, data, key, mode)
}

func (s *service) SaveReader(ctx context.Context, r io.Reader, key string, contentType string, mode uploader.UploadMode) (*string, error) {
	return s.provider.SaveFileUploadedReader(ctx, r, key, contentType, mode)
}

func (s *service) Get(ctx context.Context, key string, duration *time.Duration) (*string, error) {
	return s.provider.GetFileUploaded(ctx, key, duration)
}

func (s *service) Delete(ctx context.Context, key string) error {
	return s.provider.DeleteFileUploaded(ctx, key)
}
