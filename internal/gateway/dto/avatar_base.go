package dto

type Avatar struct {
	ImageID  uint64 `json:"image_id"`
	ImageKey string `json:"image_key"`
	ImageUrl string `json:"image_url"`
	Index    int    `json:"index"`
	IsMain   bool   `json:"is_main"`
}
