package route

import (
	"media-service/internal/gateway"
	"media-service/internal/middleware"
	"media-service/internal/pdf/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine, pdfHandler *domain.PDFHandler, userGw gateway.UserGateway) {
	pdfGroup := route.Group("/api/v1/pdf", middleware.Secured(userGw))
	{
		pdfGroup.POST("/student", pdfHandler.CreatePDF)
		pdfGroup.GET("/student", pdfHandler.GetPDFsByStudent)
		pdfGroup.DELETE("/:id", pdfHandler.DeletePDFsBy)
		pdfGroup.PUT("/:id", pdfHandler.UpdatePDFsBy)
	}
}
