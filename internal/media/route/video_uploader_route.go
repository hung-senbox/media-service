package route

import (
	"media-service/internal/gateway"
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterVideoUploaderRoutes(r *gin.Engine, h *handler.VideoUploaderHandler, userGw gateway.UserGateway) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured(userGw))
	{
		uploadAdmin := adminGroup.Group("/upload")
		{
			videoUploaderAdmin := uploadAdmin.Group("/videos")
			{
				videoUploaderAdmin.POST("", h.UploadVideoUploader)
				videoUploaderAdmin.GET("/progress/:video_uploader_id", h.GetUploaderStatus)
				videoUploaderAdmin.GET("", h.GetVideosUploader4Web)
				videoUploaderAdmin.DELETE("/:video_uploader_id", h.DeleteVideoUploader)
			}
		}
	}
}
