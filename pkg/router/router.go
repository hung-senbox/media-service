package router

import (
	"media-service/internal/gateway"
	"media-service/internal/media/handler"
	"media-service/internal/media/repository"
	"media-service/internal/media/route"
	"media-service/internal/media/service"
	"media-service/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(consulClient *api.Client, topicCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	fileGateway := gateway.NewFileGateway("go-main-service", consulClient)
	redisService := redis.NewRedisService()

	// topic
	topicRepo := repository.NewTopicRepository(topicCollection)
	topicService := service.NewTopicService(topicRepo, fileGateway, redisService, userGateway)
	topicHandler := handler.NewTopicHandler(topicService)

	// Register routes
	route.RegisterTopicRoutes(r, topicHandler)
	return r
}
