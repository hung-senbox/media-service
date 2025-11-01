package mapper

import (
	"context"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_request "media-service/internal/gateway/dto/request"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/repository"
)

func ToGetTopicResourceResponses(
	ctx context.Context,
	orgID string,
	topicResources []*model.TopicResource,
	topicRepository repository.TopicRepository,
	userGw gateway.UserGateway,
	fileGateway gateway.FileGateway,
) []*response.GetTopicResourceResponse {
	if len(topicResources) == 0 {
		return []*response.GetTopicResourceResponse{}
	}

	res := make([]*response.GetTopicResourceResponse, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}

		// reset per item
		var imageUrl string
		var student *gw_response.StudentResponse
		var createdBy *gw_response.TeacherResponse
		var topicResp *response.TopicResponse2Assign4Web

		if tr.ImageKey != "" {
			if url, err := fileGateway.GetImageUrl(ctx, gw_request.GetFileUrlRequest{Key: tr.ImageKey, Mode: "private"}); err == nil && url != nil {
				imageUrl = *url
			}
		}

		if tr.StudentID != "" {
			if studentData, err := userGw.GetStudentInfo(ctx, tr.StudentID); err == nil {
				student = studentData
			}
		}

		if tr.CreatedBy != "" {
			if createdByData, err := userGw.GetTeacherByUserAndOrganization(ctx, tr.CreatedBy, orgID); err == nil {
				createdBy = createdByData
			}
		}

		if tr.TopicID != "" {
			if topic, err := topicRepository.GetByID(ctx, tr.TopicID); err == nil && topic != nil {
				var title string
				var mainImageUrl string
				appLang := helper.GetAppLanguage(ctx, 1)
				for _, lc := range topic.LanguageConfig {
					if lc.LanguageID == appLang {
						title = lc.Title
						for _, img := range lc.Images {
							if img.ImageType == "full_background" {
								mainImageUrl = img.UploadedUrl
								break
							}
						}
						break
					}
				}
				topicResp = &response.TopicResponse2Assign4Web{ID: topic.ID.Hex(), Title: title, MainImageUrl: mainImageUrl}
			}
		}

		res = append(res, &response.GetTopicResourceResponse{
			ID:        tr.ID.Hex(),
			Topic:     topicResp,
			Student:   student,
			ImageUrl:  imageUrl,
			FileName:  tr.FileName,
			CreatedBy: createdBy,
			CreatedAt: tr.CreatedAt,
			UpdatedAt: tr.UpdatedAt,
		})
	}

	return res
}


func ToGetTopicResourceResponse(
	ctx context.Context,
	orgID string,
	topicResource *model.TopicResource,
	topicRepository repository.TopicRepository,
	userGw gateway.UserGateway,
	fileGateway gateway.FileGateway,
	) *response.GetTopicResourceResponse {
	
	var imageUrl string
	var student *gw_response.StudentResponse
	var createdBy *gw_response.TeacherResponse
	var topicResp *response.TopicResponse2Assign4Web

	if topicResource.StudentID != "" {
		if studentData, err := userGw.GetStudentInfo(ctx, topicResource.StudentID); err == nil {
			student = studentData
		}
	}

	if topicResource.CreatedBy != "" {
		if createdByData, err := userGw.GetTeacherByUserAndOrganization(ctx, topicResource.CreatedBy, orgID); err == nil {
			createdBy = createdByData
		}
	}

	if topicResource.TopicID != "" {
		if topic, err := topicRepository.GetByID(ctx, topicResource.TopicID); err == nil && topic != nil {
			var title string
			var mainImageUrl string
			appLang := helper.GetAppLanguage(ctx, 1)
			for _, lc := range topic.LanguageConfig {
				if lc.LanguageID == appLang {
					title = lc.Title
					for _, img := range lc.Images {
						if img.ImageType == "full_background" {
							mainImageUrl = img.UploadedUrl
							break
						}
					}
					break
				}
			}
			topicResp = &response.TopicResponse2Assign4Web{ID: topic.ID.Hex(), Title: title, MainImageUrl: mainImageUrl}
		}
	}

	if topicResource.ImageKey != "" {
		if url, err := fileGateway.GetImageUrl(ctx, gw_request.GetFileUrlRequest{Key: topicResource.ImageKey, Mode: "private"}); err == nil && url != nil {
			imageUrl = *url
		}
	}

	return &response.GetTopicResourceResponse{
		ID:        topicResource.ID.Hex(),
		Topic:     topicResp,
		Student:   student,
		ImageUrl:  imageUrl,
		FileName:  topicResource.FileName,
		CreatedBy: createdBy,
		CreatedAt: topicResource.CreatedAt,
		UpdatedAt: topicResource.UpdatedAt,
	}
}
