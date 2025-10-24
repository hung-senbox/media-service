package router

import (
	"media-service/internal/gateway"
	"media-service/internal/media/route"
	"media-service/internal/media/v2/handler"
	"media-service/internal/media/v2/repository"
	"media-service/internal/media/v2/service"
	"media-service/internal/pdf/domain"
	route2 "media-service/internal/pdf/route"
	"media-service/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(consulClient *api.Client, cacheClientRedis *goredis.Client, topicCollection, pdfCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	fileGateway := gateway.NewFileGateway("go-main-service", consulClient)
	redisService := redis.NewRedisService()

	// topic
	topicRepov2 := repository.NewTopicRepository(topicCollection)
	topicServicev2 := service.NewTopicService(topicRepov2, fileGateway, redisService, userGateway)
	topicHandlerv2 := handler.NewTopicHandler(topicServicev2)

	//pdf
	pdfRepov2 := domain.NewPDFRepository(pdfCollection)
	pdfServicev2 := domain.NewPDFService(pdfRepov2, fileGateway)
	pdfHandlerv2 := domain.NewPDFHandler(pdfServicev2)
	// Register routes
	route.RegisterTopicRoutes(r, topicHandlerv2)
	route2.RegisterRoutes(r, pdfHandlerv2)
	return r
}
