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
	switch v := value.(type) {
	case string:
		return r.client.Set(ctx, key, v, time.Duration(ttlSeconds)*time.Second).Err()
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
	}
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
