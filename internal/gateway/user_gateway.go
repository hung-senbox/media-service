package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway/dto/response"
	"media-service/logger"
	"media-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type UserGateway interface {
	GetCurrentUser(ctx context.Context) (*response.CurrentUser, error)
	GetUserByTeacher(ctx context.Context, teacherID string) (*response.CurrentUser, error)
	GetStudentInfo(ctx context.Context, studentID string) (*response.StudentResponse, error)
	GetTeacherInfo(ctx context.Context, teacherID string) (*response.TeacherResponse, error)
	GetTeacherByUserAndOrganization(ctx context.Context, userID, organizationID string) (*response.TeacherResponse, error)
	GetStaffByUserAndOrganization(ctx context.Context, userID, organizationID string) (*response.StaffResponse, error)
	GetParentByUser(ctx context.Context, userID string) (*response.ParentResponse, error)
}

type userGatewayImpl struct {
	serviceName string
	consul      *api.Client
}

func NewUserGateway(serviceName string, consulClient *api.Client) UserGateway {
	return &userGatewayImpl{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

// GetCurrentUser
func (g *userGatewayImpl) GetCurrentUser(ctx context.Context) (*response.CurrentUser, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		logger.WriteLogEx("warn", "token not found in context", nil)
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		logger.WriteLogEx("error", "init GatewayClient fail", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/user/current-user", nil, headers)
	if err != nil {
		logger.WriteLogEx("error", "call API user fail", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("call API user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		logger.WriteLogEx("error", "unmarshal response fail", map[string]any{
			"error": string(resp),
		})
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		logger.WriteLogEx("warn", "gateway error", map[string]any{
			"status_code": gwResp.StatusCode,
			"message":     gwResp.Message,
		})
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetStudentInfo(ctx context.Context, studentID string) (*response.StudentResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/students/"+studentID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API student fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.StudentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetTeacherInfo(ctx context.Context, teacherID string) (*response.TeacherResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/teachers/"+teacherID, nil, headers)

	if err != nil {
		return nil, fmt.Errorf("call API teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.TeacherResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetTeacherByUserAndOrganization(ctx context.Context, userID, organizationID string) (*response.TeacherResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/teachers/organization/"+organizationID+"/user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.TeacherResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetStaffByUserAndOrganization(ctx context.Context, userID, organizationID string) (*response.StaffResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/staffs/organization/"+organizationID+"/user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.StaffResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetUserByTeacher(ctx context.Context, teacherID string) (*response.CurrentUser, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/users/teacher/"+teacherID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API user by teacher fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *userGatewayImpl) GetParentByUser(ctx context.Context, userID string) (*response.ParentResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)

	resp, err := client.Call("GET", "/v1/gateway/parents/get-by-user/"+userID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API parent fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp response.APIGateWayResponse[response.ParentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
