package response

type GetUploadProgressResponse struct {
	Progress     int            `json:"progress"`
	UploadErrors map[string]any `json:"upload_errors"`
}
