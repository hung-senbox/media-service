package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoUploaderHandler struct {
	service service.VideoUploaderService
}

func NewVideoUploaderHandler(service service.VideoUploaderService) *VideoUploaderHandler {
	return &VideoUploaderHandler{service: service}
}

func (h *VideoUploaderHandler) UploadVideoUploader(c *gin.Context) {
	var req request.UploadVideoUploaderRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.UploadVideoUploader(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "upload video uploader success", res)
}

func (h *VideoUploaderHandler) GetUploaderStatus(c *gin.Context) {
	videoUploaderID := c.Param("video_uploader_id")
	if videoUploaderID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.GetUploaderStatus(c.Request.Context(), videoUploaderID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get uploader status success", res)
}

func (h *VideoUploaderHandler) GetVideosUploader4Web(c *gin.Context) {
	languageID := c.Query("language_id")

	res, err := h.service.GetVideosUploader4Web(c.Request.Context(), languageID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get videos uploader success", res)
}
