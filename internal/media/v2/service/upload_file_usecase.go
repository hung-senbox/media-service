package service

import (
	"context"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	gw_response "media-service/internal/gateway/dto/response"
	"time"
)

type UploadFileService interface {
	UploadImage(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadImageResponse, error)
	UploadPDF(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadPDFResponse, error)
}

type uploadFileService struct {
	fileGateway gateway.FileGateway
}

func NewUploadFileService(fileGateway gateway.FileGateway) UploadFileService {
	return &uploadFileService{
		fileGateway: fileGateway,
	}
}

func (uc *uploadFileService) UploadImage(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadImageResponse, error) {
	if req.FileName == "" {
		req.FileName = time.Now().Format("20060102150405")
	}
	resp, err := uc.fileGateway.UploadImage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (uc *uploadFileService) UploadPDF(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadPDFResponse, error) {
	if req.FileName == "" {
		req.FileName = time.Now().Format("20060102150405")
	}
	resp, err := uc.fileGateway.UploadPDF(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}