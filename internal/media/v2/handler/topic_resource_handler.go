package handler

import (
	"fmt"
	"media-service/helper"
	"media-service/internal/media/v2/dto/request"
	"media-service/internal/media/v2/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TopicResourceHandler struct {
	topicResourceService service.TopicResourceService
}

func NewTopicResourceHandler(topicResourceService service.TopicResourceService) *TopicResourceHandler {
	return &TopicResourceHandler{topicResourceService: topicResourceService}
}

func (h *TopicResourceHandler) CreateTopicResource(c *fiber.Ctx) error {
	var req request.CreateTopicResourceRequest

	// Parse multipart form for file upload
	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}

	// Validate required fields
	if req.TopicID == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("topic_id is required"), helper.ErrInvalidRequest)
	}
	if req.StudentID == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("student_id is required"), helper.ErrInvalidRequest)
	}
	var fileName string
	if req.FileName == "" {
		fileName = helper.GenerateFileName()
		req.FileName = fileName
	}
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("file is required"), helper.ErrInvalidRequest)
	}
	req.File = file

	res, err := h.topicResourceService.CreateTopicResource(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "create topic resource success", res)
}

func (h *TopicResourceHandler) GetTopicResources(c *fiber.Ctx) error {
	topicID := c.Query("topic_id")
	studentID := c.Query("student_id")
	orgID := c.Query("organization_id")

	res, err := h.topicResourceService.GetTopicResources(c.UserContext(), topicID, studentID, orgID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}

func (h *TopicResourceHandler) GetTopicResource(c *fiber.Ctx) error {
	topicResourceID := c.Params("topic_resource_id")
	orgID := c.Query("organization_id")
	if topicResourceID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.topicResourceService.GetTopicResource(c.UserContext(), topicResourceID, orgID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get topic resource success", res)
}

func (h *TopicResourceHandler) UpdateTopicResource(c *fiber.Ctx) error {
	topicResourceID := c.Params("topic_resource_id")
	if topicResourceID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	var req request.UpdateTopicResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	res, err := h.topicResourceService.UpdateTopicResource(c.UserContext(), topicResourceID, req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "update topic resource success", res)
}

func (h *TopicResourceHandler) DeleteTopicResource(c *fiber.Ctx) error {
	topicResourceID := c.Params("topic_resource_id")
	if topicResourceID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	err := h.topicResourceService.DeleteTopicResource(c.UserContext(), topicResourceID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "delete topic resource success", nil)
}

func (h *TopicResourceHandler) GetTopicResourcesByTopic4Web(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.topicResourceService.GetTopicResourcesByTopic4Web(c.UserContext(), topicID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}

func (h *TopicResourceHandler) GetTopicResourcesByTopicAndStudent4Web(c *fiber.Ctx) error {
	topicID := c.Params("topic_id")
	studentID := c.Params("student_id")
	if topicID == "" || studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.topicResourceService.GetTopicResourcesByTopicAndStudent4Web(c.UserContext(), topicID, studentID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}

func (h *TopicResourceHandler) SetOutputTopicResource(c *fiber.Ctx) error {
	var req request.SetOutputTopicResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
	}
	err := h.topicResourceService.SetOutputTopicResource(c.UserContext(), req)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "set output topic resource success", nil)
}

func (h *TopicResourceHandler) GetOutputResources4Web(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	topicID := c.Params("topic_id")
	if topicID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.topicResourceService.GetOutputResources4Web(c.UserContext(), topicID, studentID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get output resources success", res)
}

func (h *TopicResourceHandler) GetOutputResources4App(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	day := c.Query("day")
	month := c.Query("month")
	year := c.Query("year")

	if day != "" && month == "" && year == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("month and year are required"), helper.ErrInvalidRequest)
	}

	if month != "" && year == "" {
		return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("year is required"), helper.ErrInvalidRequest)
	}

	// topicID
	topicID := c.Query("topic_id")

	// Default: no month/year filter
	dayInt := 0
	monthInt := 0
	yearInt := 0

	// If both provided, validate and parse; otherwise skip filtering
	if day != "" {
		d, err := strconv.Atoi(day)
		if err != nil {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("day must be an integer"), helper.ErrInvalidRequest)
		}
		if d < 1 || d > 31 {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("day must be between 1 and 31"), helper.ErrInvalidRequest)
		}
		dayInt = d
	}
	if month != "" && year != "" {
		m, err := strconv.Atoi(month)
		if err != nil {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("month must be an integer"), helper.ErrInvalidRequest)
		}
		y, err := strconv.Atoi(year)
		if err != nil {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("year must be an integer"), helper.ErrInvalidRequest)
		}
		if m < 1 || m > 12 {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("month must be between 1 and 12"), helper.ErrInvalidRequest)
		}
		if y < 1 {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("year must be greater than 0"), helper.ErrInvalidRequest)
		}
		monthInt = m
		yearInt = y
	}

	if year != "" {
		y, err := strconv.Atoi(year)
		if err != nil {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("year must be an integer"), helper.ErrInvalidRequest)
		}
		if y < 1 {
			return helper.SendError(c, http.StatusBadRequest, fmt.Errorf("year must be greater than 0"), helper.ErrInvalidRequest)
		}
		yearInt = y
	}

	res, err := h.topicResourceService.GetOutputResources4App(c.UserContext(), studentID, dayInt, monthInt, yearInt, topicID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get output resources success", res)
}

func (h *TopicResourceHandler) OffOutputTopicResource(c *fiber.Ctx) error {
	topicResourceID := c.Params("topic_resource_id")
	if topicResourceID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	err := h.topicResourceService.OffOutputTopicResource(c.UserContext(), topicResourceID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "off output topic resource success", nil)
}

func (h *TopicResourceHandler) GetTopicResourcesByStudent4Web(c *fiber.Ctx) error {
	studentID := c.Params("student_id")
	if studentID == "" {
		return helper.SendError(c, http.StatusBadRequest, nil, helper.ErrInvalidRequest)
	}
	res, err := h.topicResourceService.GetTopicResourcesByStudent4Web(c.UserContext(), studentID)
	if err != nil {
		return helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
	}
	return helper.SendSuccess(c, http.StatusOK, "get topic resources success", res)
}
