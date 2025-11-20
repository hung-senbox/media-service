package router

import (
	"media-service/internal/gateway"
	"media-service/internal/media/route"
	"media-service/internal/media/v2/handler"
	"media-service/internal/media/v2/repository"
	"media-service/internal/media/v2/service"
	"media-service/internal/media/v2/usecase"
	mediaassetHandler "media-service/internal/mediaasset/handler"
	mediaassetRepo "media-service/internal/mediaasset/repository"
	mediaassetRoute "media-service/internal/mediaasset/route"
	mediaassetService "media-service/internal/mediaasset/service"
	"media-service/internal/middleware"
	"media-service/internal/pdf/domain"
	route2 "media-service/internal/pdf/route"
	"media-service/internal/redis"
	s3svc "media-service/internal/s3"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/hashicorp/consul/api"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache"
	"github.com/hung-senbox/senbox-cache-service/pkg/cache/cached"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(app *fiber.App, consulClient *api.Client, cacheClientRedis *cache.RedisCache, topicCollection, pdfCollection, topicResourceCollection, videoUploaderCollection, mediaAssetCollection, vocabularyCollection *mongo.Collection) *fiber.App {

	app.Use(fiberLogger.New())
	// Apply CORS for all routes
	app.Use(middleware.CORS())

	// gateway
	cachedMainGateway := cached.NewCachedMainGateway(cacheClientRedis)
	userGateway := gateway.NewUserGateway("go-main-service", consulClient, cachedMainGateway)
	redisService := redis.NewRedisService()

	// ========================  Topic ======================== //
	// --- Repo ---
	topicRepov2 := repository.NewTopicRepository(topicCollection)
	topicResourceRepov2 := repository.NewTopicResourceRepository(topicResourceCollection)
	vocabularyRepo := repository.NewVocabularyRepository(vocabularyCollection)

	// --- UseCase ---
	uploadTopicUseCasev2 := usecase.NewUploadTopicUseCase(topicRepov2, s3svc.NewFromConfig())
	getTopicAppUseCasev2 := usecase.NewGetTopicAppUseCase(topicRepov2, userGateway, s3svc.NewFromConfig())
	getTopicWebUseCasev2 := usecase.NewGetTopicWebUseCase(topicRepov2, topicResourceRepov2, userGateway, s3svc.NewFromConfig())
	getTopicGatewayUseCasev2 := usecase.NewGetTopicGatewayUseCase(topicRepov2, userGateway, s3svc.NewFromConfig())
	getUploadProgressUseCasev2 := usecase.NewGetUploadProgressUseCase(topicRepov2, redisService)
	deleteTopicFileUseCasev2 := usecase.NewDeleteTopicFileUseCase(topicRepov2, s3svc.NewFromConfig())
	getTopicResourcesWebUseCasev2 := usecase.NewGetTopicResourcesWebUseCase(topicResourceRepov2, s3svc.NewFromConfig())
	getTopicResourceAppUseCasev2 := usecase.NewGetTopicResourceAppUseCase(topicRepov2, topicResourceRepov2, s3svc.NewFromConfig())
	uploadVocabularyUseCase := usecase.NewUploadVocabularyUseCase(topicRepov2, vocabularyRepo, s3svc.NewFromConfig())
	getVocabularyWebUseCase := usecase.NewGetVocabularyWebUseCase(vocabularyRepo, s3svc.NewFromConfig())

	// --- Service ---
	topicServicev2 := service.NewTopicService(uploadTopicUseCasev2, getUploadProgressUseCasev2, getTopicAppUseCasev2, getTopicWebUseCasev2, getTopicGatewayUseCasev2, deleteTopicFileUseCasev2)
	vocabularyService := service.NewVocabularyService(uploadVocabularyUseCase, getVocabularyWebUseCase)

	// --- Handler ---
	topicHandlerv2 := handler.NewTopicHandler(topicServicev2)
	vocabularyHandler := handler.NewVocabularyHandler(vocabularyService)
	// ========================  Topic ======================== //

	// ========================  PDF ======================== //
	pdfRepov2 := domain.NewUserResourceRepository(pdfCollection)
	pdfServicev2 := domain.NewUserResourceService(pdfRepov2, s3svc.NewFromConfig(), userGateway)
	pdfHandlerv2 := domain.NewUserResourceHandler(pdfServicev2)
	// ========================  PDF ======================== //

	topicResourceServicev2 := service.NewTopicResourceService(topicResourceRepov2, topicRepov2, s3svc.NewFromConfig(), userGateway, getTopicResourcesWebUseCasev2, getTopicResourceAppUseCasev2)
	topicResourceHandlerv2 := handler.NewTopicResourceHandler(topicResourceServicev2)

	// ========================  Video Uploader ======================== //
	videoUploaderRepo := repository.NewVideoUploaderRepository(videoUploaderCollection)
	videoUploaderService := service.NewVideoUploaderService(videoUploaderRepo, s3svc.NewFromConfig(), userGateway)
	videoUploaderHandler := handler.NewVideoUploaderHandler(videoUploaderService)
	// ========================  Video Uploader ======================== //

	// Register routes
	route.RegisterTopicRoutes(app, topicHandlerv2, vocabularyHandler, userGateway)
	route.RegisterTopicResourceRoutes(app, topicResourceHandlerv2, userGateway)
	route.RegisterVideoUploaderRoutes(app, videoUploaderHandler, userGateway)
	route2.RegisterRoutes(app, pdfHandlerv2, userGateway)

	// ========================  Media Assets (direct S3) ======================== //
	mediaRepo := mediaassetRepo.NewMediaRepository(mediaAssetCollection)
	mediaSvc := mediaassetService.NewMediaService(mediaRepo)
	mediaHandler := mediaassetHandler.NewMediaHandler(mediaSvc)
	mediaassetRoute.RegisterMediaRoutes(app, mediaHandler)
	return app
}
