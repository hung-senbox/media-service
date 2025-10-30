package router

import (
	"media-service/internal/cache"
	"media-service/internal/gateway"
	"media-service/internal/media/route"
	"media-service/internal/media/v2/handler"
	"media-service/internal/media/v2/repository"
	"media-service/internal/media/v2/service"
	"media-service/internal/media/v2/usecase"
	"media-service/internal/pdf/domain"
	route2 "media-service/internal/pdf/route"
	"media-service/internal/redis"
	"media-service/pkg/config"

	cached_service "media-service/internal/cache/service"

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

	// cache setup
	systemCache := cache.NewRedisCache(cacheClientRedis)
	cachedUserGateway := cached_service.NewCachedUserGateway(userGateway, systemCache, config.AppConfig.Database.RedisCache.TTLSeconds)

	// ========================  Topic ======================== //
	// --- Repo ---
	topicRepov2 := repository.NewTopicRepository(topicCollection)

	// --- UseCase ---
	uploadTopicUseCasev2 := usecase.NewUploadTopicUseCase(topicRepov2, fileGateway, redisService)
	getTopicAppUseCasev2 := usecase.NewGetTopicAppUseCase(topicRepov2, cachedUserGateway)
	getTopicWebUseCasev2 := usecase.NewGetTopicWebUseCase(topicRepov2, cachedUserGateway, fileGateway)
	getTopicGatewayUseCasev2 := usecase.NewGetTopicGatewayUseCase(topicRepov2, cachedUserGateway, fileGateway)
	getUploadProgressUseCasev2 := usecase.NewGetUploadProgressUseCase(redisService)

	// --- Service ---
	topicServicev2 := service.NewTopicService(uploadTopicUseCasev2, getUploadProgressUseCasev2, getTopicAppUseCasev2, getTopicWebUseCasev2, getTopicGatewayUseCasev2)

	// --- Handler ---
	topicHandlerv2 := handler.NewTopicHandler(topicServicev2)
	// ========================  Topic ======================== //

	// ========================  PDF ======================== //
	pdfRepov2 := domain.NewUserResourceRepository(pdfCollection)
	pdfServicev2 := domain.NewUserResourceService(pdfRepov2, fileGateway, cachedUserGateway)
	pdfHandlerv2 := domain.NewUserResourceHandler(pdfServicev2)
	// ========================  PDF ======================== //

	// Register routes
	route.RegisterTopicRoutes(r, topicHandlerv2, cachedUserGateway)
	route2.RegisterRoutes(r, pdfHandlerv2, cachedUserGateway)
	return r
}
