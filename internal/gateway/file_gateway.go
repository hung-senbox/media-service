package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"media-service/internal/gateway/dto"
	gw_request "media-service/internal/gateway/dto/request"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/pkg/constants"
	"mime/multipart"

	"github.com/hashicorp/consul/api"
)

type FileGateway interface {
	UploadImage(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadImageResponse, error)
	UploadVideo(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadVideoResponse, error)
	UploadAudio(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadAudioResponse, error)
	DeleteVideo(ctx context.Context, videoKey string) error
	DeleteAudio(ctx context.Context, audioKey string) error
	DeleteImage(ctx context.Context, imageKey string) error
	GetImageUrl(ctx context.Context, req gw_request.GetFileUrlRequest) (*string, error)
	GetVideoUrl(ctx context.Context, req gw_request.GetFileUrlRequest) (*string, error)
	GetAudioUrl(ctx context.Context, req gw_request.GetFileUrlRequest) (*string, error)
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

func buildMultipartBody(req gw_request.UploadFileRequest) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// --- add file ---
	if req.File != nil {
		file, err := req.File.Open()
		if err != nil {
			return nil, "", fmt.Errorf("open file fail: %w", err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("file", req.File.Filename)
		if err != nil {
			return nil, "", fmt.Errorf("create form file fail: %w", err)
		}
		if _, err := io.Copy(part, file); err != nil {
			return nil, "", fmt.Errorf("copy file fail: %w", err)
		}
	}

	// --- add text fields ---
	_ = writer.WriteField("folder", req.Folder)
	_ = writer.WriteField("file_name", req.FileName)
	_ = writer.WriteField("mode", req.Mode)
	if req.ImageName != "" {
		_ = writer.WriteField("image_name", req.ImageName)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("close writer fail: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}

// --- Upload Image ---
func (g *fileGateway) UploadImage(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadImageResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	body, contentType, err := buildMultipartBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/images/upload", body, contentType)
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

// --- Upload Video ---
func (g *fileGateway) UploadVideo(ctx context.Context, req gw_request.UploadFileRequest) (*gw_response.UploadVideoResponse, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	body, contentType, err := buildMultipartBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/videos/upload", body, contentType)
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

	// multipart body
	body, contentType, err := buildMultipartBody(req)
	if err != nil {
		return nil, err
	}

	resp, err := client.CallWithMultipart("POST", "/v1/gateway/audios/upload", body, contentType)
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

func (g *fileGateway) DeleteAudio(ctx context.Context, audioKey string) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	resp, err := client.Call("DELETE", "/v1/gateway/audios/"+audioKey, nil)
	if err != nil {
		return err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete audio fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) DeleteVideo(ctx context.Context, videoKey string) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	resp, err := client.Call("DELETE", "/v1/gateway/videos/"+videoKey, nil)
	if err != nil {
		return err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete audio fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) DeleteImage(ctx context.Context, imageKey string) error {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return err
	}

	resp, err := client.Call("DELETE", "/v1/gateway/images/"+imageKey, nil)
	if err != nil {
		return err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return fmt.Errorf("call gateway delete image fail: %s", gwResp.Message)
	}

	return nil
}

func (g *fileGateway) GetImageUrl(ctx context.Context, req gw_request.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("POST", "/v1/gateway/images/get-url", req)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get image fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) GetAudioUrl(ctx context.Context, req gw_request.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("POST", "/v1/gateway/audios/get-url", req)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get audio fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}

func (g *fileGateway) GetVideoUrl(ctx context.Context, req gw_request.GetFileUrlRequest) (*string, error) {
	token, ok := ctx.Value(constants.Token).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	client, err := NewGatewayClient(g.serviceName, token, g.consul, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Call("POST", "/v1/gateway/videos/get-url", req)
	if err != nil {
		return nil, err
	}

	var gwResp dto.APIGateWayResponse[string]
	if err := json.Unmarshal(resp, &gwResp); err != nil {
		return nil, fmt.Errorf("unmarshal response fail: %w", err)
	}

	if gwResp.StatusCode != 200 {
		return nil, fmt.Errorf("call gateway get video fail: %s", gwResp.Message)
	}

	return &gwResp.Data, nil
}
