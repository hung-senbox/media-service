package usecase

import (
	"context"
	"fmt"
	"mime/multipart"

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
		return err
	}
	if err := uc.uploadAndSaveVideo(ctx, topic.ID.Hex(), req); err != nil {
		return err
	}
	if err := uc.uploadAndSaveImages(ctx, topic.ID.Hex(), req); err != nil {
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

	oldAudioKey := uc.getAudioKeyByLanguage(topic, req.LanguageID)
	if helper.IsValidFile(req.AudioFile) {
		// xóa file cũ nếu có
		if oldAudioKey != "" {
			_ = uc.s3Service.Delete(ctx, oldAudioKey)
		}

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

	oldVideoKey := uc.getVideoKeyByLanguage(topic, req.LanguageID)
	if helper.IsValidFile(req.VideoFile) {
		// xóa file cũ nếu có
		if oldVideoKey != "" {
			_ = uc.s3Service.Delete(ctx, oldVideoKey)
		}

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
		file *multipart.FileHeader
		link string
		typ  string
	}{
		{req.FullBackgroundFile, req.FullBackgroundLink, "full_background"},
		{req.ClearBackgroundFile, req.ClearBackgroundLink, "clear_background"},
		{req.ClipPartFile, req.ClipPartLink, "clip_part"},
		{req.DrawingFile, req.DrawingLink, "drawing"},
		{req.IconFile, req.IconLink, "icon"},
		{req.BMFile, req.BMLink, "bm"},
		{req.SignLangFile, req.SignLangLink, "sign_lang"},
		{req.GifFile, req.GifLink, "gif"},
		{req.OrderFile, req.OrderLink, "order"},
	}

	// Lấy topic 1 lần để lấy các key cũ.
	topic, err := uc.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	for _, img := range imageFiles {
		// 1) Nếu có file => đây là task upload (done() phải được gọi)
		if helper.IsValidFile(img.file) {
			oldKey := uc.getImageKeyByLanguageAndType(topic, req.LanguageID, img.typ)
			if oldKey != "" {
				// cố gắng xóa file cũ (ignore error)
				_ = uc.s3Service.Delete(ctx, oldKey)
			}

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
