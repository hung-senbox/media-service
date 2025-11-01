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
}
