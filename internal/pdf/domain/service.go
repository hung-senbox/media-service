package domain

import (
	"context"
	"fmt"
	"log"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	"media-service/internal/pdf/domain/dto"
	"media-service/internal/pdf/model"
	"media-service/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PDFService interface {
	CreatePDF(ctx context.Context, req dto.CreatePDFRequest) (string, error)
	GetPDFsByStudent(ctx context.Context, studentID string) ([]*dto.StudentReportPDFResponse, error)
}

type pdfService struct {
	PDFRepository PDFRepository
	fileGateway   gateway.FileGateway
}

func NewPDFService(pdfRepository PDFRepository,
	fileGateway gateway.FileGateway) PDFService {
	return &pdfService{
		PDFRepository: pdfRepository,
		fileGateway:   fileGateway,
	}
}

func (s *pdfService) CreatePDF(ctx context.Context, req dto.CreatePDFRequest) (string, error) {

	if req.File == nil {
		return "", fmt.Errorf("file cannot be empty")
	}

	if req.StudentID == "" {
		return "", fmt.Errorf("student ID cannot be empty")
	}

	if req.FileName == "" {
		return "", fmt.Errorf("file name cannot be empty")
	}

	resp, err := s.fileGateway.UploadPDF(ctx, gw_request.UploadFileRequest{
		File:     req.File,
		Folder:   "pdf_media",
		FileName: req.FileName + "_pdf",
		Mode:     "private",
	})
	if err != nil {
		return "", err
	}

	userID := ctx.Value(constants.UserID).(string)
	if userID == "" {
		return "", fmt.Errorf("user ID cannot be empty")
	}

	err = s.PDFRepository.CreatePDF(ctx, &model.StudentReportPDF{
		ID:        primitive.NewObjectID(),
		StudentID: req.StudentID,
		PDFName:   req.FileName + "_pdf",
		Folder:    "pdf_media",
		PDFKey:    resp.Key,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return "", err
	}

	return resp.Key, nil

}

func (s *pdfService) GetPDFsByStudent(ctx context.Context, studentID string) ([]*dto.StudentReportPDFResponse, error) {

	studentPdfs, err := s.PDFRepository.GetPDFsByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	var result []*dto.StudentReportPDFResponse

	for _, studentPdf := range studentPdfs {
		var pdfUrl *string
		pdfUrl, err = s.fileGateway.GetPDFUrl(ctx, gw_request.GetFileUrlRequest{
			Key: studentPdf.PDFKey,
			Mode: "private",
		})

		if err != nil {
			log.Printf("error get pdf url: %s", err.Error())
			pdfUrl = nil
		}

		result = append(result, &dto.StudentReportPDFResponse{
			ID:        studentPdf.ID.Hex(),
			StudentID: studentPdf.StudentID,
			PDFName:   studentPdf.PDFName,
			Folder:    studentPdf.Folder,
			PDFUrl:    pdfUrl,
			CreatedBy: studentPdf.CreatedBy,
			CreatedAt: studentPdf.CreatedAt,
			UpdatedAt: studentPdf.UpdatedAt,
		})
	}

	return result, nil
}
