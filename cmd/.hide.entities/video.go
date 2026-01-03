package entities

type Video struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	URL         string `json:"url" binding:"required,http_url"`
}
