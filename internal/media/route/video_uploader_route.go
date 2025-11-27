package route

import (
	"media-service/internal/gateway"
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterVideoUploaderRoutes(app *fiber.App, h *handler.VideoUploaderHandler, userGw gateway.UserGateway) {
	// Admin routes
	adminGroup := app.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured(userGw))

	uploadAdmin := adminGroup.Group("/upload")
	videoUploaderAdmin := uploadAdmin.Group("/video_folders")

	videoUploaderAdmin.Post("", h.UploadVideoUploader)
	videoUploaderAdmin.Get("", h.GetVideosUploader4Web)
	videoUploaderAdmin.Delete("/:video_uploader_id", h.DeleteVideoUploader)
	videoUploaderAdmin.Get("/:video_uploader_id", h.GetVideo4Web)
	videoUploaderAdmin.Get("/wiki_code/:wiki_code", h.GetVideosByWikiCode4Web)

	// gateway routes
	gatewayGroup := app.Group("/api/v1/gateway")
	gatewayGroup.Use(middleware.Secured(userGw))
	uploadGateway := gatewayGroup.Group("/upload")
	videoUploaderGateway := uploadGateway.Group("/video_folders")

	videoUploaderGateway.Get("/:video_uploader_id", h.GetVideo4Gw)
}
