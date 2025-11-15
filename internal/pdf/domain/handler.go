package domain

import (
	"fmt"
	"media-service/helper"
	"media-service/internal/pdf/domain/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserResourceHandler struct {
	userResourceService UserResourceService
}

func NewUserResourceHandler(userResourceService UserResourceService) *UserResourceHandler {
	return &UserResourceHandler{
		userResourceService: userResourceService,
	}
}

func (h *UserResourceHandler) CreateResource(c *fiber.Ctx) error {

	var req dto.CreateResourceRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	res, err := h.userResourceService.CreateResource(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *UserResourceHandler) GetResources(c *fiber.Ctx) error {

	role := c.Query("role")
	if role == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("role is required"), helper.ErrInvalidRequest)
	}

	orgID := c.Query("organization_id")
	if orgID == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("organization_id is required"), helper.ErrInvalidRequest)
	}

	res, err := h.userResourceService.GetResources(c.UserContext(), role, orgID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get pdf success", res)

}

func (h *UserResourceHandler) UploadDocumentToResource(c *fiber.Ctx) error {

	var req dto.UpdateResourceRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}

	res, err := h.userResourceService.UploadDocumentToResource(c.UserContext(), id, req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *UserResourceHandler) UploadSignatureToResource(c *fiber.Ctx) error {

	var req dto.UploadSignatureRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}

	res, err := h.userResourceService.UploadSignatureToResource(c.UserContext(), id, req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *UserResourceHandler) UpdateResourceStatus(c *fiber.Ctx) error {
	var req dto.UpdateResourceStatusRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}

	err := h.userResourceService.UpdateResourceStatus(c.UserContext(), id, req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "update resource status success", nil)
}

func (h *UserResourceHandler) UpdateResourceDownloadPermission(c *fiber.Ctx) error {
	var req dto.UpdateDownloadPermissionRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}

	err := h.userResourceService.UpdateResourceDownloadPermission(c.UserContext(), id, req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "update resource download permission success", nil)
}

func (h *UserResourceHandler) DeleteResource(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
	}

	err := h.userResourceService.DeleteResource(c.UserContext(), id)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "delete pdf success", nil)

}
