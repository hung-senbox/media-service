package route

import (
	"media-service/internal/gateway"
	"media-service/internal/middleware"
	"media-service/internal/pdf/domain"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userResource *domain.UserResourceHandler, userGw gateway.UserGateway) {
	// Authenticated routes - require valid JWT token
	pdfGroup := app.Group("/api/v1/user", middleware.Secured(userGw))

	// User resource routes (authenticated users)
	pdfGroup.Post("/resource", userResource.CreateResource)
	pdfGroup.Get("/resource", userResource.GetResources)
	pdfGroup.Delete("/resource/:id", userResource.DeleteResource)
	pdfGroup.Put("/resource/:id", userResource.UploadDocumentToResource)
	pdfGroup.Put("/resource/add-signature/:id", userResource.UploadSignatureToResource)

	// Admin only routes - require SuperAdmin role
	adminGroup := pdfGroup.Group("", middleware.RequireAdmin())
	adminGroup.Put("/resource/update-status/:id", userResource.UpdateResourceStatus)
	adminGroup.Put("/resource/download/:id", userResource.UpdateResourceDownloadPermission)

	// Commented out - might be used in future
	// pdfGroup.Put("/resource/update-download-permission/:id", userResource.UpdateResourceDownloadPermission)
}
