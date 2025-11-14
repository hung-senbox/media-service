package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strings"

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

	err := h.service.UploadVideoUploader(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "upload video uploader success", nil)
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
	title := c.Query("title")

	var sortBy []request.GetVideoUploaderSortBy

	sortParam := strings.TrimSpace(c.Query("sort_by"))
	if sortParam != "" {
		items := strings.Split(sortParam, ",")
		for _, item := range items {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}

			order := request.GetVideoUploaderSortByOrderAsc
			field := item

			// Detect prefix + or -
			if strings.HasPrefix(item, "-") {
				order = request.GetVideoUploaderSortByOrderDesc
				field = strings.TrimPrefix(item, "-")
			} else if strings.HasPrefix(item, "+") {
				order = request.GetVideoUploaderSortByOrderAsc
				field = strings.TrimPrefix(item, "+")
			}

			sortBy = append(sortBy, request.GetVideoUploaderSortBy{
				Field: field,
				Order: order,
			})
		}
	}

	res, err := h.service.GetVideosUploader4Web(c.Request.Context(), languageID, title, sortBy)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get videos uploader success", res)
}

func (h *VideoUploaderHandler) DeleteVideoUploader(c *gin.Context) {
	videoUploaderID := c.Param("video_uploader_id")
	if videoUploaderID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	err := h.service.DeleteVideoUploader(c.Request.Context(), videoUploaderID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "success", nil)
}
