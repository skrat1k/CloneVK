package dto

type CreatePostDTO struct {
	ID       int    `json:"id" validate:"required"`
	Content  string `json:"content" validate:"required"`
	ImageURL string `json:"imageURL,omitempty"`
}
