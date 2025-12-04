package usecase

import (
	"context"
	"fmt"
	"media-service/helper"
	gw_response "media-service/internal/gateway/dto/response"
	"media-service/internal/media/model"
	"media-service/internal/media/v2/dto/response"
	"media-service/internal/media/v2/mapper"
	"media-service/internal/media/v2/repository"
	"media-service/internal/s3"
	"media-service/logger"
	"media-service/pkg/constants"
)

type GetTopicResourcesWebUseCase interface {
	GetTopicResourcesByTopicAndStudent4Web(ctx context.Context, topicID string, studentID string) ([]*response.GetTopicResourcesResponse4WebV2, error)
	GetTopicResourcesByTopic4Web(ctx context.Context, topicID string) ([]*response.GetTopicResourcesResponse4Web, error)
	GetOutputResources4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error)
	GetTopicResourcesByStudent4Web(ctx context.Context, studentID string) ([]*response.GetTopicResourcesResponseByStudent4Web, error)
}

type getTopicResourcesWebUseCase struct {
	topicResourceRepository repository.TopicResourceRepository
	topicRepository         repository.TopicRepository
	s3Service               s3.Service
}

func NewGetTopicResourcesWebUseCase(
	topicResourceRepository repository.TopicResourceRepository,
	topicRepository repository.TopicRepository,
	s3Service s3.Service,
) GetTopicResourcesWebUseCase {
	return &getTopicResourcesWebUseCase{
		topicResourceRepository: topicResourceRepository,
		topicRepository:         topicRepository,
		s3Service:               s3Service,
	}
}

func (uc *getTopicResourcesWebUseCase) GetTopicResourcesByTopicAndStudent4Web(ctx context.Context, topicID string, studentID string) ([]*response.GetTopicResourcesResponse4WebV2, error) {
	currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
		return nil, fmt.Errorf("access denied")
	}
	// get topic de kiem tra isAllPic
	topic, _ := uc.topicRepository.GetByID(ctx, topicID)
	if topic == nil {
		return nil, fmt.Errorf("topic not found")
	}
	var topicResources []*model.TopicResource
	if topic.IsAllPic {
		resources, err := uc.topicResourceRepository.GetTopicResouresByStudentID(ctx, studentID)
		if err != nil {
			return nil, err
		}
		topicResources = resources
	} else {
		resources, err := uc.topicResourceRepository.GetTopicResouresByTopicAndStudent(ctx, topicID, studentID)
		if err != nil {
			return nil, err
		}
		topicResources = resources
	}

	// result := make([]*response.GetTopicResourcesResponse4Web, 0, len(topicResources))
	// for _, tr := range topicResources {
	// 	if tr == nil {
	// 		continue
	// 	}
	// 	var imageUrl string
	// 	if tr.ImageKey != "" {
	// 		if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
	// 			imageUrl = *url
	// 		}
	// 	}
	// 	result = append(result, mapper.ToGetTopicResourcesResponse4Web(ctx, tr, imageUrl, nil))
	// }
	// return result, nil

	res := mapper.ToGetTopicResourcesResponse4WebV2(topicResources)

	// loop res de lay topic va image url
	for _, res := range res {
		for _, pic := range res.Pictures {
			var imageUrl string
			if pic.ImageUrl == "" {
				if url, err := uc.s3Service.Get(ctx, pic.ImageKey, nil); err == nil && url != nil {
					imageUrl = *url
				}
			}
			pic.ImageUrl = imageUrl
		}
	}
	return res, nil
}

func (uc *getTopicResourcesWebUseCase) GetTopicResourcesByTopic4Web(ctx context.Context, topicID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	// currentUser, _ := ctx.Value(constants.CurrentUserKey).(*gw_response.CurrentUser)
	// if currentUser.IsSuperAdmin || currentUser.OrganizationAdmin.ID == "" {
	// 	return nil, fmt.Errorf("access denied")
	// }
	topicResources, err := uc.topicResourceRepository.GetTopicResouresByTopic(ctx, topicID)
	if err != nil {
		return nil, err
	}

	result := make([]*response.GetTopicResourcesResponse4Web, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.IsOutput {
			var imageUrl string
			if tr.ImageKey != "" {
				if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
					imageUrl = *url
				}
			}
			result = append(result, mapper.ToGetTopicResourcesResponse4Web(ctx, tr, imageUrl, nil))
		}
	}
	return result, nil
}

func (uc *getTopicResourcesWebUseCase) GetOutputResources4Web(ctx context.Context, topicID, studentID string) ([]*response.GetTopicResourcesResponse4Web, error) {
	var topicResources []*model.TopicResource
	// get topic de kiem tra isAllPic
	topic, err := uc.topicRepository.GetByID(ctx, topicID)
	if err != nil {
		return nil, err
	}
	if topic == nil {
		return nil, fmt.Errorf("topic not found")
	}

	if topic.IsAllPic {
		resources, err := uc.topicResourceRepository.GetTopicResouresByStudentID(ctx, studentID)
		if err != nil {
			return nil, err
		}
		topicResources = resources
	} else {
		resources, err := uc.topicResourceRepository.GetTopicResouresByTopicAndStudent(ctx, topicID, studentID)
		if err != nil {
			return nil, err
		}
		topicResources = resources
	}

	result := make([]*response.GetTopicResourcesResponse4Web, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.IsOutput {
			var imageUrl string
			if tr.ImageKey != "" {
				if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
					imageUrl = *url
				}
			}
			result = append(result, mapper.ToGetTopicResourcesResponse4Web(ctx, tr, imageUrl, nil))
		}
	}
	return result, nil
}

func (uc *getTopicResourcesWebUseCase) GetTopicResourcesByStudent4Web(ctx context.Context, studentID string) ([]*response.GetTopicResourcesResponseByStudent4Web, error) {
	// get list topic trong topicResource có studentID
	result := make([]*response.GetTopicResourcesResponseByStudent4Web, 0)
	topicResources, err := uc.topicResourceRepository.GetTopicResouresByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	// lay danh sach topicID khac nhau trong topicResources
	topicIDs := make([]string, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		topicIDs = append(topicIDs, tr.TopicID)
	}
	topicIDs = helper.RemoveDuplicatesString(topicIDs)
	// loop qua topicIDs lay topic tu database
	for _, topicID := range topicIDs {
		topic, _ := uc.topicRepository.GetByID(ctx, topicID)
		if topic == nil {
			continue
		}
		uc.populateMediaUrlsForTopic(ctx, topic)
		appLang := helper.GetAppLanguage(ctx, 1)
		topicRes := mapper.ToTopicResponse(topic, appLang)
		// lay danh sach theo stdudentid, topicid
		topicResources := uc.filterByTopicID(topicResources, topicID)
		topicResourceResponses := make([]*response.TopicResourceResponse, 0, len(topicResources))
		for _, tr := range topicResources {
			// get resource image url
			var imageUrl string
			if tr.ImageKey != "" {
				if url, err := uc.s3Service.Get(ctx, tr.ImageKey, nil); err == nil && url != nil {
					imageUrl = *url
				}
			}
			topicResourceResponses = append(topicResourceResponses, &response.TopicResourceResponse{
				ID:        tr.ID.Hex(),
				FileName:  tr.FileName,
				ImageUrl:  imageUrl,
				CreatedAt: tr.CreatedAt,
				PicID:     tr.CreatedBy,
			})
		}
		result = append(result, &response.GetTopicResourcesResponseByStudent4Web{
			Topic:     topicRes,
			Resources: topicResourceResponses,
		})
	}

	return result, nil
}

func (uc *getTopicResourcesWebUseCase) populateMediaUrlsForTopic(ctx context.Context, topic *model.Topic) {
	for li := range topic.LanguageConfig {
		langCfg := &topic.LanguageConfig[li]

		// images
		for ii := range langCfg.Images {
			img := &langCfg.Images[ii]
			if img.ImageKey != "" {
				url, err := uc.s3Service.Get(ctx, img.ImageKey, nil)
				if err == nil && url != nil {
					img.UploadedUrl = *url
				} else {
					logger.WriteLogEx("get_topic_web_usecase", "populateMediaUrlsForTopic_images", fmt.Sprintf("error getting image url: %v", err))
				}
			}
		}

		// video
		if langCfg.Video.VideoKey != "" {
			url, err := uc.s3Service.Get(ctx, langCfg.Video.VideoKey, nil)
			if err == nil && url != nil {
				langCfg.Video.UploadedUrl = *url
			} else {
				logger.WriteLogEx("get_topic_web_usecase", "populateMediaUrlsForTopic_video", fmt.Sprintf("error getting image url: %v", err))
			}
		}

		// audio
		if langCfg.Audio.AudioKey != "" {
			url, err := uc.s3Service.Get(ctx, langCfg.Audio.AudioKey, nil)
			if err == nil && url != nil {
				langCfg.Audio.UploadedUrl = *url
			} else {
				logger.WriteLogEx("get_topic_web_usecase", "populateMediaUrlsForTopic_audio", fmt.Sprintf("error getting image url: %v", err))
			}
		}
	}
}

func (uc *getTopicResourcesWebUseCase) filterByTopicID(topicResources []*model.TopicResource, topicID string) []*model.TopicResource {
	result := make([]*model.TopicResource, 0, len(topicResources))
	for _, tr := range topicResources {
		if tr == nil {
			continue
		}
		if tr.TopicID == topicID {
			result = append(result, tr)
		}
	}
	return result
}
