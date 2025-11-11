package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
)

func ToGetVideosResponse4Web(videoUploaders []model.VideoUploader, createdByName string, languageID uint) []response.GetVideoUploaderResponse4Web {
	var result = make([]response.GetVideoUploaderResponse4Web, 0)
	for _, videoUploader := range videoUploaders {
		creatorName := createdByName
		if creatorName == "" {
			creatorName = videoUploader.CreatedBy
		}
		title := ""
		videoURL := ""
		imageURL := ""
		if len(videoUploader.LanguageConfig) > 0 {
			found := false
			for _, lc := range videoUploader.LanguageConfig {
				if lc.LanguageID == languageID {
					title = lc.Title
					videoURL = lc.VideoPublicUrl
					imageURL = lc.ImagePreviewPublicUrl
					found = true
					break
				}
			}
			if !found {
				lc := videoUploader.LanguageConfig[0]
				title = lc.Title
				videoURL = lc.VideoPublicUrl
				imageURL = lc.ImagePreviewPublicUrl
			}
		}
		result = append(result, response.GetVideoUploaderResponse4Web{
			ID:              videoUploader.ID.Hex(),
			IsVisible:       videoUploader.IsVisible,
			CreatedByName:   creatorName,
			Title:           title,
			VideoUrl:        videoURL,
			ImagePreviewUrl: imageURL,
		})
	}
	return result
}
