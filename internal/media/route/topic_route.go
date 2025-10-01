package route

import (
	"media-service/internal/media/handler"
	"media-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTopicRoutes(r *gin.Engine, h *handler.TopicHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured())
	{
		topicsAdmin := adminGroup.Group("/topics")
		{
			topicsAdmin.POST("", h.UploadTopic)
			topicsAdmin.GET("/:topic_id/progress", h.GetPregressUpload)
		}
	}
}
