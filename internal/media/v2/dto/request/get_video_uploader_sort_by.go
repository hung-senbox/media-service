package request

type GetVideoUploaderSortBy struct {
	Field string `form:"field"`
	Order string `form:"order"`
}

const (
	GetVideoUploaderSortByFieldTitle      = "title"
	GetVideoUploaderSortByFieldLanguageID = "language_id"
	GetVideoUploaderSortByFieldCreatedAt  = "created_at"
	GetVideoUploaderSortByFieldUpdatedAt  = "updated_at"
)

const (
	GetVideoUploaderSortByOrderAsc  = "asc"
	GetVideoUploaderSortByOrderDesc = "desc"
)
