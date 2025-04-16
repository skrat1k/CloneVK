package services

import (
	dto "CloneVK/internal/dto/posts"
	"CloneVK/internal/models"
)

type IUserService interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)

	Register(username, email, password string) error

	Login(email, password string) (*models.User, error)
}

type IPostService interface {
	CreatePost(dto *dto.CreatePostDTO) (int, error)

	FindPostByID(id int) (*models.Post, error)

	GetAllPostsByUser(userId int) ([]models.Post, error)
}
