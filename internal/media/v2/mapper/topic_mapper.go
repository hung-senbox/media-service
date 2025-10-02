package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"strings"
)

func ToTopicResponses(topics []model.Topic) []response.TopicResponse {
	var result = make([]response.TopicResponse, 0)

	for _, t := range topics {
		resp := response.TopicResponse{
			ID:          t.ID.Hex(),
			IsPublished: t.IsPublished,
		}

		var langs []response.MessageLanguageEntry
		for _, lc := range t.LanguageConfig {

			entry := response.MessageLanguageEntry{
				LanguageID: int(lc.LanguageID),
				Contents: response.LanguageContents{
					FileName:    lc.FileName,
					Title:       lc.Title,
					Note:        lc.Note,
					Description: lc.Description,
				},
			}

			// map audio
			if lc.Audio.LinkUrl != "" {
				entry.Contents.Audio = &response.MediaContent{
					UploadedURL: &lc.Audio.AudioKey,
					LinkURL:     lc.Audio.LinkUrl,
					StartTime:   strPtr(trimQuotes(lc.Audio.StartTime)),
					EndTime:     strPtr(trimQuotes(lc.Audio.EndTime)),
				}
			}

			// map video
			if lc.Video.LinkUrl != "" {
				entry.Contents.Video = &response.MediaContent{
					UploadedURL: &lc.Video.VideoKey,
					LinkURL:     lc.Video.LinkUrl,
					StartTime:   strPtr(trimQuotes(lc.Video.StartTime)),
					EndTime:     strPtr(trimQuotes(lc.Video.EndTime)),
				}
			}

			// map images slice → object
			imgMap := make(map[string]response.ImgEntry)
			for _, img := range lc.Images {
				uploaded := img.ImageKey
				imgMap[img.ImageType] = response.ImgEntry{
					UploadedURL: &uploaded,
					LinkURL:     img.LinkUrl,
				}
			}
			entry.Contents.Images = imgMap

			langs = append(langs, entry)
		}

		resp.MessageLangs = langs
		result = append(result, resp)
	}

	return result
}

func trimQuotes(s string) string {
	return strings.Trim(s, "\"")
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
