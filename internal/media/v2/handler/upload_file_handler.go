package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/service"
	"net/http"

	gw_request "media-service/internal/gateway/dto/request"

	"github.com/gofiber/fiber/v2"
)

type UploadFileHandler struct {
	uploadFileUsecase service.UploadFileService
}

func NewUploadFileHandler(uploadFileUsecase service.UploadFileService) *UploadFileHandler {
	return &UploadFileHandler{
		uploadFileUsecase: uploadFileUsecase,
	}
}

func (h *UploadFileHandler) UploadImage(c *fiber.Ctx) error {
	_, err := c.MultipartForm()
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	folder := c.FormValue("folder")
	fileName := c.FormValue("file_name")
	mode := c.FormValue("mode")

	req := gw_request.UploadFileRequest{
		File:      file,
		Folder:    folder,
		FileName:  fileName,
		ImageName: c.FormValue("image_name"),
		Mode:      mode,
	}

	resp, err := h.uploadFileUsecase.UploadImage(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "upload image success", resp)
}

func (h *UploadFileHandler) UploadPDF(c *fiber.Ctx) error {
	_, err := c.MultipartForm()
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	
	folder := c.FormValue("folder")
	fileName := c.FormValue("file_name")
	mode := c.FormValue("mode")

	req := gw_request.UploadFileRequest{
		File:      file,
		Folder:    folder,
		FileName:  fileName,
		Mode:      mode,
	}

	resp, err := h.uploadFileUsecase.UploadPDF(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	
	return helper.SendSuccess(c, http.StatusOK, "upload pdf success", resp)
}