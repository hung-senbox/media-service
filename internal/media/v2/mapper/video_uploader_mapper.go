package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
)

func ToGetVideosResponse4Web(videoUploaders []model.VideoUploader, createdByName string) []response.GetVideoUploaderResponse4Web {
	var result = make([]response.GetVideoUploaderResponse4Web, 0)
	for _, videoUploader := range videoUploaders {
		creatorName := createdByName
		if creatorName == "" {
			creatorName = videoUploader.CreatedBy
		}
		langItems := make([]response.MessageLanguageVideo, 0, len(videoUploader.LanguageConfig))
		for _, lc := range videoUploader.LanguageConfig {
			langItems = append(langItems, response.MessageLanguageVideo{
				LanguageID: int(lc.LanguageID),
				Contents: response.VideoLanguageContents{
					Title:           lc.Title,
					VideoURL:        lc.VideoPublicUrl,
					ImagePreviewURL: lc.ImagePreviewPublicUrl,
				},
			})
		}
		result = append(result, response.GetVideoUploaderResponse4Web{
			ID:            videoUploader.ID.Hex(),
			IsVisible:     videoUploader.IsVisible,
			CreatedByName: creatorName,
			MessageLangs:  langItems,
		})
	}
	return result
}
