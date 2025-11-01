package domain

import (
	"fmt"
	"media-service/helper"
	"media-service/internal/pdf/domain/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserResourceHandler struct {
	userResourceService UserResourceService
}

func NewUserResourceHandler(userResourceService UserResourceService) *UserResourceHandler {
	return &UserResourceHandler{
		userResourceService: userResourceService,
	}
}

func (h *UserResourceHandler) CreateResource(c *gin.Context) {

	var req dto.CreateResourceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	res, err := h.userResourceService.CreateResource(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *UserResourceHandler) GetResources(c *gin.Context) {

	role := c.Query("role")
	if role == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("role is required"), helper.ErrInvalidRequest)
		return
	}

	orgID := c.Query("organization_id")
	if orgID == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("organization_id is required"), helper.ErrInvalidRequest)
		return
	}

	res, err := h.userResourceService.GetResources(c.Request.Context(), role, orgID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get pdf success", res)

}

func (h *UserResourceHandler) UploadDocumentToResource(c *gin.Context) {

	var req dto.UpdateResourceRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	res, err := h.userResourceService.UploadDocumentToResource(c.Request.Context(), id, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *UserResourceHandler) UploadSignatureToResource(c *gin.Context) {

	var req dto.UploadSignatureRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	res, err := h.userResourceService.UploadSignatureToResource(c.Request.Context(), id, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *UserResourceHandler) DeleteResource(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	err := h.userResourceService.DeleteResource(c.Request.Context(), id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "delete pdf success", nil)

}
