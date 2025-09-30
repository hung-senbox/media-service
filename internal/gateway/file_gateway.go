package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"media-service/internal/gateway/dto"
	gw_request "media-service/internal/gateway/dto/request"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/pkg/constants"

	"github.com/hashicorp/consul/api"
)

type FileGateway interface {
	UploadImage(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadImageResponse, error)
	UploadVideo(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadVideoResponse, error)
	UploadAudio(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadAudioResponse, error)
}

type fileGateway struct {
	serviceName string
	consul      *api.Client
}

func NewFileGateway(serviceName string, consulClient *api.Client) FileGateway {
	return &fileGateway{
		serviceName: serviceName,
		consul:      consulClient,
	}
}

func (g *fileGateway) UploadImage(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadImageResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("POST", "/v1/gateway/images/upload", req)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[gw_response.UploadImageResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload image fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) UploadVideo(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadVideoResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("POST", "/v1/gateway/videos/upload", req)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[gw_response.UploadVideoResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload video fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) UploadAudio(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadAudioResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("POST", "/v1/gateway/audios/upload", req)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[gw_response.UploadAudioResponse]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway upload audio fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
