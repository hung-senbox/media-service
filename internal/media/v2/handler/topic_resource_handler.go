package handler

import (
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TopicResourceHandler struct {
	topicResourceService service.TopicResourceService
}

func NewTopicResourceHandler(topicResourceService service.TopicResourceService) *TopicResourceHandler {
	return &TopicResourceHandler{topicResourceService: topicResourceService}
}

func (h *TopicResourceHandler) CreateTopicResource(c *gin.Context) {
	var req request.CreateTopicResourceRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.CreateTopicResource(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "create topic resource success", res)
}

func (h *TopicResourceHandler) GetTopicResources(c *gin.Context) {
	topicID := c.Query("topic_id")
	studentID := c.Query("student_id")
	orgID := c.Query("organization_id")

	res, err := h.topicResourceService.GetTopicResources(c.Request.Context(), topicID, studentID, orgID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}

func (h *TopicResourceHandler) GetTopicResource(c *gin.Context) {
	topicResourceID := c.Param("topic_resource_id")
	orgID := c.Query("organization_id")
	if topicResourceID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.GetTopicResource(c.Request.Context(), topicResourceID, orgID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "get topic resource success", res)
}

func (h *TopicResourceHandler) UpdateTopicResource(c *gin.Context) {
	topicResourceID := c.Param("topic_resource_id")
	if topicResourceID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	var req request.UpdateTopicResourceRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.UpdateTopicResource(c.Request.Context(), topicResourceID, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "update topic resource success", res)
}

func (h *TopicResourceHandler) DeleteTopicResource(c *gin.Context) {
	topicResourceID := c.Param("topic_resource_id")
	if topicResourceID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	err := h.topicResourceService.DeleteTopicResource(c.Request.Context(), topicResourceID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "delete topic resource success", nil)
}

func (h *TopicResourceHandler) GetTopicResourcesByTopic4Web(c *gin.Context) {
	topicID := c.Param("topic_id")
	if topicID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.GetTopicResourcesByTopic4Web(c.Request.Context(), topicID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}

func (h *TopicResourceHandler) GetTopicResourcesByTopicAndStudent4Web(c *gin.Context) {
	topicID := c.Param("topic_id")
	studentID := c.Param("student_id")
	if topicID == "" || studentID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.GetTopicResourcesByTopicAndStudent4Web(c.Request.Context(), topicID, studentID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}

func (h *TopicResourceHandler) SetOutputTopicResource(c *gin.Context) {
	var req request.SetOutputTopicResourceRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	err := h.topicResourceService.SetOutputTopicResource(c.Request.Context(), req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "set output topic resource success", nil)
}

func (h *TopicResourceHandler) GetOutputResources4Web(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	topicID := c.Param("topic_id")
	if topicID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.GetOutputResources4Web(c.Request.Context(), topicID, studentID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "get output resources success", res)
}

func (h *TopicResourceHandler) GetOutputResources4App(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	month := c.Query("month")
	year := c.Query("year")
	if month == "" || year == "" {
		helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
		return
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}
	res, err := h.topicResourceService.GetOutputResources4App(c.Request.Context(), studentID, monthInt, yearInt)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, http.StatusOK, "get output resources success", res)
}
