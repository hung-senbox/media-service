package response

type VocabularyResponse4Web struct {
	ID           string                           `json:"id"`
	IsPublished  bool                             `json:"is_published"`
	MainImageUrl string                           `json:"main_image_url"`
	MessageLangs []VocabularyMessageLanguageEntry `json:"message_languages"`
}

type VocabularyMessageLanguageEntry struct {
	LanguageID int                        `json:"language_id"`
	Contents   VocabularyLanguageContents `json:"contents"`
}

type VocabularyLanguageContents struct {
	FileName    string                        `json:"file_name"`
	Title       string                        `json:"title"`
	Note        string                        `json:"note"`
	Description string                        `json:"description"`
	Audio       VocabularyMediaContent        `json:"audio,omitempty"`
	Video       VocabularyMediaContent        `json:"video,omitempty"`
	Images      map[string]VocabularyImgEntry `json:"images,omitempty"`
}

type VocabularyMediaContent struct {
	UploadedURL string `json:"uploaded_url"`
	LinkURL     string `json:"link_url"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}

type VocabularyImgEntry struct {
	UploadedURL *string `json:"uploaded_url"`
	LinkURL     string  `json:"link_url"`
}

//// 4 App

type GetVocabulary4StudentResponse4App struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type GetVocabularyResponse4App struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type GetVocabulary4StudentResponse4Web struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type GetVocabulary4StudentResponse4Gw struct {
	ID           string `json:"id"`
	IsPublished  bool   `json:"is_published"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type VocabularyResponse4GW struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type VocabularyResponse2Assign4Web struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}

type VocabularyResponse4Gw struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	MainImageUrl string `json:"main_image_url"`
}
