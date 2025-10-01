package redis

import (
	"context"
	"fmt"
	"media-service/pkg/db"
)

type RedisService struct{}

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (s *RedisService) InitUploadProgress(ctx context.Context, topicID string, totalTasks int) error {
	totalKey := fmt.Sprintf("topic_upload:%s:total", topicID)
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)

	pipe := db.Client.TxPipeline()
	pipe.Set(ctx, totalKey, totalTasks, 0)
	pipe.Set(ctx, remainKey, totalTasks, 0)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisService) DecrementUploadTask(ctx context.Context, topicID string) (int64, error) {
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)
	return db.Client.Decr(ctx, remainKey).Result()
}

func (s *RedisService) GetUploadProgress(ctx context.Context, topicID string) (int64, error) {
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)
	return db.Client.Get(ctx, remainKey).Int64()
}

func (s *RedisService) GetTotalUploadTask(ctx context.Context, topicID string) (int64, error) {
	totalKey := fmt.Sprintf("topic_upload:%s:total", topicID)
	return db.Client.Get(ctx, totalKey).Int64()
}
