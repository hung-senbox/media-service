package dto

type TeacherResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Avatar         Avatar `json:"avatar"`
}
