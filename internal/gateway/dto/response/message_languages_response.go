package response

type MessageLanguageResponse struct {
	LangID   uint              `json:"language_id"`
	Contents map[string]string `json:"contents"`
}
