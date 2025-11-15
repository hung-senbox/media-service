package handler

import (
	"fmt"
	"media-service/helper"
	"media-service/internal/mediaasset/dto"
	"media-service/internal/mediaasset/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	svc service.MediaService
}

func NewMediaHandler(svc service.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

func (h *MediaHandler) Upload(c *fiber.Ctx) error {
	folder := c.FormValue("folder")
	mode := c.FormValue("mode")
	mtStr := c.FormValue("media_type")
	var mtPtr *string
	if mtStr != "" {
		mtPtr = &mtStr
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	opened, err := fileHeader.Open()
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	_ = opened.Close()

	meta, url, err := h.svc.Upload(c.UserContext(), fileHeader, folder, mode, mtPtr)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "upload success", dto.UploadResponse{
		ID:  meta.ID.Hex(),
		Key: meta.Key,
		URL: url,
	})
}

func (h *MediaHandler) GetURL(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}
	ttlStr := c.Query("ttl_seconds")
	var duration *time.Duration
	if ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr + "s"); err == nil {
			duration = &d
		}
	}
	url, err := h.svc.GetURL(c.UserContext(), id, duration)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "ok", fiber.Map{"url": url})
}

func (h *MediaHandler) GetMeta(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}
	meta, err := h.svc.GetMeta(c.UserContext(), id)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "ok", meta)
}

func (h *MediaHandler) GetURLByKey(c *fiber.Ctx) error {
	key := c.Query("key")
	if key == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("key is required"), helper.ErrInvalidRequest)
	}
	ttlStr := c.Query("ttl_seconds")
	var duration *time.Duration
	if ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr + "s"); err == nil {
			duration = &d
		}
	}
	url, err := h.svc.GetURLByKey(c.UserContext(), key, duration)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "ok", fiber.Map{"url": url})
}

func (h *MediaHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}
	if err := h.svc.Delete(c.UserContext(), id); err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "deleted", nil)
}
