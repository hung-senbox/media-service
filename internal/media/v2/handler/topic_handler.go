package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TopicHandler struct {
	service service.TopicService
}

func NewTopicHandler(service service.TopicService) *TopicHandler {
	return &TopicHandler{service: service}
}

func (h *TopicHandler) UploadTopic(c *fiber.Ctx) error {
	// Parse multipart form
	_, err := c.MultipartForm()
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	// Build request manually
	req := request.UploadTopicRequest{
		TopicID:     c.FormValue("topic_id"),
		FileName:    c.FormValue("file_name"),
		Title:       c.FormValue("title"),
		Note:        c.FormValue("note"),
		Description: c.FormValue("description"),

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

	err = h.service.UploadTopic(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "upload topic success", nil)
}

func (h TopicHandler) GetPregressUpload(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetUploadProgress(c.UserContext(), c.Params("topic_id"))
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get progress upload success", res)
}

func (h TopicHandler) GetTopics4Web(c *fiber.Ctx) error {

	res, err := h.service.GetTopics4Web(c.UserContext())
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) GetTopic4Web(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetTopic4Web(c.UserContext(), topicID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topic success", res)
}

func (h TopicHandler) GetTopics4Student4App(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetTopics4Student4App(c.UserContext(), studentID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) GetTopic4Gw(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetTopic4Gw(c.UserContext(), topicID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topic success", res)
}

func (h TopicHandler) GetAllTopicsByOrganization4Gw(c *fiber.Ctx) error {
	organizationID := c.Params("organization_id")
	if organizationID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}

	res, err := h.service.GetAllTopicsByOrganization4Gw(c.UserContext(), organizationID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get all topics success", res)
}

func (h TopicHandler) GetTopics4Student4Web(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetTopics4Student4Web(c.UserContext(), studentID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) GetTopics4Student4Gw(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetTopics4Student4Gw(c.UserContext(), studentID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) GetTopics2Assign4Web(c *fiber.Ctx) error {
	res, err := h.service.GetTopics2Assign4Web(c.UserContext())
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) GetTopics4App(c *fiber.Ctx) error {
	organizationID := c.Params("organization_id")
	if organizationID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.service.GetTopics4App(c.UserContext(), organizationID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}

	return helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) DeleteTopicAudioKey(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageID := c.Params("language_id")
	if languageID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageIDUint, err := strconv.ParseUint(languageID, 10, 64)
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	err = h.service.DeleteTopicAudioKey(c.UserContext(), topicID, uint(languageIDUint))
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "delete topic audio key success", nil)
}

func (h TopicHandler) DeleteTopicVideoKey(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageID := c.Params("language_id")
	if languageID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageIDUint, err := strconv.ParseUint(languageID, 10, 64)
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	err = h.service.DeleteTopicVideoKey(c.UserContext(), topicID, uint(languageIDUint))
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "delete topic video key success", nil)
}

func (h TopicHandler) DeleteTopicImageKey(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageID := c.Params("language_id")
	if languageID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	languageIDUint, err := strconv.ParseUint(languageID, 10, 64)
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	imageType := c.Params("image_type")
	if imageType == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	err = h.service.DeleteTopicImageKey(c.UserContext(), topicID, uint(languageIDUint), imageType)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "delete topic image key success", nil)
}
