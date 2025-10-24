package route

import (
	"media-service/internal/gateway"
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTopicRoutes(r *gin.Engine, hv2 *handler.TopicHandler, userGw gateway.UserGateway) {
	// Admin routes
	adminGroup := r.Group("/api/v2/admin")
	adminGroup.Use(middleware.Secured(userGw))
	{
		topicsAdmin := adminGroup.Group("/topics")
		{
			topicsAdmin.POST("", hv2.CreateTopic)
			topicsAdmin.GET("", hv2.GetTopics4Web)
			topicsAdmin.GET("/:topic_id/progress", hv2.GetPregressUpload)
			topicsAdmin.GET("/:topic_id", hv2.GetTopic4Web)
			topicsAdmin.GET("/student/:student_id", hv2.GetTopics4Student4Web)
			topicsAdmin.PUT("/:topic_id", hv2.UpdateTopic)
			topicsAdmin.GET("/assign", hv2.GetTopics2Assign4Web)
		}
	}

	// User routes
	userGroup := r.Group("/api/v2/user")
	userGroup.Use(middleware.Secured(userGw))
	{
		topicsUser := userGroup.Group("/topics")
		{
			topicsUser.GET("/student/:student_id", hv2.GetTopics4Student4App)
		}
	}

	// gateway
	gatewayGroup := r.Group("/api/v2/gateway")
	gatewayGroup.Use(middleware.Secured(userGw))
	{
		topicsGateway := gatewayGroup.Group("/topics")
		{
			topicsGateway.GET("/:topic_id", hv2.GetTopic4GW)
			topicsGateway.GET("/student/:student_id", hv2.GetTopics4Student4Gw)
		}
	}
}
