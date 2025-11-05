package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"media-service/pkg/db"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisService struct{}

func NewRedisService() *RedisService {
	return &RedisService{}
}

// helper để build key
func buildKey(organizationID, topicID, field string) string {
	return fmt.Sprintf("topic_upload:%s:%s:%s", organizationID, topicID, field)
}

func (s *RedisService) SetUploaderStatus(ctx context.Context, key string, values map[string]interface{}) error {
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal redis value: %w", err)
	}

	return db.Client.Set(ctx, key, data, 0).Err()
}

func (s *RedisService) GetUploaderStatus(ctx context.Context, key string) (map[string]interface{}, error) {
	val, err := db.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get redis value: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal redis value: %w", err)
	}

	return result, nil
}

func (s *RedisService) DeleteUploaderStatusKey(ctx context.Context, key string) error {
	return db.Client.Del(ctx, key).Err()
}

// Khởi tạo progress upload
func (s *RedisService) InitUploadProgress(ctx context.Context, organizationID, topicID string, totalTasks int) error {
	totalKey := buildKey(organizationID, topicID, "total")
	remainKey := buildKey(organizationID, topicID, "remaining")

	pipe := db.Client.TxPipeline()
	pipe.Set(ctx, totalKey, totalTasks, 0)
	pipe.Set(ctx, remainKey, totalTasks, 0)
	_, err := pipe.Exec(ctx)
	return err
}

// Giảm số task còn lại
func (s *RedisService) DecrementUploadTask(ctx context.Context, organizationID, topicID string) (int64, error) {
	remainKey := buildKey(organizationID, topicID, "remaining")
	return db.Client.Decr(ctx, remainKey).Result()
}

// SetUploadProgress cập nhật progress (%) cho topic
func (s *RedisService) SetUploadProgress(ctx context.Context, organizationID, topicID string, progress int) error {
	key := fmt.Sprintf("topic_upload:%s:%s:progress", organizationID, topicID)
	return db.Client.Set(ctx, key, progress, 0).Err()
}

// Lấy số task còn lại
func (s *RedisService) GetUploadProgress(ctx context.Context, organizationID, topicID string) (int64, error) {
	remainKey := buildKey(organizationID, topicID, "remaining")
	val, err := db.Client.Get(ctx, remainKey).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// Lấy tổng task
func (s *RedisService) GetTotalUploadTask(ctx context.Context, organizationID, topicID string) (int64, error) {
	totalKey := buildKey(organizationID, topicID, "total")
	val, err := db.Client.Get(ctx, totalKey).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(val, 10, 64)
}

// Lưu lỗi upload
func (s *RedisService) SetUploadError(ctx context.Context, organizationID, topicID, key, errMsg string) error {
	redisKey := buildKey(organizationID, topicID, "errors")
	return db.Client.HSet(ctx, redisKey, key, errMsg).Err()
}

// Lấy tất cả lỗi
func (s *RedisService) GetUploadErrors(ctx context.Context, organizationID, topicID string) (map[string]string, error) {
	redisKey := buildKey(organizationID, topicID, "errors")
	return db.Client.HGetAll(ctx, redisKey).Result()
}

// Xoá progress + errors
func (s *RedisService) DeleteUploadProgress(ctx context.Context, organizationID, topicID string) error {
	totalKey := buildKey(organizationID, topicID, "total")
	remainKey := buildKey(organizationID, topicID, "remaining")
	errorKey := buildKey(organizationID, topicID, "errors")
	progressKey := buildKey(organizationID, topicID, "progress")
	pipe := db.Client.TxPipeline()
	pipe.Del(ctx, totalKey)
	pipe.Del(ctx, remainKey)
	pipe.Del(ctx, errorKey)
	pipe.Del(ctx, progressKey)
	_, err := pipe.Exec(ctx)
	return err
}

// HasAnyUploadInProgress check xem trong org còn topic nào chưa upload xong không
func (s *RedisService) HasAnyUploadInProgress(ctx context.Context, organizationID string) (bool, error) {
	// pattern cho tất cả remaining key của org này
	pattern := fmt.Sprintf("topic_upload:%s:*:remaining", organizationID)

	iter := db.Client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		remainVal, err := db.Client.Get(ctx, iter.Val()).Int64()
		if err != nil {
			// Nếu key không tồn tại hoặc lỗi parse → bỏ qua
			continue
		}
		if remainVal > 0 {
			return true, nil
		}
	}
	if err := iter.Err(); err != nil {
		return false, err
	}
	return false, nil
}
