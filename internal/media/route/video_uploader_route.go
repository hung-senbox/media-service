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
	videoUploaderAdmin := uploadAdmin.Group("/videos")

	videoUploaderAdmin.Post("", h.UploadVideoUploader)
	videoUploaderAdmin.Get("/progress/:video_uploader_id", h.GetUploaderStatus)
	videoUploaderAdmin.Get("", h.GetVideosUploader4Web)
	videoUploaderAdmin.Delete("/:video_uploader_id", h.DeleteVideoUploader)
}
