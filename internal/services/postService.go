package services

import (
	dto "CloneVK/internal/dto/posts"
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type postService struct {
	PostRepository repositories.IPostRepositories
}

func NewPostService(postRepository repositories.IPostRepositories) IPostService {
	return &postService{postRepository}
}

func (ps *postService) CreatePost(dto *dto.CreatePostDTO) (int, error) {
	post := models.Post{
		UserID:  dto.ID,
		Content: dto.Content,
		ImgURL:  dto.ImageURL,
	}

	err := ps.PostRepository.CreatePost(&post)

	return post.ID, err
}

func (ps *postService) FindPostByID(id int) (*models.Post, error) {
	return ps.PostRepository.FindPostByID(id)
}

func (ps *postService) GetAllPostsByUser(userId int) ([]models.Post, error) {
	return ps.PostRepository.GetAllPostsByUser(userId)
}

func (ps *postService) DeletePost(id int) error {
	return ps.PostRepository.DeletePost(id)
}
