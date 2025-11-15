package route

import (
	"media-service/internal/gateway"
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterTopicResourceRoutes(app *fiber.App, h *handler.TopicResourceHandler, userGw gateway.UserGateway) {
	topicResourceGroup := app.Group("/api/v2/topic-resources")
	topicResourceGroup.Use(middleware.Secured(userGw))
	
	topicResourceGroup.Post("", h.CreateTopicResource)
	topicResourceGroup.Get("", h.GetTopicResources)
	topicResourceGroup.Get("/:topic_resource_id", h.GetTopicResource)
	topicResourceGroup.Put("/:topic_resource_id", h.UpdateTopicResource)
	topicResourceGroup.Delete("/:topic_resource_id", h.DeleteTopicResource)

	adminGroup := app.Group("/api/v2/admin")
	adminGroup.Use(middleware.Secured(userGw))
	
	topicResourceAdmin := adminGroup.Group("/resources")
	topicResourceAdmin.Get("/topic/:topic_id", h.GetTopicResourcesByTopic4Web)
	topicResourceAdmin.Get("/topic/:topic_id/student/:student_id", h.GetTopicResourcesByTopicAndStudent4Web)
	topicResourceAdmin.Post("/output", h.SetOutputTopicResource)
	topicResourceAdmin.Delete("/output/:topic_resource_id", h.OffOutputTopicResource)
	topicResourceAdmin.Get("/output/topic/:topic_id/student/:student_id", h.GetOutputResources4Web)

	userGroup := app.Group("/api/v2/user")
	userGroup.Use(middleware.Secured(userGw))
	
	topicResourceUser := userGroup.Group("/resources")
	topicResourceUser.Get("/output/student/:student_id", h.GetOutputResources4App)
}
