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
	key := fmt.Sprintf("topic_upload:%s:progress", topicID)
	return db.Client.Set(ctx, key, totalTasks, 0).Err()
}

func (s *RedisService) DecrementUploadTask(ctx context.Context, topicID string) (int64, error) {
	key := fmt.Sprintf("topic_upload:%s:progress", topicID)
	return db.Client.Decr(ctx, key).Result()
}

func (s *RedisService) GetUploadProgress(ctx context.Context, topicID string) (int64, error) {
	key := fmt.Sprintf("topic_upload:%s:progress", topicID)
	return db.Client.Get(ctx, key).Int64()
}
