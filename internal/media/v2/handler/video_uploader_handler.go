package handler

import (
	"fmt"
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strconv"
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
	// Parse multipart form
	req := request.UploadVideoUploaderRequest{
		VideoFolderID:         c.FormValue("video_folder_id"),
		Title:                 c.FormValue("title"),
		WikiCode:              c.FormValue("wiki_code"),
		Note:                  c.FormValue("note"),
		Transcript:            c.FormValue("transcript"),
		IsVisible:             c.FormValue("is_visible") == "true",
		IsDeletedVideo:        c.FormValue("is_deleted_video") == "true",
		IsDeletedImagePreview: c.FormValue("is_deleted_image") == "true",
	}

	// Parse uint fields
	if langID := c.FormValue("language_id"); langID != "" {
		if val, err := strconv.ParseUint(langID, 10, 32); err == nil {
			req.LanguageID = uint(val)
		}
	}

	// Parse file fields
	if videoFile, err := c.FormFile("video_file"); err == nil {
		req.VideoFile = videoFile
	}
	if imagePreviewFile, err := c.FormFile("image_preview_file"); err == nil {
		req.ImagePreviewFile = imagePreviewFile
	}

	videoUploader, err := h.service.UploadVideoUploader(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "upload video uploader success", videoUploader)
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

func (h *VideoUploaderHandler) GetVideo4Web(c *fiber.Ctx) error {
	videoUploaderID := c.Params("video_uploader_id")
	if videoUploaderID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetVideo4Web(c.UserContext(), videoUploaderID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get video success", res)
}

func (h *VideoUploaderHandler) GetVideosByWikiCode4Web(c *fiber.Ctx) error {
	wikiCode := c.Params("wiki_code")
	if wikiCode == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageID := c.Query("language_id")
	if languageID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	var langID uint
	if _, err := fmt.Sscanf(languageID, "%d", &langID); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetVideosByWikiCode4Web(c.UserContext(), wikiCode, langID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get videos by wiki code success", res)
}
