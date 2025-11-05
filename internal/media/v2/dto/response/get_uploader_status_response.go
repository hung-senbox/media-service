package response

type GetUploaderStatusResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	UpdatedAt string `json:"updated_at"`
}
