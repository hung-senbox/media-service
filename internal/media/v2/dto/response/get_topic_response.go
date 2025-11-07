package response

type TopicResponse4Web struct {
	ID           string                 `json:"id"`
	IsPublished  bool                   `json:"is_published"`
	MainImageUrl string                 `json:"main_image_url"`
	MessageLangs []MessageLanguageEntry `json:"message_languages"`
}

type MessageLanguageEntry struct {
	LanguageID int              `json:"language_id"`
	Contents   LanguageContents `json:"contents"`
}

type LanguageContents struct {
	FileName    string              `json:"file_name"`
	Title       string              `json:"title"`
	Note        string              `json:"note"`
	Description string              `json:"description"`
	Audio       MediaContent        `json:"audio,omitempty"`
	Video       MediaContent        `json:"video,omitempty"`
	Images      map[string]ImgEntry `json:"images,omitempty"`
}

type MediaContent struct {
	UploadedURL string `json:"uploaded_url"`
	LinkURL     string `json:"link_url"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type ImgEntry struct {
	UploadedURL *string `json:"uploaded_url"`
	LinkURL     string  `json:"link_url"`
}

//// 4 App

type GetTopic4StudentResponse4App struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type GetTopicResponse4App struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type GetTopic4StudentResponse4Web struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type GetTopic4StudentResponse4Gw struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type TopicResponse4GW struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type TopicResponse2Assign4Web struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}
