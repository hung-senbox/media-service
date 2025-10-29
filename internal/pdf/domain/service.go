package domain

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/internal/pdf/domain/dto"
	"media-service/internal/pdf/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResourceService interface {
	CreateResource(ctx context.Context, req dto.CreateResource) (string, error)
	// GetPDFsByStudent(ctx context.Context, studentID string) ([]*dto.StudentReportPDFResponse, error)
	// UpdatePDFsBy(ctx context.Context, id string, req dto.UpdatePDFRequest) error
	// DeletePDFsBy(ctx context.Context, id string) error
}

type userResourceService struct {
	UserResourceRepository UserResourceRepository
	fileGateway            gateway.FileGateway
}

func NewUserResourceService(userResourceRepository UserResourceRepository,
	fileGateway gateway.FileGateway) UserResourceService {
	return &userResourceService{
		UserResourceRepository: userResourceRepository,
		fileGateway:            fileGateway,
	}
}

func (s *userResourceService) CreateResource(ctx context.Context, req dto.CreateResource) (string, error) {

	if req.UploaderID == nil {
		return "", fmt.Errorf("uploader id cannot be empty")
	}

	if req.TargetID == nil {
		return "", fmt.Errorf("target id cannot be empty")
	}

	if req.Folder == "" {
		return "", fmt.Errorf("folder cannot be empty")
	}

	ID := primitive.NewObjectID()

	err := s.UserResourceRepository.CreateResource(ctx, &model.UserResource{
		ID:           ID,
		UploaderID:   req.UploaderID,
		TargetID:     req.TargetID,
		ResourceType: "",
		FileName:     "",
		Folder:       req.Folder,
		Color:        req.Color,
		SignatureKey: "",
		URL:          nil,
		PDFKey:       nil,
		CreatedBy:    helper.GetUserID(ctx),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		return "", err
	}

	return ID.Hex(), nil

}

// func (s *pdfService) GetPDFsByStudent(ctx context.Context, studentID string) ([]*dto.StudentReportPDFResponse, error) {

// 	studentPdfs, err := s.PDFRepository.GetPDFsByStudent(ctx, studentID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []*dto.StudentReportPDFResponse

// 	for _, studentPdf := range studentPdfs {
// 		var pdfUrl *string
// 		pdfUrl, err = s.fileGateway.GetPDFUrl(ctx, gw_request.GetFileUrlRequest{
// 			Key:  studentPdf.PDFKey,
// 			Mode: "private",
// 		})

// 		if err != nil {
// 			log.Printf("error get pdf url: %s", err.Error())
// 			pdfUrl = nil
// 		}

// 		result = append(result, &dto.StudentReportPDFResponse{
// 			ID:        studentPdf.ID.Hex(),
// 			StudentID: studentPdf.StudentID,
// 			PDFName:   studentPdf.PDFName,
// 			Folder:    studentPdf.Folder,
// 			PDFUrl:    pdfUrl,
// 			Color:     studentPdf.Color,
// 			CreatedBy: studentPdf.CreatedBy,
// 			CreatedAt: studentPdf.CreatedAt,
// 			UpdatedAt: studentPdf.UpdatedAt,
// 		})
// 	}

// 	return result, nil
// }

// func (s *pdfService) UpdatePDFsBy(ctx context.Context, id string, req dto.UpdatePDFRequest) error {

// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	pdfData, err := s.PDFRepository.GetPDFByID(ctx, objectID)
// 	if err != nil {
// 		return err
// 	}

// 	if pdfData.PDFKey != "" {
// 		err = s.fileGateway.DeletePDF(ctx, pdfData.PDFKey)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if pdfData == nil {
// 		return fmt.Errorf("pdf not found")
// 	}

// 	if req.FileName != "" {
// 		pdfData.PDFName = req.FileName
// 	}

// 	if req.Color != "" {
// 		pdfData.Color = req.Color
// 	}

// 	if req.File == nil {
// 		return fmt.Errorf("file is required")
// 	}

// 	resp, err := s.fileGateway.UploadPDF(ctx, gw_request.UploadFileRequest{
// 		File:     req.File,
// 		Folder:   "pdf_media",
// 		FileName: pdfData.PDFName,
// 		Mode:     "private",
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	pdfData.PDFKey = resp.Key
// 	pdfData.UpdatedAt = time.Now()

// 	return s.PDFRepository.UpdatePDFByID(ctx, objectID, pdfData)

// }

// func (s *pdfService) DeletePDFsBy(ctx context.Context, id string) error {

// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return err
// 	}

// 	pdf, err := s.PDFRepository.GetPDFByID(ctx, objectID)
// 	if err != nil {
// 		return err
// 	}

// 	if pdf == nil {
// 		return fmt.Errorf("pdf not found")
// 	}

// 	err = s.fileGateway.DeletePDF(ctx, pdf.PDFKey)
// 	if err != nil {
// 		return err
// 	}

// 	return s.PDFRepository.DeletePDFByID(ctx, objectID)

// }
