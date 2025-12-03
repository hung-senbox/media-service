package helper

import (
	"context"
	"fmt"
	"media-service/internal/media/model"
	"media-service/pkg/constants"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

func RemoveDuplicateString(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func BuildVideoUploaderRedisKey(videoUploaderID string) string {
	return fmt.Sprintf("media_video:upload_status:%s", videoUploaderID)
}

func BuildObjectKeyS3(folder, originalFilename, baseName string) string {
	ext := strings.ToLower(filepath.Ext(originalFilename))
	name := strings.TrimSpace(baseName)
	if name == "" {
		name = strings.TrimSuffix(originalFilename, ext)
	}
	name = strings.NewReplacer(" ", "-", "/", "-", "\\", "-").Replace(name)
	return fmt.Sprintf("%s/%d_%s%s", strings.Trim(folder, "/"), time.Now().UnixNano(), name, ext)
}

func GetVocabularyAudioKeyByLanguage(vocabulary *model.Vocabulary, languageID uint) string {
	for _, lc := range vocabulary.LanguageConfig {
		if lc.LanguageID == languageID {
			if lc.Audio.AudioKey != "" {
				return lc.Audio.AudioKey
			}
			break
		}
	}
	return ""
}

func GetVocabularyVideoKeyByLanguage(vocabulary *model.Vocabulary, languageID uint) string {
	for _, lc := range vocabulary.LanguageConfig {
		if lc.LanguageID == languageID {
			if lc.Video.VideoKey != "" {
				return lc.Video.VideoKey
			}
			break
		}
	}
	return ""
}

func GetVocabularyImageKeyByLanguageAndType(vocabulary *model.Vocabulary, languageID uint, imageType string) string {
	for _, lc := range vocabulary.LanguageConfig {
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

func RemoveDuplicatesString(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GenerateFileName() string {
	return fmt.Sprintf("%d_%s", time.Now().UnixNano(), uuid.New().String())
}
