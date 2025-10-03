package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"strings"
)

func ToTopicResponses4Web(topics []model.Topic) []response.TopicResponse4Web {
	var result = make([]response.TopicResponse4Web, 0)

	for _, t := range topics {
		resp := response.TopicResponse4Web{
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
					UploadedURL: &lc.Audio.UploadedUrl,
					LinkURL:     lc.Audio.LinkUrl,
					StartTime:   strPtr(trimQuotes(lc.Audio.StartTime)),
					EndTime:     strPtr(trimQuotes(lc.Audio.EndTime)),
				}
			}

			// map video
			if lc.Video.LinkUrl != "" {
				entry.Contents.Video = &response.MediaContent{
					UploadedURL: &lc.Video.UploadedUrl,
					LinkURL:     lc.Video.LinkUrl,
					StartTime:   strPtr(trimQuotes(lc.Video.StartTime)),
					EndTime:     strPtr(trimQuotes(lc.Video.EndTime)),
				}
			}

			// map images slice → object
			imgMap := make(map[string]response.ImgEntry)
			for _, img := range lc.Images {
				uploaded := img.UploadedUrl
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

func ToTopicResponse4Web(t *model.Topic) *response.TopicResponse4Web {
	if t == nil {
		return nil
	}

	resp := &response.TopicResponse4Web{
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
				UploadedURL: &lc.Audio.UploadedUrl,
				LinkURL:     lc.Audio.LinkUrl,
				StartTime:   strPtr(trimQuotes(lc.Audio.StartTime)),
				EndTime:     strPtr(trimQuotes(lc.Audio.EndTime)),
			}
		}

		// map video
		if lc.Video.LinkUrl != "" {
			entry.Contents.Video = &response.MediaContent{
				UploadedURL: &lc.Video.UploadedUrl,
				LinkURL:     lc.Video.LinkUrl,
				StartTime:   strPtr(trimQuotes(lc.Video.StartTime)),
				EndTime:     strPtr(trimQuotes(lc.Video.EndTime)),
			}
		}

		// map images slice → object
		imgMap := make(map[string]response.ImgEntry)
		for _, img := range lc.Images {
			uploaded := img.UploadedUrl
			imgMap[img.ImageType] = response.ImgEntry{
				UploadedURL: &uploaded,
				LinkURL:     img.LinkUrl,
			}
		}
		entry.Contents.Images = imgMap

		langs = append(langs, entry)
	}

	resp.MessageLangs = langs
	return resp
}

func ToTopic4StudentResponses4App(topics []model.Topic, appLanguage uint) []*response.GetTopic4StudentResponse4App {
	var res = make([]*response.GetTopic4StudentResponse4App, 0)

	for _, t := range topics {
		// is published = false
		if !t.IsPublished {
			continue
		}
		// chọn language config
		var langConfig *model.TopicLanguageConfig
		for _, lc := range t.LanguageConfig {
			if lc.LanguageID == appLanguage {
				langConfig = &lc
				break
			}
		}

		if langConfig == nil {
			continue
		}

		res = append(res, &response.GetTopic4StudentResponse4App{
			ID:          t.ID.Hex(),
			IsPublished: t.IsPublished,
			Title:       langConfig.Title,
		})
	}

	return res
}

func ToTopic4StudentResponses4Web(topics []model.Topic, appLanguage uint) []*response.GetTopic4StudentResponse4Web {
	var res = make([]*response.GetTopic4StudentResponse4Web, 0)

	for _, t := range topics {
		if !t.IsPublished {
			continue
		}
		// chọn language config
		var langConfig *model.TopicLanguageConfig
		for _, lc := range t.LanguageConfig {
			if lc.LanguageID == appLanguage {
				langConfig = &lc
				break
			}
		}

		if langConfig == nil {
			continue
		}

		res = append(res, &response.GetTopic4StudentResponse4Web{
			ID:          t.ID.Hex(),
			IsPublished: t.IsPublished,
			Title:       langConfig.Title,
		})
	}

	return res
}
