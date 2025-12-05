package mapper

import (
	"context"
	"media-service/helper"
	"media-service/internal/gateway"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
	"time"
)

func ToGetTopicResourceResponses(
	ctx context.Context,
	orgID string,
	topicResources []*model.TopicResource,
	topicRepository repository.TopicRepository,
	userGw gateway.UserGateway,
	s3Service s3.Service,
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
			if url, err := s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
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
	s3Service s3.Service,
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
		if url, err := s3Service.Get(ctx, topicResource.ImageKey, nil); err == nil && url != nil {
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

func ToGetTopicResourcesResponse4Web(
	ctx context.Context,
	topicResources *model.TopicResource,
	resourceImageUrl string,
	topic *response.TopicResponse2Assign4Web,
) *response.GetTopicResourcesResponse4Web {

	loc := time.FixedZone("GMT+7", 7*60*60)
	return &response.GetTopicResourcesResponse4Web{
		ID:        topicResources.ID.Hex(),
		FileName:  topicResources.FileName,
		ImageUrl:  resourceImageUrl,
		CreatedAt: topicResources.CreatedAt,
		PicID:     topicResources.CreatedAt.In(loc).Format("02 Jan 2006 15:04"),
		Topic:     topic,
	}
}

func ToGetTopicResourcesResponse4WebV2(
	topicResources []*model.TopicResource,
) []*response.GetTopicResourcesResponse4WebV2 {
	loc := time.FixedZone("GMT+7", 7*60*60)

	// group theo created_at (theo ngày)
	grouped := make(map[string][]*response.TopicResourceResponseV2)
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}

		dateKey := tr.CreatedAt.In(loc).Format("02 Jan 2006")
		pic := &response.TopicResourceResponseV2{
			ID:        tr.ID.Hex(),
			TopicID:   tr.TopicID,
			ImageKey:  tr.ImageKey,
			FileName:  tr.FileName,
			CreatedAt: tr.CreatedAt,
			PicID:     tr.CreatedAt.In(loc).Format("02 Jan 2006 15:04"),
		}

		grouped[dateKey] = append(grouped[dateKey], pic)
	}

	res := make([]*response.GetTopicResourcesResponse4WebV2, 0, len(grouped))
	for date, pictures := range grouped {
		res = append(res, &response.GetTopicResourcesResponse4WebV2{
			Date:     date,
			Pictures: pictures,
		})
	}

	return res
}

func ToGetOutputTopicResourcesResponse4Web(
	topicResources []*model.TopicResource,
) []*response.GetTopicResourcesResponse4WebV2 {
	loc := time.FixedZone("GMT+7", 7*60*60)

	// group theo created_at (theo ngày)
	grouped := make(map[string][]*response.TopicResourceResponseV2)
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if !tr.IsOutput {
			continue
		}

		dateKey := tr.CreatedAt.In(loc).Format("02 Jan 2006")
		pic := &response.TopicResourceResponseV2{
			ID:        tr.ID.Hex(),
			TopicID:   tr.TopicID,
			ImageKey:  tr.ImageKey,
			FileName:  tr.FileName,
			CreatedAt: tr.CreatedAt,
			PicID:     tr.CreatedAt.In(loc).Format("02 Jan 2006 15:04"),
		}

		grouped[dateKey] = append(grouped[dateKey], pic)
	}

	res := make([]*response.GetTopicResourcesResponse4WebV2, 0, len(grouped))
	for date, pictures := range grouped {
		res = append(res, &response.GetTopicResourcesResponse4WebV2{
			Date:     date,
			Pictures: pictures,
		})
	}

	return res
}

func ToGetTopicResourcesResponse4App(
	ctx context.Context,
	topicResources *model.TopicResource,
	resourceImageUrl string,
	topic response.GetTopicResponse4App,
) *response.GetTopicResourcesResponse4App {

	return &response.GetTopicResourcesResponse4App{
		ID:        topicResources.ID.Hex(),
		FileName:  topicResources.FileName,
		ImageUrl:  resourceImageUrl,
		CreatedAt: topicResources.CreatedAt,
		PicID:     topicResources.CreatedAt.Format("02 Jan 2006 15:04"),
		Topic:     topic,
	}
}
