package route

import (
	"media-service/internal/middleware"
	"media-service/internal/pdf/domain"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine, pdfHandler *domain.PDFHandler) {
	pdfGroup := route.Group("/api/v1/pdf", middleware.Secured())
	{
		pdfGroup.POST("/student", pdfHandler.CreatePDF)
		pdfGroup.GET("/student", pdfHandler.GetPDFsByStudent)
	}
}
