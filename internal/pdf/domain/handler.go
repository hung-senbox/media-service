package domain

import (
	"media-service/helper"
	"media-service/internal/pdf/domain/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PDFHandler struct {
	PDFService PDFService
}

func NewPDFHandler(pdfService PDFService) *PDFHandler {
	return &PDFHandler{
		PDFService: pdfService,
	}
}

func (h *PDFHandler) CreatePDF(c *gin.Context) {

	var req dto.CreatePDFRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	res, err := h.PDFService.CreatePDF(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "create pdf success", res)

}

func (h *PDFHandler) GetPDFsByStudent(c *gin.Context) {

	studentID := c.Query("student_id")

	res, err := h.PDFService.GetPDFsByStudent(c.Request.Context(), studentID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get pdfs success", res)

}

func (h *PDFHandler) DeletePDFsBy(c *gin.Context) {

	id := c.Param("id")

	err := h.PDFService.DeletePDFsBy(c.Request.Context(), id)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "delete pdfs success", nil)

}

func (h *PDFHandler) UpdatePDFsBy(c *gin.Context) {

	id := c.Param("id")

	var req dto.UpdatePDFRequest

	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.PDFService.UpdatePDFsBy(c.Request.Context(), id, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "update pdfs success", nil)
} 
