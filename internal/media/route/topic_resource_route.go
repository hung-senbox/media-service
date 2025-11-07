package route

import (
	"media-service/internal/gateway"
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTopicResourceRoutes(r *gin.Engine, h *handler.TopicResourceHandler, userGw gateway.UserGateway) {
	topicResourceGroup := r.Group("/api/v2/topic-resources")
	topicResourceGroup.Use(middleware.Secured(userGw))
	{
		topicResourceGroup.POST("", h.CreateTopicResource)
		topicResourceGroup.GET("", h.GetTopicResources)
		topicResourceGroup.GET("/:topic_resource_id", h.GetTopicResource)
		topicResourceGroup.PUT("/:topic_resource_id", h.UpdateTopicResource)
		topicResourceGroup.DELETE("/:topic_resource_id", h.DeleteTopicResource)
	}

	adminGroup := r.Group("/api/v2/admin")
	adminGroup.Use(middleware.Secured(userGw))
	{
		topicResourceAdmin := adminGroup.Group("/resources")
		{
			topicResourceAdmin.GET("/topic/:topic_id", h.GetTopicResourcesByTopic4Web)
			topicResourceAdmin.GET("/topic/:topic_id/student/:student_id", h.GetTopicResourcesByTopicAndStudent4Web)
			topicResourceAdmin.POST("/output", h.SetOutputTopicResource)
			topicResourceAdmin.DELETE("/output/:topic_resource_id", h.OffOutputTopicResource)
			topicResourceAdmin.GET("/output/topic/:topic_id/student/:student_id", h.GetOutputResources4Web)
		}
	}

	userGroup := r.Group("/api/v2/user")
	userGroup.Use(middleware.Secured(userGw))
	{
		topicResourceUser := userGroup.Group("/resources")
		{
			topicResourceUser.GET("/output/student/:student_id", h.GetOutputResources4App)
		}
	}
}
