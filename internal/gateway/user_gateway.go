package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"media-service/internal/gateway/dto"
	"media-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserGateway interface {
	GetCurrentUser(ctx context.Context) (*dto.CurrentUser, error)
	GetUserInfo(ctx context.Context, userID string) (*dto.CurrentUser, error)
	GetTeachersByUser(ctx context.Context, userID string) ([]*dto.TeacherResponse, error)
	GetStaffsByUser(ctx context.Context, userID string) ([]*dto.StaffResponse, error)
	GetTeacherInfo(ctx context.Context, teacherID string) (*dto.TeacherResponse, error)
	GetStaffInfo(ctx context.Context, staffID string) (*dto.StaffResponse, error)
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
func (g *userGatewayImpl) GetCurrentUser(ctx context.Context) (*dto.CurrentUser, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/user/current-user", nil)
	if err != nil {
		return nil, fmt.Errorf("call API user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

// GetStudentInfo
func (g *userGatewayImpl) GetStudentInfo(ctx context.Context, studentID string) (*dto.StudentResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/gateway/students/"+studentID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API student fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.StudentResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

// GetUserInfo
func (g *userGatewayImpl) GetUserInfo(ctx context.Context, userID string) (*dto.CurrentUser, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/gateway/users/"+userID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.CurrentUser]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

// get all teacher by user
func (g *userGatewayImpl) GetTeachersByUser(ctx context.Context, userID string) ([]*dto.TeacherResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/gateway/teachers/get-by-user/"+userID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API get teacher by user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[[]*dto.TeacherResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

// get all staff by user
func (g *userGatewayImpl) GetStaffsByUser(ctx context.Context, userID string) ([]*dto.StaffResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/gateway/staffs/get-by-user/"+userID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API get staff by user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[[]*dto.StaffResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *userGatewayImpl) GetTeacherInfo(ctx context.Context, teacherID string) (*dto.TeacherResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/gateway/teachers/"+teacherID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API get teacher by user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.TeacherResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

// get staff by user
func (g *userGatewayImpl) GetStaffInfo(ctx context.Context, staffID string) (*dto.StaffResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	resp, err := client.Call("GET", "/v1/gateway/staffs/"+staffID, nil)
	if err != nil {
		return nil, fmt.Errorf("call API get teacher by user fail: %w", err)
	}

	// Unmarshal response theo format Gateway
	var gwResp dto.APIGateWayResponse[dto.StaffResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	// Check status_code trả về
	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
