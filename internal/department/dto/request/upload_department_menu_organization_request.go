package request

type UploadDepartmentMenuOrganizationRequest DepartmentSectionMenuOrganizationItem

type DepartmentSectionMenuOrganizationItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	DepartmentID       string                       `json:"department_id"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
