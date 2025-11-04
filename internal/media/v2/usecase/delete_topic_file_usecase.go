package usecase

import (
	"context"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway"
	"media-service/internal/media/v2/repository"
)

type DeleteTopicFileUseCase interface {
	DeleteTopicAudioKey(ctx context.Context, topicID string, languageID uint) error
	DeleteTopicVideoKey(ctx context.Context, topicID string, languageID uint) error
	DeleteTopicImageKey(ctx context.Context, topicID string, languageID uint, imageType string) error
}

type deleteTopicFileUseCase struct {
	topicRepo   repository.TopicRepository
	fileGateway gateway.FileGateway
}

func NewDeleteTopicFileUseCase(topicRepo repository.TopicRepository, fileGateway gateway.FileGateway) DeleteTopicFileUseCase {
	return &deleteTopicFileUseCase{topicRepo: topicRepo, fileGateway: fileGateway}
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

	err = uc.fileGateway.DeleteAudio(ctx, audioKey)
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

	err = uc.fileGateway.DeleteVideo(ctx, videoKey)
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

	err = uc.fileGateway.DeleteImage(ctx, imageKey)
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
