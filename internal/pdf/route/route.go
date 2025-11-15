package route

import (
	"media-service/internal/gateway"
	"media-service/internal/middleware"
	"media-service/internal/pdf/domain"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userResource *domain.UserResourceHandler, userGw gateway.UserGateway) {
	pdfGroup := app.Group("/api/v1/user", middleware.Secured(userGw))
	
	pdfGroup.Post("/resource", userResource.CreateResource)
	pdfGroup.Get("/resource", userResource.GetResources)
	pdfGroup.Delete("/resource/:id", userResource.DeleteResource)
	pdfGroup.Put("/resource/:id", userResource.UploadDocumentToResource)
	pdfGroup.Put("/resource/add-signature/:id", userResource.UploadSignatureToResource)
	pdfGroup.Put("/resource/update-status/:id", userResource.UpdateResourceStatus)
	// pdfGroup.Put("/resource/update-download-permission/:id", userResource.UpdateResourceDownloadPermission)
	pdfGroup.Put("/resource/download/:id", userResource.UpdateResourceDownloadPermission)
}
