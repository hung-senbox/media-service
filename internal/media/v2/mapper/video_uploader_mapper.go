package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
)

func ToGetVideosResponse4Web(videoUploaders []model.VideoUploader, createdByName string) []response.GetVideoUploaderResponse4Web {
	var result = make([]response.GetVideoUploaderResponse4Web, 0)
	for _, videoUploader := range videoUploaders {
		result = append(result, response.GetVideoUploaderResponse4Web{
			ID:              videoUploader.ID.Hex(),
			Title:           videoUploader.Title,
			IsVisible:       videoUploader.IsVisible,
			CreatedByName:   createdByName,
			CreatedAt:       videoUploader.CreatedAt,
			UpdatedAt:       videoUploader.UpdatedAt,
			VideoURL:        videoUploader.VideoPublicUrl,
			ImagePreviewURL: videoUploader.ImagePreviewPublicUrl,
		})
	}
	return result
}
