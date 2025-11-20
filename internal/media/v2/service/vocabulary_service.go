package service

import (
	"context"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/usecase"
)

type VocabularyService interface {
	UploadVocabulary(ctx context.Context, req request.UploadVocabularyRequest) error
	GetVocabularies4Web(ctx context.Context, topicID string) (*response.VocabularyResponse4Web, error)
}

type vocabularyService struct {
	uploadVocabularyUseCase usecase.UploadVocabularyUseCase
	getVocabularyWebUseCase usecase.GetVocabularyWebUseCase
}

func NewVocabularyService(uploadVocabularyUseCase usecase.UploadVocabularyUseCase, getVocabularyWebUseCase usecase.GetVocabularyWebUseCase) VocabularyService {
	return &vocabularyService{
		uploadVocabularyUseCase: uploadVocabularyUseCase,
		getVocabularyWebUseCase: getVocabularyWebUseCase,
	}
}

func (s *vocabularyService) UploadVocabulary(ctx context.Context, req request.UploadVocabularyRequest) error {
	return s.uploadVocabularyUseCase.UploadVocabulary(ctx, req)
}

func (s *vocabularyService) GetVocabularies4Web(ctx context.Context, topicID string) (*response.VocabularyResponse4Web, error) {
	return s.getVocabularyWebUseCase.GetVocabularies4Web(ctx, topicID)
}
