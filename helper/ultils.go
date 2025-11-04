package helper

import (
	"context"
	"media-service/internal/media/model"
	"media-service/pkg/constants"
	"mime/multipart"
	"strconv"
	"strings"
)

func ParseAppLanguage(header string, defaultVal uint) uint {
	header = strings.TrimSpace(strings.Trim(header, "\""))
	if val, err := strconv.Atoi(header); err == nil {
		return uint(val)
	}
	return defaultVal
}

func GetHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)

	if lang, ok := ctx.Value(constants.AppLanguage).(uint); ok {
		headers["X-App-Language"] = strconv.Itoa(int(lang))
	}

	return headers
}

func GetAppLanguage(ctx context.Context, defaultVal uint) uint {
	if lang, ok := ctx.Value(constants.AppLanguage).(uint); ok {
		return lang
	}
	return defaultVal
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(constants.UserID).(string); ok {
		return userID
	}
	return ""
}

func IsValidFile(f *multipart.FileHeader) bool {
	return f != nil && f.Size > 0
}

func GetAudioKeyByLanguage(topic *model.Topic, languageID uint) string {
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

func GetVideoKeyByLanguage(topic *model.Topic, languageID uint) string {
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

func GetImageKeyByLanguageAndType(topic *model.Topic, languageID uint, imageType string) string {
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
