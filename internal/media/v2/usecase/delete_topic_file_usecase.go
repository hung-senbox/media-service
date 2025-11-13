package usecase

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
)

type DeleteTopicFileUseCase interface {
	DeleteTopicAudioKey(ctx context.Context, topicID string, languageID uint) error
	DeleteTopicVideoKey(ctx context.Context, topicID string, languageID uint) error
	DeleteTopicImageKey(ctx context.Context, topicID string, languageID uint, imageType string) error
}

type deleteTopicFileUseCase struct {
	topicRepo repository.TopicRepository
	s3Service s3.Service
}

func NewDeleteTopicFileUseCase(topicRepo repository.TopicRepository, s3Service s3.Service) DeleteTopicFileUseCase {
	return &deleteTopicFileUseCase{topicRepo: topicRepo, s3Service: s3Service}
}

func (uc *deleteTopicFileUseCase) DeleteTopicAudioKey(ctx context.Context, topicID string, languageID uint) error {
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	if topic == nil {
		return fmt.Errorf("topic not found")
	}

	audioKey := helper.GetAudioKeyByLanguage(topic, languageID)
	if audioKey == "" {
		return fmt.Errorf("audio key not found")
	}

	err = uc.s3Service.Delete(ctx, audioKey)
	if err != nil {
		return err
	}

	// goi repo xoa audio key
	err = uc.topicRepo.DeleteAudioKey(ctx, topicID, languageID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *deleteTopicFileUseCase) DeleteTopicVideoKey(ctx context.Context, topicID string, languageID uint) error {
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	if topic == nil {
		return fmt.Errorf("topic not found")
	}

	videoKey := helper.GetVideoKeyByLanguage(topic, languageID)
	if videoKey == "" {
		return fmt.Errorf("video key not found")
	}

	err = uc.s3Service.Delete(ctx, videoKey)
	if err != nil {
		return err
	}

	// goi repo xoa video key
	err = uc.topicRepo.DeleteVideoKey(ctx, topicID, languageID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *deleteTopicFileUseCase) DeleteTopicImageKey(ctx context.Context, topicID string, languageID uint, imageType string) error {
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	if topic == nil {
		return fmt.Errorf("topic not found")
	}

	imageKey := helper.GetImageKeyByLanguageAndType(topic, languageID, imageType)
	if imageKey == "" {
		return fmt.Errorf("image key not found")
	}

	err = uc.s3Service.Delete(ctx, imageKey)
	if err != nil {
		return err
	}

	// goi repo xoa image key
	err = uc.topicRepo.DeleteImageKey(ctx, topicID, languageID, imageType)
	if err != nil {
		return err
	}

	return nil
}
