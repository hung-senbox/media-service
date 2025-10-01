package handler

import (
	"media-service/helper"
	"media-service/internal/media/dto/request"
	"media-service/internal/media/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	service service.TopicService
}

func NewTopicHandler(service service.TopicService) *TopicHandler {
	return &TopicHandler{service: service}
}

func (h *TopicHandler) UploadTopic(c *gin.Context) {
	var req request.UploadTopicRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	res, err := h.service.UploadTopic(c.Request.Context(), req)
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
