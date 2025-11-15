package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type VideoUploaderHandler struct {
	service service.VideoUploaderService
}

func NewVideoUploaderHandler(service service.VideoUploaderService) *VideoUploaderHandler {
	return &VideoUploaderHandler{service: service}
}

func (h *VideoUploaderHandler) UploadVideoUploader(c *fiber.Ctx) error {
	var req request.UploadVideoUploaderRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	err := h.service.UploadVideoUploader(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "upload video uploader success", nil)
}

func (h *VideoUploaderHandler) GetUploaderStatus(c *fiber.Ctx) error {
	videoUploaderID := c.Params("video_uploader_id")
	if videoUploaderID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}

	res, err := h.service.GetUploaderStatus(c.UserContext(), videoUploaderID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get uploader status success", res)
}

func (h *VideoUploaderHandler) GetVideosUploader4Web(c *fiber.Ctx) error {
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

	res, err := h.service.GetVideosUploader4Web(c.UserContext(), languageID, title, sortBy)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get videos uploader success", res)
}

func (h *VideoUploaderHandler) DeleteVideoUploader(c *fiber.Ctx) error {
	videoUploaderID := c.Params("video_uploader_id")
	if videoUploaderID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	err := h.service.DeleteVideoUploader(c.UserContext(), videoUploaderID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "success", nil)
}
