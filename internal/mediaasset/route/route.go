package route

import (
	"media-service/internal/mediaasset/handler"

	"github.com/gin-gonic/gin"
)

func RegisterMediaRoutes(r *gin.Engine, h *handler.MediaHandler) {
	v2 := r.Group("/v2/media")
	{
		v2.POST("/upload", h.Upload)
		v2.GET("/:id/url", h.GetURL)
		v2.GET("/:id", h.GetMeta)
		v2.DELETE("/:id", h.Delete)
		v2.GET("/url", h.GetURLByKey)
	}
}
