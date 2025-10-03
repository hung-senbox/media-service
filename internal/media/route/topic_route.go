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
			topicsAdmin.GET("", hv2.GetTopics4Web)
			topicsAdmin.GET("/:topic_id/progress", hv2.GetPregressUpload)
			topicsAdmin.GET("/:topic_id", hv2.GetTopic4Web)
		}
	}

	// User routes
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(middleware.Secured())
	{
		topicsUser := userGroup.Group("/topics")
		{
			topicsUser.GET("/student/:student_id", hv2.GetTopics4Student)
		}
	}
}
