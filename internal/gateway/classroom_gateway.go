package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"media-service/helper"
	"media-service/internal/gateway/dto"
	"media-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type ClassroomGateway interface {
	GetClassroomByID(ctx context.Context, locationID string) (*dto.ClassroomResponse, error)
}

type classroomGateway struct {
	serviceName string
	consul      *api.Client
}

func NewClassroomGateway(serviceName string, consulClient *api.Client) ClassroomGateway {
	return &classroomGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *classroomGateway) GetClassroomByID(ctx context.Context, locationID string) (*dto.ClassroomResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, fmt.Errorf("init GatewayClient fail: %w", err)
	}

	headers := helper.GetHeaders(ctx)
	resp, err := client.Call("GET", "/api/v1/storage/"+locationID, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("call API get classroom by id fail: %w", err)
	}

	var gwResp dto.APIGateWayResponse[dto.ClassroomResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("gateway error: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
