package dto

type CreateArticleRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	TagList     []string `json:"tagList"`
}

type UpdateArticleRequest struct {
	Title       string `json:"title"`       // Optional
	Description string `json:"description"` // Optional
	Body        string `json:"body"`        // Optional
}
