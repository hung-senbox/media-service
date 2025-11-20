package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"media-service/helper"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
	"media-service/logger"
	"media-service/pkg/constants"
	"media-service/pkg/uploader"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadVocabularyUseCase interface {
	UploadVocabulary(ctx context.Context, req request.UploadVocabularyRequest) error
}

type uploadVocabularyUseCase struct {
	topicRepo      repository.TopicRepository
	vocabularyRepo repository.VocabularyRepository
	s3Service      s3.Service
}

func NewUploadVocabularyUseCase(topicRepo repository.TopicRepository, vocabularyRepo repository.VocabularyRepository, s3Svc s3.Service) UploadVocabularyUseCase {
	return &uploadVocabularyUseCase{
		topicRepo:      topicRepo,
		vocabularyRepo: vocabularyRepo,
		s3Service:      s3Svc,
	}
}

// ------------------- UploadVocabulary main flow -------------------
func (uc *uploadVocabularyUseCase) UploadVocabulary(ctx context.Context, req request.UploadVocabularyRequest) error {

	var vocabulary *model.Vocabulary
	var err error

	if req.VocabularyID != "" {
		// Case update existing vocabulary
		if req.TopicID == "" {
			return fmt.Errorf("topic id is required")
		}
		topic, err := uc.topicRepo.GetByID(ctx, req.TopicID)
		if err != nil {
			return fmt.Errorf("get topic failed: %w", err)
		}
		if topic == nil {
			return fmt.Errorf("topic not found")
		}
		vocabulary, err = uc.updateVocabulary(ctx, req)
		if err != nil {
			return err
		}
	} else {
		// Case create new vocabulary
		vocabulary, err = uc.createVocabulary(ctx, req)
		if err != nil {
			return err
		}
	}

	// Thực thi upload đồng bộ, không dùng Redis
	if err := uc.uploadAndSaveAudio(ctx, vocabulary, req); err != nil {
		logger.WriteLogMsg("error", "Failed to upload and save audio")
		logger.WriteLogEx("error", "Failed to upload and save audio", err)
		return err
	}
	if err := uc.uploadAndSaveVideo(ctx, vocabulary, req); err != nil {
		logger.WriteLogMsg("error", "Failed to upload and save video")
		logger.WriteLogEx("error", "Failed to upload and save video", err)
		return err
	}
	if err := uc.uploadAndSaveImages(ctx, vocabulary, req); err != nil {
		logger.WriteLogMsg("error", "Failed to upload and save images")
		logger.WriteLogEx("error", "Failed to upload and save images", err)
		return err
	}
	return nil
}

// ------------------- Upload handlers -------------------
func (uc *uploadVocabularyUseCase) uploadAndSaveAudio(ctx context.Context, vocabulary *model.Vocabulary, req request.UploadVocabularyRequest) error {
	vocabularyID := vocabulary.ID.Hex()

	if req.IsDeletedAudio {
		audioKey := helper.GetVocabularyAudioKeyByLanguage(vocabulary, req.LanguageID)
		if audioKey != "" {
			_ = uc.s3Service.Delete(ctx, audioKey)
		}
		// goi repo xoa audio key
		if err := uc.vocabularyRepo.DeleteAudioKey(ctx, vocabularyID, req.LanguageID); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveAudio] Failed to delete audio key", err)
		}
		// Nếu không có file mới, giữ trạng thái xóa (key rỗng) và chỉ cập nhật metadata
		if !helper.IsValidFile(req.AudioFile) {
			return uc.vocabularyRepo.SetAudio(ctx, vocabularyID, req.LanguageID, model.VocabularyAudioConfig{
				AudioKey:  "",
				LinkUrl:   req.AudioLinkUrl,
				StartTime: req.AudioStart,
				EndTime:   req.AudioEnd,
			})
		}
	}

	oldAudioKey := helper.GetVocabularyAudioKeyByLanguage(vocabulary, req.LanguageID)
	if helper.IsValidFile(req.AudioFile) {

		key := helper.BuildObjectKeyS3("vocabulary_media/audio", req.AudioFile.Filename, fmt.Sprintf("%s_audio", req.Title))
		f, openErr := req.AudioFile.Open()
		if openErr != nil {
			return openErr
		}
		defer f.Close()
		ct := req.AudioFile.Header.Get("Content-Type")
		_, err := uc.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPrivate)
		if err != nil {
			return err
		}
		// cập nhật metadata + key (mới hoặc cũ)
		err = uc.vocabularyRepo.SetAudio(ctx, vocabularyID, req.LanguageID, model.VocabularyAudioConfig{
			AudioKey:  key,
			LinkUrl:   req.AudioLinkUrl,
			StartTime: req.AudioStart,
			EndTime:   req.AudioEnd,
		})
		if err != nil {
			return err
		}
	} else {
		// cập nhật metadata + key (mới hoặc cũ)
		err := uc.vocabularyRepo.SetAudio(ctx, vocabularyID, req.LanguageID, model.VocabularyAudioConfig{
			AudioKey:  oldAudioKey,
			LinkUrl:   req.AudioLinkUrl,
			StartTime: req.AudioStart,
			EndTime:   req.AudioEnd,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *uploadVocabularyUseCase) uploadAndSaveVideo(ctx context.Context, vocabulary *model.Vocabulary, req request.UploadVocabularyRequest) error {
	vocabularyID := vocabulary.ID.Hex()

	if req.IsDeletedVideo {
		videoKey := helper.GetVocabularyVideoKeyByLanguage(vocabulary, req.LanguageID)
		if videoKey == "" {
			return fmt.Errorf("video key not found")
		}
		_ = uc.s3Service.Delete(ctx, videoKey)

		// goi repo xoa video key (ignore error -> chi ra log)
		if err := uc.vocabularyRepo.DeleteVideoKey(ctx, vocabularyID, req.LanguageID); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveVideo] Failed to delete video key", err)
		}
		// Nếu không có file mới, giữ trạng thái xóa (key rỗng) và chỉ cập nhật metadata
		if !helper.IsValidFile(req.VideoFile) {
			return uc.vocabularyRepo.SetVideo(ctx, vocabularyID, req.LanguageID, model.VocabularyVideoConfig{
				VideoKey:  "",
				LinkUrl:   req.VideoLinkUrl,
				StartTime: req.VideoStart,
				EndTime:   req.VideoEnd,
			})
		}
	}

	oldVideoKey := helper.GetVocabularyVideoKeyByLanguage(vocabulary, req.LanguageID)
	if helper.IsValidFile(req.VideoFile) {

		key := helper.BuildObjectKeyS3("vocabulary_media/video", req.VideoFile.Filename, fmt.Sprintf("%s_video", req.Title))
		f, openErr := req.VideoFile.Open()
		if openErr != nil {
			return openErr
		}
		defer f.Close()
		ct := req.VideoFile.Header.Get("Content-Type")
		_, err := uc.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPrivate)
		if err != nil {
			return err
		}
		err = uc.vocabularyRepo.SetVideo(ctx, vocabularyID, req.LanguageID, model.VocabularyVideoConfig{
			VideoKey:  key,
			LinkUrl:   req.VideoLinkUrl,
			StartTime: req.VideoStart,
			EndTime:   req.VideoEnd,
		})
		if err != nil {
			return err
		}
	} else {
		// cập nhật metadata + key (mới hoặc cũ)
		err := uc.vocabularyRepo.SetVideo(ctx, vocabularyID, req.LanguageID, model.VocabularyVideoConfig{
			VideoKey:  oldVideoKey,
			LinkUrl:   req.VideoLinkUrl,
			StartTime: req.VideoStart,
			EndTime:   req.VideoEnd,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *uploadVocabularyUseCase) uploadAndSaveImages(ctx context.Context, vocabulary *model.Vocabulary, req request.UploadVocabularyRequest) error {
	vocabularyID := vocabulary.ID.Hex()
	imageFiles := []struct {
		file      *multipart.FileHeader
		link      string
		typ       string
		isDeleted bool
	}{
		{req.FullBackgroundFile, req.FullBackgroundLink, string(constants.TopicImageTypeFullBackground), req.IsDeletedFullBackground},
		{req.ClearBackgroundFile, req.ClearBackgroundLink, string(constants.TopicImageTypeClearBackground), req.IsDeletedClearBackground},
		{req.ClipPartFile, req.ClipPartLink, string(constants.TopicImageTypeClipPart), req.IsDeletedClipPart},
		{req.DrawingFile, req.DrawingLink, string(constants.TopicImageTypeDrawing), req.IsDeletedDrawing},
		{req.IconFile, req.IconLink, string(constants.TopicImageTypeIcon), req.IsDeletedIcon},
		{req.BMFile, req.BMLink, string(constants.TopicImageTypeBM), req.IsDeletedBM},
		{req.SignLangFile, req.SignLangLink, string(constants.TopicImageTypeSignLang), req.IsDeletedSignLang},
		{req.GifFile, req.GifLink, string(constants.TopicImageTypeGif), req.IsDeletedGif},
		{req.OrderFile, req.OrderLink, string(constants.TopicImageTypeOrder), req.IsDeletedOrder},
	}

	// vocabulary is already available

	// ================================ DELETE IMAGES ================================
	// khong handle error khi xoa image -> chi ra log
	if req.IsDeletedFullBackground {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeFullBackground)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete full background image", err)
		}
	}
	if req.IsDeletedClearBackground {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeClearBackground)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete clear background image", err)
		}
	}
	if req.IsDeletedClipPart {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeClipPart)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete clip part image", err)
		}
	}
	if req.IsDeletedDrawing {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeDrawing)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete drawing image", err)
		}
	}
	if req.IsDeletedIcon {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeIcon)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete icon image", err)
		}
	}
	if req.IsDeletedBM {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeBM)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete bm image", err)
		}
	}
	if req.IsDeletedSignLang {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeSignLang)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete sign lang image", err)
		}
	}
	if req.IsDeletedGif {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeGif)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete gif image", err)
		}
	}
	if req.IsDeletedOrder {
		if err := uc.deleteImageKeyByLanguageAndType(ctx, vocabularyID, req.LanguageID, string(constants.TopicImageTypeOrder)); err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete order image", err)
		}
	}
	// ================================ DELETE IMAGES ================================

	for _, img := range imageFiles {
		// Nếu có cờ xóa và không upload file mới → giữ trạng thái xóa (key rỗng), chỉ cập nhật link nếu cần
		if img.isDeleted && !helper.IsValidFile(img.file) {
			if err := uc.vocabularyRepo.SetImage(ctx, vocabularyID, req.LanguageID, model.VocabularyImageConfig{
				ImageType: img.typ,
				ImageKey:  "",
				LinkUrl:   img.link,
			}); err != nil {
				return err
			}
			continue
		}

		if helper.IsValidFile(img.file) {

			key := helper.BuildObjectKeyS3("vocabulary_media/image", img.file.Filename, fmt.Sprintf("%s_%s_image", req.Title, img.typ))
			f, openErr := img.file.Open()
			if openErr != nil {
				return openErr
			}
			ct := img.file.Header.Get("Content-Type")
			if _, err := uc.s3Service.SaveReader(ctx, f, key, ct, uploader.UploadPrivate); err != nil {
				_ = f.Close()
				return err
			}
			_ = f.Close()

			// Lưu key + metadata mới
			if err := uc.vocabularyRepo.SetImage(ctx, vocabularyID, req.LanguageID, model.VocabularyImageConfig{
				ImageKey:  key,
				ImageType: img.typ,
				LinkUrl:   img.link,
			}); err != nil {
				return err
			}

			continue
		} else {
			oldKey := helper.GetVocabularyImageKeyByLanguageAndType(vocabulary, req.LanguageID, img.typ)
			if err := uc.vocabularyRepo.SetImage(ctx, vocabularyID, req.LanguageID, model.VocabularyImageConfig{
				ImageKey:  oldKey,
				ImageType: img.typ,
				LinkUrl:   img.link,
			}); err != nil {
				// chỉ log warning, không ghi Redis error
				logger.WriteLogData("[uploadAndSaveImages] Failed to update metadata case2", err)
			}
		}
	}

	return nil
}

// ------------------- Vocabulary helpers -------------------
func (uc *uploadVocabularyUseCase) updateVocabulary(ctx context.Context, req request.UploadVocabularyRequest) (*model.Vocabulary, error) {
	oldVocabulary, err := uc.vocabularyRepo.GetByID(ctx, req.VocabularyID)
	if err != nil {
		return nil, fmt.Errorf("get vocabulary failed: %w", err)
	}

	found := false
	for i, lc := range oldVocabulary.LanguageConfig {
		if lc.LanguageID == req.LanguageID {
			oldVocabulary.LanguageConfig[i].FileName = req.FileName
			oldVocabulary.LanguageConfig[i].Title = req.Title
			oldVocabulary.LanguageConfig[i].Note = req.Note
			oldVocabulary.LanguageConfig[i].Description = req.Description
			found = true
			if len(lc.Images) == 0 {
				err = uc.vocabularyRepo.InitImages(ctx, oldVocabulary.ID.Hex(), req.LanguageID)
				if err != nil {
					return nil, fmt.Errorf("init images fail: %w", err)
				}
			}
			break
		}
	}

	if !found {
		oldVocabulary.LanguageConfig = append(oldVocabulary.LanguageConfig, model.VocabularyLanguageConfig{
			LanguageID:  req.LanguageID,
			FileName:    req.FileName,
			Title:       req.Title,
			Note:        req.Note,
			Description: req.Description,
		})
		// init images
		err = uc.vocabularyRepo.InitImages(ctx, oldVocabulary.ID.Hex(), req.LanguageID)
		if err != nil {
			return nil, fmt.Errorf("init images fail: %w", err)
		}
	}

	oldVocabulary.IsPublished = req.IsPublished
	return uc.vocabularyRepo.UpdateVocabulary(ctx, oldVocabulary)
}

func (uc *uploadVocabularyUseCase) createVocabulary(ctx context.Context, req request.UploadVocabularyRequest) (*model.Vocabulary, error) {
	vocabulary := &model.Vocabulary{
		ID:             primitive.NewObjectID(),
		TopicID:        req.TopicID,
		IsPublished:    req.IsPublished,
		LanguageConfig: []model.VocabularyLanguageConfig{},
	}

	newVocabulary, err := uc.vocabularyRepo.CreateVocabulary(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("create topic fail: %w", err)
	}

	langConfig := model.VocabularyLanguageConfig{
		LanguageID:  req.LanguageID,
		FileName:    req.FileName,
		Title:       req.Title,
		Note:        req.Note,
		Description: req.Description,
		Images:      []model.VocabularyImageConfig{},
		Audio:       model.VocabularyAudioConfig{},
		Video:       model.VocabularyVideoConfig{},
	}
	if err := uc.vocabularyRepo.SetLanguageConfig(ctx, vocabulary.ID.Hex(), langConfig); err != nil {
		return nil, fmt.Errorf("set language config fail: %w", err)
	}

	// init images
	if err := uc.vocabularyRepo.InitImages(ctx, vocabulary.ID.Hex(), req.LanguageID); err != nil {
		return nil, fmt.Errorf("init vocabulary images fail: %w", err)
	}

	return newVocabulary, nil
}

func (uc *uploadVocabularyUseCase) deleteImageKeyByLanguageAndType(ctx context.Context, vocabularyID string, languageID uint, imageType string) error {
	vocabulary, err := uc.vocabularyRepo.GetByID(ctx, vocabularyID)
	if err != nil {
		return err
	}
	oldKey := helper.GetVocabularyImageKeyByLanguageAndType(vocabulary, languageID, imageType)
	if oldKey != "" {
		err = uc.s3Service.Delete(ctx, oldKey)
		if err != nil {
			logger.WriteLogEx("error", "Failed to delete s3 service image", err)
		}
	}
	// goi repo xoa image key
	err = uc.vocabularyRepo.DeleteImageKey(ctx, vocabularyID, languageID, imageType)
	if err != nil {
		return err
	}
	return nil

}
