package mapper

import (
	"media-service/internal/department/dto/response"
	"media-service/internal/department/model"
	"media-service/internal/gateway/dto"
)

func MapDepartmentToResponse(dept *model.Department, menus []dto.MenuResponse, iconUrl string) *response.DepartmentResponseDTO {
	staffs := dept.Staffs
	if staffs == nil {
		staffs = []model.Staff{}
	}
	if menus == nil {
		menus = []dto.MenuResponse{}
	}
	return &response.DepartmentResponseDTO{
		ID:             dept.ID.Hex(),
		LocationID:     dept.LocationID,
		OrganizationID: dept.OrganizationID,
		RegionID:       dept.RegionID,
		Name:           dept.Name,
		Description:    dept.Description,
		Note:           dept.Note,
		Icon:           dept.Icon,
		IconUrl:        iconUrl,
		Leader:         dept.Leader,
		Staffs:         staffs,
		Menus:          menus,
	}
}

func MapDepartmentsToResponses(depts []*model.Department) []*response.DepartmentResponseDTO {
	responses := make([]*response.DepartmentResponseDTO, len(depts))
	for i, dept := range depts {
		responses[i] = MapDepartmentToResponse(dept, nil, "")
	}
	return responses
}

func MapDepartmentsToGroupedResponses4Web(
	depts []*model.Department,
	homeMenusMap map[string][]dto.MenuResponse,
	iconUrls map[string]string,
	leaders map[string]response.LeaderResponseDTO,
	staffsMap map[string][]response.StaffResponseDTO,
	organizationMenusMap map[string][]dto.MenuResponse,
	messageLanguagesMaps map[string][]dto.MessageLanguageResponse,
	locationNamesMap map[string]string,
) []*response.DepartmentGroupResponse {
	groupMap := make(map[string][]*response.GetDepartment4Web)

	for _, dept := range depts {
		homeMenus := homeMenusMap[dept.ID.Hex()]
		iconUrl := iconUrls[dept.ID.Hex()]
		leader := leaders[dept.ID.Hex()]
		staffs := staffsMap[dept.ID.Hex()]
		organizationMenus := organizationMenusMap[dept.ID.Hex()]
		messageLanguages := messageLanguagesMaps[dept.ID.Hex()]
		locationName := locationNamesMap[dept.LocationID]

		resp := MapDepartmentToResponse4Web(dept, homeMenus, iconUrl, leader, staffs, organizationMenus, messageLanguages, locationName)
		groupMap[dept.RegionID] = append(groupMap[dept.RegionID], resp)
	}

	// convert map -> slice
	var result []*response.DepartmentGroupResponse
	for regjonIdx, depts := range groupMap {
		result = append(result, &response.DepartmentGroupResponse{
			RegionID:    regjonIdx,
			Departments: depts,
		})
	}

	return result
}

func MapDepartmentToResponse4App(dept *model.Department, msg dto.MessageLanguageResponse) *response.GetDepartment4App {
	name := dept.Name
	description := dept.Description
	note := dept.Note
	mess := ""

	// Nếu gateway trả về có giá trị thì override
	if val, ok := msg.Contents["name"]; ok && val != "" {
		name = val
	}
	if val, ok := msg.Contents["description"]; ok && val != "" {
		description = val
	}
	if val, ok := msg.Contents["note"]; ok && val != "" {
		note = val
	}
	if val, ok := msg.Contents["message"]; ok && val != "" {
		mess = val
	}

	return &response.GetDepartment4App{
		ID:             dept.ID.Hex(),
		LocationID:     dept.LocationID,
		OrganizationID: dept.OrganizationID,
		Name:           name,
		Description:    description,
		Note:           note,
		Icon:           dept.Icon,
		Message:        mess,
	}
}

// func MapDepartmentsToResponses4App(depts []*model.Department) []*response.GetDepartment4App {
// 	responses := make([]*response.GetDepartment4App, len(depts))
// 	for i, dept := range depts {
// 		responses[i] = MapDepartmentToResponse4App(dept)
// 	}
// 	return responses
// }

func MapDepartmentToResponse4Gateway(dept *model.Department, msg dto.MessageLanguageResponse) *response.GetDepartment4Gateway {
	name := dept.Name

	if val, ok := msg.Contents["name"]; ok && val != "" {
		name = val
	}

	return &response.GetDepartment4Gateway{
		ID:   dept.ID.Hex(),
		Name: name,
		Icon: dept.Icon,
	}
}

func MapDepartmentsToResponses4Gateway(depts []*model.Department) []*response.GetDepartment4Gateway {
	responses := make([]*response.GetDepartment4Gateway, len(depts))
	for i, dept := range depts {
		responses[i] = MapDepartmentToResponse4Gateway(dept, dto.MessageLanguageResponse{})
	}
	return responses
}

func MapDepartmentToResponse4Web(
	dept *model.Department,
	homeMenus []dto.MenuResponse,
	iconUrl string,
	leader response.LeaderResponseDTO,
	staffs []response.StaffResponseDTO,
	organizationMenus []dto.MenuResponse,
	messageLanguages []dto.MessageLanguageResponse,
	locationName string,
) *response.GetDepartment4Web {

	if homeMenus == nil {
		homeMenus = []dto.MenuResponse{}
	}
	if organizationMenus == nil {
		organizationMenus = []dto.MenuResponse{}
	}
	if staffs == nil {
		staffs = []response.StaffResponseDTO{}
	}
	if messageLanguages == nil {
		messageLanguages = []dto.MessageLanguageResponse{}
	}
	return &response.GetDepartment4Web{
		ID:                 dept.ID.Hex(),
		LocationID:         dept.LocationID,
		OrganizationID:     dept.OrganizationID,
		RegionID:           dept.RegionID,
		Icon:               dept.Icon,
		IconUrl:            iconUrl,
		Url:                dept.Url,
		Leader:             leader,
		Staffs:             staffs,
		HomeMenus:          homeMenus,
		OrganizationMenus:  organizationMenus,
		MessageLanguages:   messageLanguages,
		IsPublishedMessage: dept.IsPublishedMessage,
		LocationName:       locationName,
	}
}

func MapDepartmentsToResponses4Web(
	depts []*model.Department,
	homeMenusMap map[string][]dto.MenuResponse,
	iconUrls map[string]string,
	leaders map[string]response.LeaderResponseDTO,
	staffsMap map[string][]response.StaffResponseDTO,
	organizationMenusMap map[string][]dto.MenuResponse,
	messageLanguagesMap map[string][]dto.MessageLanguageResponse,
	locationNamesMap map[string]string,
) []*response.GetDepartment4Web {
	result := make([]*response.GetDepartment4Web, 0, len(depts))
	for _, dept := range depts {
		homeMenus := homeMenusMap[dept.ID.Hex()]
		iconUrl := iconUrls[dept.ID.Hex()]
		leader := leaders[dept.ID.Hex()]
		staffs := staffsMap[dept.ID.Hex()]
		orgMenus := organizationMenusMap[dept.ID.Hex()]
		messageLanguages := messageLanguagesMap[dept.ID.Hex()]
		locationName := locationNamesMap[dept.LocationID]

		resp := MapDepartmentToResponse4Web(dept, homeMenus, iconUrl, leader, staffs, orgMenus, messageLanguages, locationName)
		result = append(result, resp)
	}
	return result
}
