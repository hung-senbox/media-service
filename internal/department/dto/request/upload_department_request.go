package request

type UploadDepartmentRequest struct {
	LocationID         string `json:"location_id"`
	ComponentID        string `json:"component_id"`
	RegionID           string `json:"region_id" binding:"required"`
	Name               string `json:"name" binding:"required"`
	Description        string `json:"description"`
	Note               string `json:"note"`
	Icon               string `json:"icon"`
	IsPublishedMessage bool   `json:"is_published_message"`
	Url                string `json:"url"`
	LanguageID         uint   `json:"language_id" binding:"required"`
	Message            string `json:"message"`
}
