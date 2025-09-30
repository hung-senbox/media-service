package handler

import (
	"errors"
	"fmt"
	"media-service/helper"
	"media-service/internal/department/dto/request"
	"media-service/internal/department/service"
	"media-service/pkg/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(s service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: s}
}

func (h *DepartmentHandler) UploadDepartment(c *gin.Context) {
	var req request.UploadDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.UploadDepartment(c.Request.Context(), req)

	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Department uploaded successfully", res)
}

func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	var req request.UpdateDepartmentRequest

	// lấy id từ param
	req.ID = c.Param("department_id")

	if req.ID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), helper.ErrInvalidRequest)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.UpdateDepartment(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Department updated successfully", res)
}

func (h *DepartmentHandler) UploadDepartmentMenu(c *gin.Context) {
	var req request.UploadSectionMenuDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.UploadDepartmentMenu(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Department menu uploaded successfully", nil)
}

func (h *DepartmentHandler) GetDepartments4Web(c *gin.Context) {
	res, err := h.service.GetDepartments4Web(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get department successfully", res)
}

func (h *DepartmentHandler) GetDepartmentDetail4Web(c *gin.Context) {
	// lấy id từ param
	departmentID := c.Param("department_id")

	if departmentID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), helper.ErrInvalidRequest)
		return
	}
	res, err := h.service.GetDeparmentDetail4Web(c.Request.Context(), departmentID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get department successfully", res)
}

func (h *DepartmentHandler) AssignLeader(c *gin.Context) {
	var req request.AssignLeaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	// check owner role valid
	if !constants.OwnerRole(req.OwnerRole).IsValid() {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid owner role: %s", req.OwnerRole), helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.AssignLeader(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Assign leader successfully", res)
}

func (h *DepartmentHandler) AssignStaff(c *gin.Context) {
	var req request.AssignStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	// check owner role valid
	if !constants.OwnerRole(req.OwnerRole).IsValid() {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("invalid owner role: %s", req.OwnerRole), helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.AssignStaff(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Assign staff successfully", res)
}

func (h *DepartmentHandler) RemoveStaff(c *gin.Context) {
	var req request.RemoveStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.RemoveStaffByIndex(c.Request.Context(), req.DepartmentID, req.Index)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Remove staff successfully", nil)
}

func (h *DepartmentHandler) GetDepartments4App(c *gin.Context) {
	res, err := h.service.GetDepartments4App(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get department successfully", res)
}

func (h *DepartmentHandler) GetDepartments4Gateway(c *gin.Context) {
	res, err := h.service.GetDepartments4Gateway(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get department successfully", res)
}

func (h *DepartmentHandler) UploadDepartmentMenuOrganization(c *gin.Context) {
	var req request.UploadDepartmentMenuOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.UploadDepartmentMenuOrganization(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Department menu uploaded successfully", nil)
}

func (h *DepartmentHandler) RemoveLeader(c *gin.Context) {
	var req request.RemoveLeaderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.RemoveLeader(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Remove leader successfully", nil)
}

func (h *DepartmentHandler) GetDepartmentsByOrganization4Gateway(c *gin.Context) {
	organizationID := c.Param("organization_id")
	if organizationID == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("organization_id is required"), helper.ErrInvalidRequest)
		return
	}

	res, err := h.service.GetDepartmentsByOrganization4Gateway(c.Request.Context(), organizationID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Get department successfully", res)
}
