package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	service service.TopicService
}

func NewTopicHandler(service service.TopicService) *TopicHandler {
	return &TopicHandler{service: service}
}

func (h *TopicHandler) CreateParentTopic(c *gin.Context) {
	var req request.CreateTopicRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	res, err := h.service.CreateParentTopic(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "waiting for upload file", res)
}

func (h TopicHandler) GetPregressUpload(c *gin.Context) {
	topicID := c.Param("topic_id")
	if topicID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	res, err := h.service.GetUploadProgress(c.Request.Context(), c.Param("topic_id"))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get progress upload success", res)
}

func (h TopicHandler) GetParentTopics4Web(c *gin.Context) {
	res, err := h.service.GetParentTopics4Web(c.Request.Context())
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get topics success", res)
}

func (h TopicHandler) GetParentTopic4Web(c *gin.Context) {
	topicID := c.Param("topic_id")
	if topicID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	res, err := h.service.GetParentTopic4Web(c.Request.Context(), topicID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "get topic success", res)
}
