package response

type UploadPDFResponse struct {
	PDFName        string `json:"pdf_name"`
	Key            string `json:"key"`
	OrganizationID string `json:"organization_id"`
	Url            string `json:"url"`
	Extension      string `json:"extension"`
}
