package handler

import (
	"errors"
	"media-service/helper"
	"media-service/internal/department/dto/request"
	"media-service/internal/department/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegionHandler struct {
	service service.RegionService
}

func NewRegionHandler(s service.RegionService) *RegionHandler {
	return &RegionHandler{service: s}
}

func (h *RegionHandler) CreateRegion(c *gin.Context) {
	var req request.CreateRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.service.CreateRegion(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Region created successfully", nil)
}

func (h *RegionHandler) UpdateRegionName(c *gin.Context) {
	var req request.UpdateRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	regionId := c.Param("region_id")
	if regionId == "" {
		helper.SendError(c, http.StatusBadRequest, errors.New("id is required"), helper.ErrInvalidRequest)
		return
	}
	req.ID = regionId

	err := h.service.UpdateRegionName(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Region updated successfully", nil)
}
