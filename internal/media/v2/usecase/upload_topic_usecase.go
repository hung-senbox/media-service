package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"media-service/helper"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
	"media-service/logger"
	"media-service/pkg/constants"
	"media-service/pkg/uploader"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadTopicUseCase interface {
	UploadTopic(ctx context.Context, req request.UploadTopicRequest) error
}

type uploadTopicUseCase struct {
	topicRepo repository.TopicRepository
	s3Service s3.Service
}

func NewUploadTopicUseCase(topicRepo repository.TopicRepository, s3Svc s3.Service) UploadTopicUseCase {
	return &uploadTopicUseCase{
		topicRepo: topicRepo,
		s3Service: s3Svc,
	}
}

// ------------------- UploadTopic main flow -------------------
func (uc *uploadTopicUseCase) UploadTopic(ctx context.Context, req request.UploadTopicRequest) error {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser == nil || !currentUser.IsSuperAdmin {
		return fmt.Errorf("access denied")
	}

	var topic *model.Topic
	var err error

	if req.TopicID != "" {
		// Case update existing topic
		topic, err = uc.updateTopicLanguage(ctx, req)
		if err != nil {
			return err
		}
	} else {
		// Case create new topic
		topic, err = uc.createTopicLanguage(ctx, req)
		if err != nil {
			return err
		}
	}

	// Thực thi upload đồng bộ, không dùng Redis
	if err := uc.uploadAndSaveAudio(ctx, topic.ID.Hex(), req); err != nil {
		logger.WriteLogMsg("error", "Failed to upload and save audio")
		logger.WriteLogEx("error", "Failed to upload and save audio", err)
		return err
	}
	if err := uc.uploadAndSaveVideo(ctx, topic.ID.Hex(), req); err != nil {
		logger.WriteLogMsg("error", "Failed to upload and save video")
		logger.WriteLogEx("error", "Failed to upload and save video", err)
		return err
	}
	if err := uc.uploadAndSaveImages(ctx, topic.ID.Hex(), req); err != nil {
		logger.WriteLogMsg("error", "Failed to upload and save images")
		logger.WriteLogEx("error", "Failed to upload and save images", err)
		return err
	}
	return nil
}

// ------------------- Upload handlers -------------------
func (uc *uploadTopicUseCase) uploadAndSaveAudio(ctx context.Context, topicID string, req request.UploadTopicRequest) error {

	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	if req.IsDeletedAudio {
		videoKey := helper.GetVideoKeyByLanguage(topic, req.LanguageID)
		if videoKey == "" {
			return fmt.Errorf("video key not found")
		}
		_ = uc.s3Service.Delete(ctx, videoKey)

		// goi repo xoa video key
		err = uc.topicRepo.DeleteVideoKey(ctx, topicID, req.LanguageID)
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveAudio] Failed to delete audio key", err)
		}
	}

	oldAudioKey := uc.getAudioKeyByLanguage(topic, req.LanguageID)
	if helper.IsValidFile(req.AudioFile) {

		key := helper.BuildObjectKeyS3("topic_media/audio", req.AudioFile.Filename, fmt.Sprintf("%s_audio", req.Title))
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
		err = uc.topicRepo.SetAudio(ctx, topicID, req.LanguageID, model.TopicAudioConfig{
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
		err := uc.topicRepo.SetAudio(ctx, topicID, req.LanguageID, model.TopicAudioConfig{
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

func (uc *uploadTopicUseCase) uploadAndSaveVideo(ctx context.Context, topicID string, req request.UploadTopicRequest) error {

	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	if req.IsDeletedVideo {
		videoKey := helper.GetVideoKeyByLanguage(topic, req.LanguageID)
		if videoKey == "" {
			return fmt.Errorf("video key not found")
		}
		_ = uc.s3Service.Delete(ctx, videoKey)

		// goi repo xoa video key
		// ignore error -> chi ra log
		err = uc.topicRepo.DeleteVideoKey(ctx, topicID, req.LanguageID)
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveVideo] Failed to delete video key", err)
		}
	}

	oldVideoKey := uc.getVideoKeyByLanguage(topic, req.LanguageID)
	if helper.IsValidFile(req.VideoFile) {

		key := helper.BuildObjectKeyS3("topic_media/video", req.VideoFile.Filename, fmt.Sprintf("%s_video", req.Title))
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
		err = uc.topicRepo.SetVideo(ctx, topicID, req.LanguageID, model.TopicVideoConfig{
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
		err := uc.topicRepo.SetVideo(ctx, topicID, req.LanguageID, model.TopicVideoConfig{
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

func (uc *uploadTopicUseCase) uploadAndSaveImages(ctx context.Context, topicID string, req request.UploadTopicRequest) error {
	imageFiles := []struct {
		file      *multipart.FileHeader
		link      string
		typ       string
		isDeleted bool
	}{
		{req.FullBackgroundFile, req.FullBackgroundLink, "full_background", req.IsDeletedFullBackground},
		{req.ClearBackgroundFile, req.ClearBackgroundLink, "clear_background", req.IsDeletedClearBackground},
		{req.ClipPartFile, req.ClipPartLink, "clip_part", req.IsDeletedClipPart},
		{req.DrawingFile, req.DrawingLink, "drawing", req.IsDeletedDrawing},
		{req.IconFile, req.IconLink, "icon", req.IsDeletedIcon},
		{req.BMFile, req.BMLink, "bm", req.IsDeletedBM},
		{req.SignLangFile, req.SignLangLink, "sign_lang", req.IsDeletedSignLang},
		{req.GifFile, req.GifLink, "gif", req.IsDeletedGif},
		{req.OrderFile, req.OrderLink, "order", req.IsDeletedOrder},
	}

	// Lấy topic 1 lần để lấy các key cũ.
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	// ================================ DELETE IMAGES ================================
	// khong handle error khi xoa image -> chi ra log
	if req.IsDeletedFullBackground {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "full_background")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete full background image", err)
		}
	}
	if req.IsDeletedClearBackground {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "clear_background")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete clear background image", err)
		}
	}
	if req.IsDeletedClipPart {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "clip_part")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete clip part image", err)
		}
	}
	if req.IsDeletedDrawing {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "drawing")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete drawing image", err)
		}
	}
	if req.IsDeletedIcon {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "icon")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete icon image", err)
		}
	}
	if req.IsDeletedBM {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "bm")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete bm image", err)
		}
	}
	if req.IsDeletedSignLang {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "sign_lang")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete sign lang image", err)
		}
	}
	if req.IsDeletedGif {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "gif")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete gif image", err)
		}
	}
	if req.IsDeletedOrder {
		err = uc.deleteImageKeyByLanguageAndType(ctx, topicID, req.LanguageID, "order")
		if err != nil {
			logger.WriteLogData("[Time: "+time.Now().Format("2006-01-02 15:04:05")+"] [uploadAndSaveImages] Failed to delete order image", err)
		}
	}
	// ================================ DELETE IMAGES ================================

	for _, img := range imageFiles {
		// Nếu có cờ xóa và không upload file mới → giữ trạng thái xóa (key rỗng), chỉ cập nhật link nếu cần
		if img.isDeleted && !helper.IsValidFile(img.file) {
			if err := uc.topicRepo.SetImage(ctx, topicID, req.LanguageID, model.TopicImageConfig{
				ImageType: img.typ,
				ImageKey:  "",
				LinkUrl:   img.link,
			}); err != nil {
				return err
			}
			continue
		}

		if helper.IsValidFile(img.file) {

			key := helper.BuildObjectKeyS3("topic_media/image", img.file.Filename, fmt.Sprintf("%s_%s_image", req.Title, img.typ))
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
			if err := uc.topicRepo.SetImage(ctx, topicID, req.LanguageID, model.TopicImageConfig{
				ImageKey:  key,
				ImageType: img.typ,
				LinkUrl:   img.link,
			}); err != nil {
				return err
			}

			continue
		} else {
			oldKey := uc.getImageKeyByLanguageAndType(topic, req.LanguageID, img.typ)
			if err := uc.topicRepo.SetImage(ctx, topicID, req.LanguageID, model.TopicImageConfig{
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

// ------------------- Topic helpers -------------------
func (uc *uploadTopicUseCase) updateTopicLanguage(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	oldTopic, err := uc.topicRepo.GetByID(ctx, req.TopicID)
	if err != nil {
		return nil, fmt.Errorf("get topic failed: %w", err)
	}

	found := false
	for i, lc := range oldTopic.LanguageConfig {
		if lc.LanguageID == req.LanguageID {
			oldTopic.LanguageConfig[i].FileName = req.FileName
			oldTopic.LanguageConfig[i].Title = req.Title
			oldTopic.LanguageConfig[i].Note = req.Note
			oldTopic.LanguageConfig[i].Description = req.Description
			found = true
			if len(lc.Images) == 0 {
				err = uc.topicRepo.InitImages(ctx, oldTopic.ID.Hex(), req.LanguageID)
				if err != nil {
					return nil, fmt.Errorf("init images fail: %w", err)
				}
			}
			if len(lc.Images) < 8 {
				err = uc.topicRepo.InitImages(ctx, oldTopic.ID.Hex(), req.LanguageID)
				if err != nil {
					return nil, fmt.Errorf("init images fail: %w", err)
				}
			}
			break
		}
	}

	if !found {
		oldTopic.LanguageConfig = append(oldTopic.LanguageConfig, model.TopicLanguageConfig{
			LanguageID:  req.LanguageID,
			FileName:    req.FileName,
			Title:       req.Title,
			Note:        req.Note,
			Description: req.Description,
		})
		// init images
		err = uc.topicRepo.InitImages(ctx, oldTopic.ID.Hex(), req.LanguageID)
		if err != nil {
			return nil, fmt.Errorf("init images fail: %w", err)
		}
	}

	oldTopic.IsPublished = req.IsPublished
	return uc.topicRepo.UpdateTopic(ctx, oldTopic)
}

func (uc *uploadTopicUseCase) createTopicLanguage(ctx context.Context, req request.UploadTopicRequest) (*model.Topic, error) {
	topic := &model.Topic{
		ID:             primitive.NewObjectID(),
		IsPublished:    req.IsPublished,
		LanguageConfig: []model.TopicLanguageConfig{},
	}

	newTopic, err := uc.topicRepo.CreateTopic(ctx, topic)
	if err != nil {
		return nil, fmt.Errorf("create topic fail: %w", err)
	}

	langConfig := model.TopicLanguageConfig{
		LanguageID:  req.LanguageID,
		FileName:    req.FileName,
		Title:       req.Title,
		Note:        req.Note,
		Description: req.Description,
		Images:      []model.TopicImageConfig{},
		Audio:       model.TopicAudioConfig{},
		Video:       model.TopicVideoConfig{},
	}
	if err := uc.topicRepo.SetLanguageConfig(ctx, newTopic.ID.Hex(), langConfig); err != nil {
		return nil, fmt.Errorf("set language config fail: %w", err)
	}

	// init images
	err = uc.topicRepo.InitImages(ctx, newTopic.ID.Hex(), req.LanguageID)
	if err != nil {
		return nil, fmt.Errorf("init images fail: %w", err)
	}

	return newTopic, nil
}

func (uc *uploadTopicUseCase) getAudioKeyByLanguage(topic *model.Topic, languageID uint) string {
	for _, lc := range topic.LanguageConfig {
		if lc.LanguageID == languageID {
			if lc.Audio.AudioKey != "" {
				return lc.Audio.AudioKey
			}
			break
		}
	}
	return ""
}

func (uc *uploadTopicUseCase) getVideoKeyByLanguage(topic *model.Topic, languageID uint) string {
	for _, lc := range topic.LanguageConfig {
		if lc.LanguageID == languageID {
			if lc.Video.VideoKey != "" {
				return lc.Video.VideoKey
			}
			break
		}
	}
	return ""
}

func (uc *uploadTopicUseCase) getImageKeyByLanguageAndType(topic *model.Topic, languageID uint, imageType string) string {
	for _, lc := range topic.LanguageConfig {
		if lc.LanguageID == languageID {
			for _, img := range lc.Images {
				if img.ImageType == imageType {
					return img.ImageKey
				}
			}
			break
		}
	}
	return ""
}

func (uc *uploadTopicUseCase) deleteImageKeyByLanguageAndType(ctx context.Context, topicID string, languageID uint, imageType string) error {
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}
	oldKey := uc.getImageKeyByLanguageAndType(topic, languageID, imageType)
	if oldKey != "" {
		err = uc.s3Service.Delete(ctx, oldKey)
		if err != nil {
			logger.WriteLogEx("error", "Failed to delete s3 service image", err)
		}
	}
	// goi repo xoa image key
	err = uc.topicRepo.DeleteImageKey(ctx, topicID, languageID, imageType)
	if err != nil {
		return err
	}
	return nil

}
