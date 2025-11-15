package route

import (
	"media-service/internal/mediaasset/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterMediaRoutes(app *fiber.App, h *handler.MediaHandler) {
	v2 := app.Group("/v2/media")
	
	v2.Post("/upload", h.Upload)
	v2.Get("/:id/url", h.GetURL)
	v2.Get("/:id", h.GetMeta)
	v2.Delete("/:id", h.Delete)
	v2.Get("/url", h.GetURLByKey)
}
