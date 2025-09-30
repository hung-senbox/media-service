package request

type CreateRegionRequest struct {
	Name string `json:"name" binding:"required"`
}
