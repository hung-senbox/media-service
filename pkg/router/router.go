package router

import (
	"media-service/internal/gateway"
	"media-service/internal/media/route"
	"media-service/internal/media/v2/handler"
	"media-service/internal/media/v2/repository"
	"media-service/internal/media/v2/service"
	"media-service/internal/media/v2/usecase"
	"media-service/internal/pdf/domain"
	route2 "media-service/internal/pdf/route"
	"media-service/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/cached"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(consulClient *api.Client, cacheClientRedis *cache.RedisCache, topicCollection, pdfCollection, topicResourceCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	cachedMainGateway := cached.NewCachedMainGateway(cacheClientRedis)
	userGateway := gateway.NewUserGateway("go-main-service", consulClient, cachedMainGateway)
	fileGateway := gateway.NewFileGateway("go-main-service", consulClient)
	redisService := redis.NewRedisService()

	// ========================  Topic ======================== //
	// --- Repo ---
	topicRepov2 := repository.NewTopicRepository(topicCollection)

	// --- UseCase ---
	uploadTopicUseCasev2 := usecase.NewUploadTopicUseCase(topicRepov2, fileGateway, redisService)
	getTopicAppUseCasev2 := usecase.NewGetTopicAppUseCase(topicRepov2, userGateway)
	getTopicWebUseCasev2 := usecase.NewGetTopicWebUseCase(topicRepov2, userGateway, fileGateway)
	getTopicGatewayUseCasev2 := usecase.NewGetTopicGatewayUseCase(topicRepov2, userGateway, fileGateway)
	getUploadProgressUseCasev2 := usecase.NewGetUploadProgressUseCase(topicRepov2, redisService)

	// --- Service ---
	topicServicev2 := service.NewTopicService(uploadTopicUseCasev2, getUploadProgressUseCasev2, getTopicAppUseCasev2, getTopicWebUseCasev2, getTopicGatewayUseCasev2)

	// --- Handler ---
	topicHandlerv2 := handler.NewTopicHandler(topicServicev2)
	// ========================  Topic ======================== //

	// ========================  PDF ======================== //
	pdfRepov2 := domain.NewUserResourceRepository(pdfCollection)
	pdfServicev2 := domain.NewUserResourceService(pdfRepov2, fileGateway, userGateway)
	pdfHandlerv2 := domain.NewUserResourceHandler(pdfServicev2)
	// ========================  PDF ======================== //

	topicResourceRepov2 := repository.NewTopicResourceRepository(topicResourceCollection)
	topicResourceServicev2 := service.NewTopicResourceService(topicResourceRepov2, topicRepov2, fileGateway, userGateway)
	topicResourceHandlerv2 := handler.NewTopicResourceHandler(topicResourceServicev2)
	// Register routes
	route.RegisterTopicRoutes(r, topicHandlerv2, userGateway)
	route.RegisterTopicResourceRoutes(r, topicResourceHandlerv2, userGateway)
	route2.RegisterRoutes(r, pdfHandlerv2, userGateway)
	return r
}
