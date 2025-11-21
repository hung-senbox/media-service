package mapper

import (
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/pkg/constants"
)

func ToVocabulariesResponses4Web(vocabularies []model.Vocabulary) []*response.VocabularyResponse4Web {
	var result = make([]*response.VocabularyResponse4Web, 0)

	for _, v := range vocabularies {
		resp := &response.VocabularyResponse4Web{
			ID:          v.ID.Hex(),
			IsPublished: v.IsPublished,
		}

		var langs []response.VocabularyMessageLanguageEntry
		for _, lc := range v.LanguageConfig {

			entry := response.VocabularyMessageLanguageEntry{
				LanguageID: int(lc.LanguageID),
				Contents: response.VocabularyLanguageContents{
					FileName:    lc.FileName,
					Title:       lc.Title,
					Note:        lc.Note,
					Description: lc.Description,
				},
			}

			// map audio
			entry.Contents.Audio = response.VocabularyMediaContent{
				UploadedURL: lc.Audio.UploadedUrl,
				LinkURL:     lc.Audio.LinkUrl,
				StartTime:   strPtr(trimQuotes(lc.Audio.StartTime)),
				EndTime:     strPtr(trimQuotes(lc.Audio.EndTime)),
			}

			// map video
			entry.Contents.Video = response.VocabularyMediaContent{
				UploadedURL: lc.Video.UploadedUrl,
				LinkURL:     lc.Video.LinkUrl,
				StartTime:   strPtr(trimQuotes(lc.Video.StartTime)),
				EndTime:     strPtr(trimQuotes(lc.Video.EndTime)),
			}

			// map images slice → object
			imgMap := make(map[string]response.VocabularyImgEntry)
			for _, img := range lc.Images {
				uploaded := img.UploadedUrl
				imgMap[img.ImageType] = response.VocabularyImgEntry{
					UploadedURL: &uploaded,
					LinkURL:     img.LinkUrl,
				}
			}
			entry.Contents.Images = imgMap

			langs = append(langs, entry)
		}

		mainImageUrl := ""
		if len(v.LanguageConfig) > 0 {
			for _, lc := range v.LanguageConfig {
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

func ToVocabulariesResponses4App(vocabularies []*model.Vocabulary, appLanguage uint) []*response.GetVocabularyResponse4App {
	var res = make([]*response.GetVocabularyResponse4App, 0)

	for _, v := range vocabularies {
		// is published = false
		if !v.IsPublished {
			continue
		}
		// chọn language config
		var langConfig *model.VocabularyLanguageConfig
		for _, lc := range v.LanguageConfig {
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

		res = append(res, &response.GetVocabularyResponse4App{
			ID:           v.ID.Hex(),
			IsPublished:  v.IsPublished,
			Title:        langConfig.Title,
			MainImageUrl: mainImageUrl,
		})
	}

	return res
}
