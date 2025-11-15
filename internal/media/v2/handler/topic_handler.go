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
	var req request.UploadTopicRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	err := h.service.UploadTopic(c.UserContext(), req)
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

	studentID := c.Query("student_id")

	res, err := h.service.GetTopics4Web(c.UserContext(), studentID)
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
