package request

import "github.com/gofrs/uuid"

type UploadSectionMenuDepartmentRequest DepartmentSectionMenuItem

type DepartmentSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	DepartmentID       string                       `json:"department_id" binding:"required"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}

type CreateMenuComponentRequest struct {
	ID        *uuid.UUID `json:"id"`
	SectionId string     `json:"section_id"`
	Name      string     `json:"name" binding:"required"`
	Type      string     `json:"type" binding:"required"`
	Key       string     `json:"key" binding:"required" default:""`
	Value     string     `json:"value" binding:"required"`
	Order     int        `json:"order" binding:"required"`
	IsShow    bool       `json:"is_show"`
}
