package route

import (
	"media-service/internal/gateway"
	"media-service/internal/media/v2/handler"
	"media-service/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterTopicRoutes(app *fiber.App, hv2 *handler.TopicHandler, hv *handler.VocabularyHandler, userGw gateway.UserGateway) {
	// Admin routes
	adminGroup := app.Group("/api/v2/admin")
	adminGroup.Use(middleware.Secured(userGw))

	topicsAdmin := adminGroup.Group("/topics")

	vocabularyAdmin := topicsAdmin.Group("/:topic_id/vocabularies")
	vocabularyAdmin.Get("", hv.GetVocabularies4Web)
	vocabularyAdmin.Post("", middleware.RequireAdmin(), hv.UploadVocabulary)

	topicsAdmin.Post("", middleware.RequireAdmin(), hv2.UploadTopic)
	topicsAdmin.Get("", hv2.GetTopics4Web)
	// Static routes MUST come before dynamic routes
	topicsAdmin.Get("/assign", hv2.GetTopics2Assign4Web)
	topicsAdmin.Get("/student/:student_id", hv2.GetTopics4Student4Web)
	// Dynamic routes come after static routes
	topicsAdmin.Get("/:topic_id/progress", hv2.GetPregressUpload)
	topicsAdmin.Get("/:topic_id", hv2.GetTopic4Web)
	topicsAdmin.Delete("/audio/:topic_id/language/:language_id", hv2.DeleteTopicAudioKey)
	topicsAdmin.Delete("/video/:topic_id/language/:language_id", hv2.DeleteTopicVideoKey)
	topicsAdmin.Delete("/image/:topic_id/language/:language_id/type/:image_type", hv2.DeleteTopicImageKey)

	// User routes
	userGroup := app.Group("/api/v2/user")
	userGroup.Use(middleware.Secured(userGw))

	topicsUser := userGroup.Group("/topics")
	topicsUser.Get("/student/:student_id", hv2.GetTopics4Student4App)
	topicsUser.Get("/organization/:organization_id", hv2.GetTopics4App)

	// gateway
	gatewayGroup := app.Group("/api/v2/gateway")
	gatewayGroup.Use(middleware.Secured(userGw))

	topicsGateway := gatewayGroup.Group("/topics")
	topicsGateway.Get("/organization/:organization_id", hv2.GetAllTopicsByOrganization4Gw)
	topicsGateway.Get("/:topic_id", hv2.GetTopic4Gw)
	topicsGateway.Get("/student/:student_id", hv2.GetTopics4Student4Gw)
}
