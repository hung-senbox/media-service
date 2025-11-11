package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
)

func ToGetVideosResponse4Web(videoUploaders []model.VideoUploader, createdByName string) []response.GetVideoUploaderResponse4Web {
	var result = make([]response.GetVideoUploaderResponse4Web, 0)
	for _, videoUploader := range videoUploaders {
		// optional filter by language
		creatorName := createdByName
		if creatorName == "" {
			creatorName = videoUploader.CreatedBy
		}
		result = append(result, response.GetVideoUploaderResponse4Web{
			ID:              videoUploader.ID.Hex(),
			LanguageID:      videoUploader.LanguageID,
			IsVisible:       videoUploader.IsVisible,
			CreatedByName:   creatorName,
			Title:           videoUploader.Title,
			VideoUrl:        videoUploader.VideoPublicUrl,
			ImagePreviewUrl: videoUploader.ImagePreviewPublicUrl,
		})
	}
	return result
}
