package usecase

import (
	"context"
	"fmt"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	s3svc "media-service/internal/s3"
	"media-service/logger"
	"media-service/pkg/constants"
)

type GetVocabularyWebUseCase interface {
	GetVocabularies4Web(ctx context.Context, topicID string) ([]*response.VocabularyResponse4Web, error)
}

type getVocabularyWebUseCase struct {
	vocabularyRepo repository.VocabularyRepository
	s3Service      s3svc.Service
}

func NewGetVocabularyWebUseCase(vocabularyRepo repository.VocabularyRepository, s3Service s3svc.Service) GetVocabularyWebUseCase {
	return &getVocabularyWebUseCase{vocabularyRepo: vocabularyRepo, s3Service: s3Service}
}

func (uc *getVocabularyWebUseCase) GetVocabularies4Web(ctx context.Context, topicID string) ([]*response.VocabularyResponse4Web, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser == nil {
		return nil, fmt.Errorf("access denied")
	}

	if currentUser.IsSuperAdmin {
		vocabularies, err := uc.vocabularyRepo.GetAllVocabulariesByTopicID(ctx, topicID)
		if err != nil {
			return nil, err
		}
		for ti := range vocabularies {
			uc.populateMediaUrlsForVocabulary(ctx, &vocabularies[ti])
		}
		return mapper.ToVocabulariesResponses4Web(vocabularies), nil
	}

	if currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}

	vocabularies, err := uc.vocabularyRepo.GetAllVocabulariesByTopicID(ctx, topicID)
	if err != nil {
		return nil, err
	}

	// filter is published
	vocabularies = filterIsPublished(vocabularies)

	for ti := range vocabularies {
		uc.populateMediaUrlsForVocabulary(ctx, &vocabularies[ti])
	}

	return mapper.ToVocabulariesResponses4Web(vocabularies), nil
}

// populateMediaUrlsForTopic enriches a topic's language configs with signed media URLs when keys exist
func (uc *getVocabularyWebUseCase) populateMediaUrlsForVocabulary(ctx context.Context, vocabulary *model.Vocabulary) {
	for li := range vocabulary.LanguageConfig {
		langCfg := &vocabulary.LanguageConfig[li]

		// images
		for ii := range langCfg.Images {
			img := &langCfg.Images[ii]
			if img.ImageKey != "" {
				url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
				if err == nil && url != nil {
					img.UploadedUrl = *url
				} else {
					logger.WriteLogEx("get_vocabulary_web_usecase", "populateMediaUrlsForVocabulary_images", fmt.Sprintf("error getting image url: %v", err))
				}
			}
		}

		// video
		if langCfg.Video.VideoKey != "" {
			url, err := uc.s3Service.Get(ctx, langCfg.Video.VideoKey, nil)
			if err == nil && url != nil {
				langCfg.Video.UploadedUrl = *url
			} else {
				logger.WriteLogEx("get_vocabulary_web_usecase", "populateMediaUrlsForVocabulary_video", fmt.Sprintf("error getting image url: %v", err))
			}
		}

		// audio
		if langCfg.Audio.AudioKey != "" {
			url, err := uc.s3Service.Get(ctx, langCfg.Audio.AudioKey, nil)
			if err == nil && url != nil {
				langCfg.Audio.UploadedUrl = *url
			} else {
				logger.WriteLogEx("get_vocabulary_web_usecase", "populateMediaUrlsForVocabulary_audio", fmt.Sprintf("error getting image url: %v", err))
			}
		}
	}
}

func filterIsPublished(vocabularies []model.Vocabulary) []model.Vocabulary {
	var filteredVocabularies []model.Vocabulary
	for _, vocabulary := range vocabularies {
		if vocabulary.IsPublished {
			filteredVocabularies = append(filteredVocabularies, vocabulary)
		}
	}
	return filteredVocabularies
}
