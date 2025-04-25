package services

import (
	postdto "CloneVK/internal/dto/posts"
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type postService struct {
	PostRepository repositories.IPostRepositories
}

func NewPostService(postRepository repositories.IPostRepositories) IPostService {
	return &postService{postRepository}
}

func (s *postService) CreatePost(dto *postdto.CreatePostDTO) (int, error) {
	post := models.Post{
		UserID:  dto.ID,
		Content: dto.Content,
		ImgURL:  dto.ImageURL,
	}

	err := s.PostRepository.CreatePost(&post)

	return post.ID, err
}

func (s *postService) FindPostByID(id int) (*models.Post, error) {
	return s.PostRepository.FindPostByID(id)
}

func (s *postService) GetAllPostsByUser(userId int) ([]models.Post, error) {
	return s.PostRepository.GetAllPostsByUser(userId)
}

func (s *postService) DeletePost(id int) error {
	return s.PostRepository.DeletePost(id)
}

func (s *postService) UpdatePost(newPost *postdto.UpdatePostDTO) error {
	post, err := s.FindPostByID(newPost.ID)
	if err != nil {
		return err
	}

	if newPost.Content != "" {
		post.Content = newPost.Content
	}

	if newPost.ImageURL != "" {
		post.ImgURL = newPost.ImageURL
	}

	return s.PostRepository.UpdatePost(post)
}
