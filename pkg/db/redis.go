package db

import (
	"context"
	"fmt"
	"media-service/pkg/config"

	goredis "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *goredis.Client

func ConnectRedis() {
	cfg := config.AppConfig.Redis
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	Client = goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// test kết nối
	if err := Client.Ping(Ctx).Err(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
}
