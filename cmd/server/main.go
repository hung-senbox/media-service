package main

import (
	"fmt"
	"log"
	"os"
	"time"

	// "os"

	"media-service/pkg/config"
	"media-service/pkg/consul"
	"media-service/pkg/db"
	"media-service/pkg/router"

	"media-service/pkg/zap"

	"github.com/gofiber/fiber/v2"
	consulapi "github.com/hashicorp/consul/api"
	redis "github.com/hung-senbox/senbox-cache-service/pkg/redis"
)

func main() {
	filePath := os.Args[1]
	if filePath == "" {
		filePath = "configs/config.yaml"
	}

	config.LoadConfig(filePath)

	cfg := config.AppConfig

	//logger
	logger, err := zap.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	//consul
	consulConn := consul.NewConsulConn(logger, cfg)
	consulClient := consulConn.Connect()
	defer consulConn.Deregister()

	if err := waitPassing(consulClient, "go-main-service", 60*time.Second); err != nil {
		logger.Fatalf("Dependency not ready: %v", err)
	}

	//db
	db.ConnectMongoDB()

	//redis
	db.ConnectRedis()
	defer db.Client.Close()

	// redis cache
	cacheClientRedis, err := redis.InitRedisCache(config.AppConfig.Database.RedisCache.Host, config.AppConfig.Database.RedisCache.Port, config.AppConfig.Database.RedisCache.Password, config.AppConfig.Database.RedisCache.DB)
	if err != nil {
		logger.Fatalf("Failed to initialize Redis cache: %v", err)
		return
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
	})

	app = router.SetupRouter(
		app,
		consulClient,
		cacheClientRedis,
		db.TopicCollection,
		db.PDFCollection,
		db.TopicResourceCollection,
		db.VideoUploaderCollection,
		db.MediaAssetCollection,
		db.VocabularyCollection,
	)
	port := cfg.Server.Port
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

func waitPassing(cli *consulapi.Client, name string, timeout time.Duration) error {
	dl := time.Now().Add(timeout)
	for time.Now().Before(dl) {
		entries, _, err := cli.Health().Service(name, "", true, nil)
		if err == nil && len(entries) > 0 {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("%s not ready in consul", name)
}
