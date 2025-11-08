package handler

import (
	"fmt"
	"media-service/helper"
	"media-service/internal/mediaasset/dto"
	"media-service/internal/mediaasset/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	svc service.MediaService
}

func NewMediaHandler(svc service.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	folder := c.PostForm("folder")
	mode := c.PostForm("mode")
	mtStr := c.PostForm("media_type")
	var mtPtr *string
	if mtStr != "" {
		mtPtr = &mtStr
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	opened, err := fileHeader.Open()
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	_ = opened.Close()

	meta, url, err := h.svc.Upload(c.Request.Context(), fileHeader, folder, mode, mtPtr)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "upload success", dto.UploadResponse{
		ID:  meta.ID.Hex(),
		Key: meta.Key,
		URL: url,
	})
}

func (h *MediaHandler) GetURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}
	ttlStr := c.Query("ttl_seconds")
	var duration *time.Duration
	if ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr + "s"); err == nil {
			duration = &d
		}
	}
	url, err := h.svc.GetURL(c.Request.Context(), id, duration)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "ok", gin.H{"url": url})
}

func (h *MediaHandler) GetMeta(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}
	meta, err := h.svc.GetMeta(c.Request.Context(), id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "ok", meta)
}

func (h *MediaHandler) GetURLByKey(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("key is required"), helper.ErrInvalidRequest)
		return
	}
	ttlStr := c.Query("ttl_seconds")
	var duration *time.Duration
	if ttlStr != "" {
		if d, err := time.ParseDuration(ttlStr + "s"); err == nil {
			duration = &d
		}
	}
	url, err := h.svc.GetURLByKey(c.Request.Context(), key, duration)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "ok", gin.H{"url": url})
}

func (h *MediaHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "deleted", nil)
}
