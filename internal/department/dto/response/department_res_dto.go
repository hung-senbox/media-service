package response

import (
	"media-service/internal/department/model"
	"media-service/internal/gateway/dto"
)

type DepartmentResponseDTO struct {
	ID             string             `json:"id"`
	LocationID     string             `json:"location_id"`
	OrganizationID string             `json:"organization_id"`
	RegionID       string             `json:"region_id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Note           string             `json:"note"`
	Icon           string             `json:"icon"`
	IconUrl        string             `json:"icon_url"`
	Leader         model.Leader       `json:"leader"`
	Staffs         []model.Staff      `json:"staffs"`
	Menus          []dto.MenuResponse `json:"menus"`
}

type GetDepartment4Web struct {
	ID                 string                        `json:"id"`
	LocationID         string                        `json:"location_id"`
	LocationName       string                        `json:"location_name"`
	OrganizationID     string                        `json:"organization_id"`
	RegionID           string                        `json:"region_id"`
	Icon               string                        `json:"icon"`
	IconUrl            string                        `json:"icon_url"`
	Url                string                        `json:"url"`
	IsPublishedMessage bool                          `json:"is_published_message"`
	Leader             LeaderResponseDTO             `json:"leader"`
	Staffs             []StaffResponseDTO            `json:"staffs"`
	HomeMenus          []dto.MenuResponse            `json:"home_menus"`
	OrganizationMenus  []dto.MenuResponse            `json:"organization_menus"`
	MessageLanguages   []dto.MessageLanguageResponse `json:"message_languages"`
}

type LeaderResponseDTO struct {
	OwnerID   string `json:"owner_id"`
	OwnerRole string `json:"owner_role"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type StaffResponseDTO struct {
	Index     int    `json:"index"`
	OwnerID   string `json:"owner_id"`
	OwnerRole string `json:"owner_role"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type DepartmentGroupResponse struct {
	RegionID    string               `json:"region_id"`
	RegionName  string               `json:"region_name"`
	Departments []*GetDepartment4Web `json:"departments"`
}

type GetDepartment4App struct {
	ID             string `json:"id"`
	LocationID     string `json:"location_id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Note           string `json:"note"`
	Icon           string `json:"icon"`
	Message        string `json:"message"`
}

type GetDepartment4Gateway struct {
	ID   string `json:"id"`
	Icon string `json:"icon"`
	Name string `json:"name"`
}

type GetDepartments4Web struct {
	RegionID    string               `json:"region_id"`
	RegionName  string               `json:"region_name"`
	Departments []*GetDepartment4Web `json:"departments"`
}
