package cache

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *goredis.Client
}

func NewRedisCache(client *goredis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return nil // not found
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) InvalidateUserCache(ctx context.Context, userID string) error {
	return r.Delete(ctx, r.UserCacheKey(userID))
}

func (r *RedisCache) InvalidateStudentCache(ctx context.Context, studentID string) error {
	return r.Delete(ctx, r.StudentCacheKey(studentID))
}

func (r *RedisCache) InvalidateTeacherCache(ctx context.Context, teacherID string) error {
	return r.Delete(ctx, r.TeacherCacheKey(teacherID))
}

// ==============================
// === Implement interface keys ===
// ==============================
func (r *RedisCache) UserCacheKey(userID string) string {
	return "user:" + userID
}

func (r *RedisCache) StudentCacheKey(studentID string) string {
	return "student:" + studentID
}

func (r *RedisCache) TeacherCacheKey(teacherID string) string {
	return "teacher:" + teacherID
}
