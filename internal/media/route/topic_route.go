package route

import (
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTopicRoutes(r *gin.Engine, hv2 *handler.TopicHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v2/admin")
	adminGroup.Use(middleware.Secured())
	{
		topicsAdmin := adminGroup.Group("/topics")
		{
			topicsAdmin.POST("", hv2.CreateTopic)
			topicsAdmin.GET("/parents", hv2.GetParentTopics4Web)
			topicsAdmin.GET("/:topic_id/progress", hv2.GetPregressUpload)
		}
	}
}
