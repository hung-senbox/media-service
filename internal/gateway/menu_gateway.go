package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"media-service/internal/department/dto/request"
	"media-service/internal/gateway/dto"
	"media-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type MenuGateway interface {
	UploadDepartmentMenu(ctx context.Context, req request.UploadSectionMenuDepartmentRequest) error
	GetDepartmentMenu(ctx context.Context, departmentID string) ([]dto.MenuResponse, error)
	GetImageURL(ctx context.Context, imageKey, mode string) (string, error)
	GetAvatarUrl(ctx context.Context, OwnerID string, OwnerRole string) (string, error)
	UploadDepartmentMenuOrganization(ctx context.Context, req request.UploadDepartmentMenuOrganizationRequest) error
	GetDepartmentMenuOrganization(ctx context.Context, departmentID string, organizationID string) ([]dto.MenuResponse, error)
}

type menuGateway struct {
	serviceName string
	consul      *api.Client
}

func NewMenuGateway(serviceName string, consulClient *api.Client) MenuGateway {
	return &menuGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *menuGateway) UploadDepartmentMenu(ctx context.Context, req request.UploadSectionMenuDepartmentRequest) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	resp, err := client.Call("POST", "/v1/gateway/menus/department", req)
	if err != nil {
		return err
	}

	var gwResp dto.APIGateWayResponse[any]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway upload department menu fail: %s", gwResp.Message)
	}

	return nil
}

func (g *menuGateway) GetDepartmentMenu(ctx context.Context, departmentID string) ([]dto.MenuResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("GET", "/v1/gateway/menus/department/"+departmentID, nil)
	if err != nil {
		return nil, err
	}

	type DepartmentMenuResponse struct {
		Components []dto.MenuResponse `json:"components"`
	}

	var gwResp dto.APIGateWayResponse[DepartmentMenuResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get department menu fail: %s", gwResp.Message)
	}

	// chỉ return components
	return gwResp.Data.Components, nil
}

func (g *menuGateway) GetImageURL(ctx context.Context, imageKey, mode string) (string, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return "", fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return "", err
	}

	reqBody := map[string]string{
		"key":  imageKey,
		"mode": mode,
	}

	resp, err := client.Call("POST", "/v1/gateway/images/get-url", reqBody)
	if err != nil {
		return "", err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get image url fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil
}

func (g *menuGateway) GetAvatarUrl(ctx context.Context, OwnerID string, OwnerRole string) (string, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return "", fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return "", err
	}

	reqBody := map[string]string{
		"owner_id":   OwnerID,
		"owner_role": OwnerRole,
	}

	resp, err := client.Call("POST", "/v1/gateway/images/avatar/get-url", reqBody)
	if err != nil {
		return "", err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return "", fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return "", fmt.Errorf("call gateway get image url fail: %s", gwResp.Message)
	}

	return gwResp.Data, nil

}

func (g *menuGateway) UploadDepartmentMenuOrganization(ctx context.Context, req request.UploadDepartmentMenuOrganizationRequest) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	resp, err := client.Call("POST", "/v1/gateway/menus/department/organization", req)
	if err != nil {
		return err
	}

	var gwResp dto.APIGateWayResponse[any]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway upload department menu fail: %s", gwResp.Message)
	}

	return nil
}

func (g *menuGateway) GetDepartmentMenuOrganization(ctx context.Context, departmentID string, organizationID string) ([]dto.MenuResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("GET", "/v1/gateway/menus/department/"+departmentID+"/organization/"+organizationID, nil)
	if err != nil {
		return nil, err
	}

	type DepartmentMenuResponse struct {
		Components []dto.MenuResponse `json:"components"`
	}

	var gwResp dto.APIGateWayResponse[DepartmentMenuResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get department menu fail: %s", gwResp.Message)
	}

	return gwResp.Data.Components, nil
}
