package models

type SEO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	AuthorName  string `json:"author_name"`
	ImageURL    string `json:"image_url"`
	PageURL     string `json:"page_url"`
}
