package response

type GetUploadProgressResponse struct {
	Progress     int            `json:"progress"`
	FileName     string         `json:"file_name"`
	UploadErrors map[string]any `json:"upload_errors"`
}
