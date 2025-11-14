package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/pkg/constants"
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
			entry.Contents.Audio = response.MediaContent{
				UploadedURL: lc.Audio.UploadedUrl,
				LinkURL:     lc.Audio.LinkUrl,
				StartTime:   strPtr(trimQuotes(lc.Audio.StartTime)),
				EndTime:     strPtr(trimQuotes(lc.Audio.EndTime)),
			}

			// map video
			entry.Contents.Video = response.MediaContent{
				UploadedURL: lc.Video.UploadedUrl,
				LinkURL:     lc.Video.LinkUrl,
				StartTime:   strPtr(trimQuotes(lc.Video.StartTime)),
				EndTime:     strPtr(trimQuotes(lc.Video.EndTime)),
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

		mainImageUrl := ""
		if len(t.LanguageConfig) > 0 {
			for _, lc := range t.LanguageConfig {
				if lc.LanguageID == 1 {
					for _, img := range lc.Images {
						if img.ImageType == string(constants.TopicImageTypeBM) && img.UploadedUrl != "" {
							mainImageUrl = img.UploadedUrl
							break
						}
					}
					break
				}
			}
		}
		resp.MainImageUrl = mainImageUrl
		resp.MessageLangs = langs
		result = append(result, resp)
	}

	return result
}

func trimQuotes(s string) string {
	return strings.Trim(s, "\"")
}

func strPtr(s string) string {
	if s == "" {
		return ""
	}
	return s
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
		entry.Contents.Audio = response.MediaContent{
			UploadedURL: lc.Audio.UploadedUrl,
			LinkURL:     lc.Audio.LinkUrl,
			StartTime:   strPtr(trimQuotes(lc.Audio.StartTime)),
			EndTime:     strPtr(trimQuotes(lc.Audio.EndTime)),
		}

		// map video
		entry.Contents.Video = response.MediaContent{
			UploadedURL: lc.Video.UploadedUrl,
			LinkURL:     lc.Video.LinkUrl,
			StartTime:   strPtr(trimQuotes(lc.Video.StartTime)),
			EndTime:     strPtr(trimQuotes(lc.Video.EndTime)),
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

		mainImageUrl := ""
		if len(langConfig.Images) > 0 {
			for _, img := range langConfig.Images {
				if img.ImageType == string(constants.TopicImageTypeBM) {
					mainImageUrl = img.UploadedUrl
					break
				}
			}
		}

		res = append(res, &response.GetTopic4StudentResponse4App{
			ID:           t.ID.Hex(),
			IsPublished:  t.IsPublished,
			Title:        langConfig.Title,
			MainImageUrl: mainImageUrl,
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

		mainImageUrl := ""
		if len(langConfig.Images) > 0 {
			for _, img := range langConfig.Images {
				if img.ImageType == string(constants.TopicImageTypeBM) {
					mainImageUrl = img.UploadedUrl
					break
				}
			}
		}

		res = append(res, &response.GetTopic4StudentResponse4Web{
			ID:           t.ID.Hex(),
			IsPublished:  t.IsPublished,
			Title:        langConfig.Title,
			MainImageUrl: mainImageUrl,
		})
	}

	return res
}

func ToTopic4StudentResponses4Gw(topics []model.Topic, appLanguage uint) []*response.GetTopic4StudentResponse4Gw {
	var res = make([]*response.GetTopic4StudentResponse4Gw, 0)

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

		mainImageUrl := ""
		if len(langConfig.Images) > 0 {
			for _, img := range langConfig.Images {
				if img.ImageType == string(constants.TopicImageTypeBM) {
					mainImageUrl = img.UploadedUrl
					break
				}
			}
		}

		res = append(res, &response.GetTopic4StudentResponse4Gw{
			ID:           t.ID.Hex(),
			IsPublished:  t.IsPublished,
			Title:        langConfig.Title,
			MainImageUrl: mainImageUrl,
		})
	}

	return res
}

func ToTopicResponses4GW(topic *model.Topic, appLanguage uint) *response.TopicResponse4GW {
	if topic == nil {
		return nil
	}
	// test
	// Tìm config ngôn ngữ tương ứng
	var langConfig *model.TopicLanguageConfig
	for _, lc := range topic.LanguageConfig {
		if lc.LanguageID == appLanguage {
			langConfig = &lc
			break
		}
	}

	if langConfig == nil {
		return nil
	}

	if !topic.IsPublished {
		return nil
	}

	mainImageUrl := ""
	if len(langConfig.Images) > 0 {
		for _, img := range langConfig.Images {
			if img.ImageType == string(constants.TopicImageTypeBM) {
				mainImageUrl = img.UploadedUrl
				break
			}
		}
	}

	return &response.TopicResponse4GW{
		ID:           topic.ID.Hex(),
		Title:        langConfig.Title,
		MainImageUrl: mainImageUrl,
	}
}

func ToTopic2Assign4Web(topics []model.Topic, appLanguage uint) []*response.TopicResponse2Assign4Web {
	var res = make([]*response.TopicResponse2Assign4Web, 0)

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

		mainImageUrl := ""
		if len(langConfig.Images) > 0 {
			for _, img := range langConfig.Images {
				if img.ImageType == string(constants.TopicImageTypeBM) {
					mainImageUrl = img.UploadedUrl
					break
				}
			}
		}

		res = append(res, &response.TopicResponse2Assign4Web{
			ID:           t.ID.Hex(),
			Title:        langConfig.Title,
			MainImageUrl: mainImageUrl,
		})
	}

	return res
}

func ToTopicResponse4App(t *model.Topic, appLanguage uint) *response.GetTopicResponse4App {
	// Nếu topic chưa publish thì bỏ qua
	if !t.IsPublished {
		return nil
	}

	// Tìm language config phù hợp
	var langConfig *model.TopicLanguageConfig
	for _, lc := range t.LanguageConfig {
		if lc.LanguageID == appLanguage {
			langConfig = &lc
			break
		}
	}

	if langConfig == nil {
		return nil
	}

	// Lấy ảnh full_background (nếu có)
	mainImageUrl := ""
	for _, img := range langConfig.Images {
		if img.ImageType == string(constants.TopicImageTypeBM) {
			mainImageUrl = img.UploadedUrl
			break
		}
	}

	// Trả về response
	return &response.GetTopicResponse4App{
		ID:           t.ID.Hex(),
		IsPublished:  t.IsPublished,
		Title:        langConfig.Title,
		MainImageUrl: mainImageUrl,
	}
}
