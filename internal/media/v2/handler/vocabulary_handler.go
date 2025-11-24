package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type VocabularyHandler struct {
	vocabularyService service.VocabularyService
}

func NewVocabularyHandler(vocabularyService service.VocabularyService) *VocabularyHandler {
	return &VocabularyHandler{vocabularyService: vocabularyService}
}

func (h *VocabularyHandler) UploadVocabulary(c *fiber.Ctx) error {
	// Parse multipart form
	_, err := c.MultipartForm()
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}

	// Build request manually
	req := request.UploadVocabularyRequest{
		TopicID:      topicID,
		VocabularyID: c.FormValue("vocabulary_id"),
		FileName:     c.FormValue("file_name"),
		Title:        c.FormValue("title"),
		Note:         c.FormValue("note"),
		Description:  c.FormValue("description"),

		// Audio fields
		AudioLinkUrl:   c.FormValue("audio_link_url"),
		AudioStart:     c.FormValue("audio_start_time"),
		AudioEnd:       c.FormValue("audio_end_time"),
		IsDeletedAudio: c.FormValue("is_deleted_audio") == "true",

		// Video fields
		VideoLinkUrl:   c.FormValue("video_link_url"),
		VideoStart:     c.FormValue("video_start_time"),
		VideoEnd:       c.FormValue("video_end_time"),
		IsDeletedVideo: c.FormValue("is_deleted_video") == "true",

		// Image fields
		FullBackgroundLink:       c.FormValue("full_background_link_url"),
		IsDeletedFullBackground:  c.FormValue("is_deleted_full_background") == "true",
		ClearBackgroundLink:      c.FormValue("clear_background_link_url"),
		IsDeletedClearBackground: c.FormValue("is_deleted_clear_background") == "true",
		ClipPartLink:             c.FormValue("clip_part_link_url"),
		IsDeletedClipPart:        c.FormValue("is_deleted_clip_part") == "true",
		DrawingLink:              c.FormValue("drawing_link_url"),
		IsDeletedDrawing:         c.FormValue("is_deleted_drawing") == "true",
		IconLink:                 c.FormValue("icon_link_url"),
		IsDeletedIcon:            c.FormValue("is_deleted_icon") == "true",
		BMLink:                   c.FormValue("bm_link_url"),
		IsDeletedBM:              c.FormValue("is_deleted_bm") == "true",
		SignLangLink:             c.FormValue("sign_lang_link_url"),
		IsDeletedSignLang:        c.FormValue("is_deleted_sign_lang") == "true",
		GifLink:                  c.FormValue("gif_link_url"),
		IsDeletedGif:             c.FormValue("is_deleted_gif") == "true",
		OrderLink:                c.FormValue("order_link_url"),
		IsDeletedOrder:           c.FormValue("is_deleted_order") == "true",
	}

	// Parse uint fields
	if langID := c.FormValue("language_id"); langID != "" {
		if val, err := strconv.ParseUint(langID, 10, 32); err == nil {
			req.LanguageID = uint(val)
		} else {
			return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		}
	}

	if isPublished := c.FormValue("is_published"); isPublished != "" {
		if val, err := strconv.ParseBool(isPublished); err == nil {
			req.IsPublished = val
		} else {
			return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		}
	}

	// Parse file fields
	if audioFile, err := c.FormFile("audio_file"); err == nil {
		req.AudioFile = audioFile
	}
	if videoFile, err := c.FormFile("video_file"); err == nil {
		req.VideoFile = videoFile
	}
	if fullBgFile, err := c.FormFile("full_background_file"); err == nil {
		req.FullBackgroundFile = fullBgFile
	}
	if clearBgFile, err := c.FormFile("clear_background_file"); err == nil {
		req.ClearBackgroundFile = clearBgFile
	}
	if clipPartFile, err := c.FormFile("clip_part_file"); err == nil {
		req.ClipPartFile = clipPartFile
	}
	if drawingFile, err := c.FormFile("drawing_file"); err == nil {
		req.DrawingFile = drawingFile
	}
	if iconFile, err := c.FormFile("icon_file"); err == nil {
		req.IconFile = iconFile
	}
	if bmFile, err := c.FormFile("bm_file"); err == nil {
		req.BMFile = bmFile
	}
	if signLangFile, err := c.FormFile("sign_lang_file"); err == nil {
		req.SignLangFile = signLangFile
	}
	if gifFile, err := c.FormFile("gif_file"); err == nil {
		req.GifFile = gifFile
	}
	if orderFile, err := c.FormFile("order_file"); err == nil {
		req.OrderFile = orderFile
	}

	err = h.vocabularyService.UploadVocabulary(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "upload vocabulary success", nil)
}

func (h *VocabularyHandler) GetVocabularies4Web(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.vocabularyService.GetVocabularies4Web(c.UserContext(), topicID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get vocabularies success", res)
}

func (h *VocabularyHandler) GetVocabularies4Gw(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.vocabularyService.GetVocabularies4Gw(c.UserContext(), topicID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get vocabulary success", res)
}
