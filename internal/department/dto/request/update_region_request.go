package request

type UpdateRegionRequest struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}
