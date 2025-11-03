package route

import (
	"media-service/internal/gateway"
	"media-service/internal/middleware"
	"media-service/internal/pdf/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine, userResource *domain.UserResourceHandler, userGw gateway.UserGateway) {
	pdfGroup := route.Group("/api/v1/user", middleware.Secured(userGw))
	{
		pdfGroup.POST("/resource", userResource.CreateResource)
		pdfGroup.GET("/resource", userResource.GetResources)
		pdfGroup.DELETE("/resource/:id", userResource.DeleteResource)
		pdfGroup.PUT("/resource/:id", userResource.UploadDocumentToResource)
		pdfGroup.PUT("/resource/add-signature/:id", userResource.UploadSignatureToResource)
		pdfGroup.PUT("/resource/update-status/:id", userResource.UpdateResourceStatus)
		pdfGroup.PUT("/resource/update-download-permission/:id", userResource.UpdateResourceDownloadPermission)
	}
}
