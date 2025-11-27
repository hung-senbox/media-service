package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
)

// langID == 0 => lấy config đầu tiên; ngược lại ưu tiên đúng languageID nếu tồn tại
func ToGetVideosResponse4Web(videoUploaders []model.VideoUploader, createdByName string, langID uint) []response.GetVideoUploaderResponse4Web {
	var result = make([]response.GetVideoUploaderResponse4Web, 0)
	for _, videoUploader := range videoUploaders {
		if len(videoUploader.LanguageConfig) == 0 {
			continue
		}

		// chọn language config phù hợp
		cfg := &videoUploader.LanguageConfig[0]
		if langID != 0 {
			for i := range videoUploader.LanguageConfig {
				if videoUploader.LanguageConfig[i].LanguageID == langID {
					cfg = &videoUploader.LanguageConfig[i]
					break
				}
			}
		}

		creatorName := createdByName
		if creatorName == "" {
			creatorName = videoUploader.CreatedBy
		}

		result = append(result, response.GetVideoUploaderResponse4Web{
			ID:               videoUploader.ID.Hex(),
			LanguageID:       cfg.LanguageID,
			LanguageConfigID: cfg.ID.Hex(),
			IsVisible:        videoUploader.IsVisible,
			CreatedByName:    creatorName,
			Title:            videoUploader.Title,
			WikiCode:         videoUploader.WikiCode,
			VideoUrl:         cfg.VideoPublicUrl,
			ImagePreviewUrl:  cfg.ImagePreviewPublicUrl,
			Note:             cfg.Note,
			Transcript:       cfg.Transcript,
			CreatedAt:        videoUploader.CreatedAt,
		})
	}
	return result
}

func ToGetDetailVideo4WebResponse(videoUploader *model.VideoUploader) *response.GetDetailVideo4WebResponse {
	var result = make([]response.DetailVideoMessageLanguageEntry, 0)
	for _, cfg := range videoUploader.LanguageConfig {
		result = append(result, response.DetailVideoMessageLanguageEntry{
			LanguageID: int(cfg.LanguageID),
			Contents: response.DetailVideoLanguageContents{
				Note:            cfg.Note,
				Transcript:      cfg.Transcript,
				VideoUrl:        cfg.VideoPublicUrl,
				ImagePreviewUrl: cfg.ImagePreviewPublicUrl,
			},
		})
	}
	return &response.GetDetailVideo4WebResponse{
		ID:            videoUploader.ID.Hex(),
		IsVisible:     videoUploader.IsVisible,
		Title:         videoUploader.Title,
		WikiCode:      videoUploader.WikiCode,
		CreatedByName: videoUploader.CreatedBy,
		MessageLangs:  result,
		CreatedAt:     videoUploader.CreatedAt,
	}
}

func ToVideo4GwResponse(videoUploader *model.VideoUploader, languageID uint) *response.GetVideo4GwResponse {
	videoUrl := ""
	imagePreviewUrl := ""
	if languageID != 0 {
		for _, cfg := range videoUploader.LanguageConfig {
			if cfg.LanguageID == languageID {
				videoUrl = cfg.VideoPublicUrl
				imagePreviewUrl = cfg.ImagePreviewPublicUrl
			}
		}
	}
	return &response.GetVideo4GwResponse{
		ID:              videoUploader.ID.Hex(),
		Title:           videoUploader.Title,
		WikiCode:        videoUploader.WikiCode,
		VideoUrl:        videoUrl,
		ImagePreviewUrl: imagePreviewUrl,
		CreatedAt:       videoUploader.CreatedAt,
	}
}
