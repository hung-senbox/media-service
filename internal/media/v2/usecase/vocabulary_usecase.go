package usecase

import (
	"context"
	"media-service/helper"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	s3svc "media-service/internal/s3"
)

type VocabularyUseCase interface {
	GetVocabulariesByTopicID4App(ctx context.Context, topicID string) ([]*response.GetVocabularyResponse4App, error)
}

type vocabularyUseCase struct {
	vocabularyRepo repository.VocabularyRepository
	s3Service      s3svc.Service
}

func NewVocabularyUseCase(vocabularyRepo repository.VocabularyRepository, s3Service s3svc.Service) VocabularyUseCase {
	return &vocabularyUseCase{
		vocabularyRepo: vocabularyRepo,
		s3Service:      s3Service,
	}
}

func (uc *vocabularyUseCase) GetVocabulariesByTopicID4App(ctx context.Context, topicID string) ([]*response.GetVocabularyResponse4App, error) {
	appLanguage := helper.GetAppLanguage(ctx, 1)
	vocabularies, err := uc.vocabularyRepo.GetAllVocabulariesByTopicIDAndIsPublished(ctx, topicID)
	if err != nil {
		return nil, err
	}
	for vi := range vocabularies {
		// duyệt qua từng language config
		for li := range vocabularies[vi].LanguageConfig {
			langCfg := &vocabularies[vi].LanguageConfig[li]

			// images
			for ii := range langCfg.Images {
				img := &langCfg.Images[ii]
				if img.ImageKey != "" {
					url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
					if err == nil && url != nil {
						img.UploadedUrl = *url
					}
				}
			}
		}
	}
	return mapper.ToVocabulariesResponses4App(vocabularies, appLanguage), nil
}
