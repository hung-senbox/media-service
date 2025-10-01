package redis

import (
	"context"
	"fmt"
	"media-service/pkg/db"
	"strconv"
)

type RedisService struct{}

func NewRedisService() *RedisService {
	return &RedisService{}
}

// Khởi tạo progress upload
func (s *RedisService) InitUploadProgress(ctx context.Context, topicID string, totalTasks int) error {
	totalKey := fmt.Sprintf("topic_upload:%s:total", topicID)
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)

	pipe := db.Client.TxPipeline()
	pipe.Set(ctx, totalKey, totalTasks, 0)
	pipe.Set(ctx, remainKey, totalTasks, 0)
	_, err := pipe.Exec(ctx)
	return err
}

// Giảm số task còn lại
func (s *RedisService) DecrementUploadTask(ctx context.Context, topicID string) (int64, error) {
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)
	return db.Client.Decr(ctx, remainKey).Result()
}

// Lấy số task còn lại
func (s *RedisService) GetUploadProgress(ctx context.Context, topicID string) (int64, error) {
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)
	val, err := db.Client.Get(ctx, remainKey).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// Lấy tổng task
func (s *RedisService) GetTotalUploadTask(ctx context.Context, topicID string) (int64, error) {
	totalKey := fmt.Sprintf("topic_upload:%s:total", topicID)
	val, err := db.Client.Get(ctx, totalKey).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// Lưu lỗi upload
func (s *RedisService) SetUploadError(ctx context.Context, topicID, key, errMsg string) error {
	redisKey := fmt.Sprintf("topic_upload:%s:errors", topicID)
	return db.Client.HSet(ctx, redisKey, key, errMsg).Err()
}

// Lấy tất cả lỗi
func (s *RedisService) GetUploadErrors(ctx context.Context, topicID string) (map[string]string, error) {
	redisKey := fmt.Sprintf("topic_upload:%s:errors", topicID)
	return db.Client.HGetAll(ctx, redisKey).Result()
}

// Xoá progress + errors
func (s *RedisService) DeleteUploadProgress(ctx context.Context, topicID string) error {
	totalKey := fmt.Sprintf("topic_upload:%s:total", topicID)
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)
	errorKey := fmt.Sprintf("topic_upload:%s:errors", topicID)

	pipe := db.Client.TxPipeline()
	pipe.Del(ctx, totalKey)
	pipe.Del(ctx, remainKey)
	pipe.Del(ctx, errorKey)
	_, err := pipe.Exec(ctx)
	return err
}

// Cập nhật remaining (dùng khi cần)
func (s *RedisService) SetUploadProgress(ctx context.Context, topicID string, remaining int64) error {
	remainKey := fmt.Sprintf("topic_upload:%s:remaining", topicID)
	return db.Client.Set(ctx, remainKey, remaining, 0).Err()
}
