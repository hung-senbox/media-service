package request

type GetFileUrlRequest struct {
	Key  string `json:"key" binding:"required"`
	Mode string `json:"mode" binding:"required"`
}
