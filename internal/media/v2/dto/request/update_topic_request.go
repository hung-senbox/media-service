package request

type UpdateTopicRequest struct {
	TopicID     string `json:"topic_id"`
	LanguageID  uint   `json:"language_id" binding:"required"`
	IsPublished bool   `json:"is_published" binding:"required"`
	FileName    string `json:"file_name"`
	Title       string `json:"title"`
	Node        string `json:"note"`
	Description string `json:"description"`
}
