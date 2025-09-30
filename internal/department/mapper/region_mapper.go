package mapper

import (
	"media-service/internal/department/dto/response"
	"media-service/internal/department/model"
)

func MapRegionToResponse(region *model.Region) *response.RegionResponseDTO {
	return &response.RegionResponseDTO{
		ID:   region.ID.Hex(),
		Name: region.Name,
	}
}

func MapRegionsToResponse(regions []*model.Region) []*response.RegionResponseDTO {
	res := []*response.RegionResponseDTO{}
	for _, region := range regions {
		res = append(res, MapRegionToResponse(region))
	}
	return res
}
