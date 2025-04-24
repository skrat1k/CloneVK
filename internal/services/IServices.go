package services

import (
	dto "CloneVK/internal/dto/posts"
	"CloneVK/internal/models"
	"time"
)

type IUserService interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)

	Register(username, email, password string) error

	Login(email, password string) (*models.User, error)

	SaveRefreshToken(userID int, token string, expiresAt time.Time) error
}

type IPostService interface {
	CreatePost(dto *dto.CreatePostDTO) (int, error)

	FindPostByID(id int) (*models.Post, error)

	GetAllPostsByUser(userId int) ([]models.Post, error)

	DeletePost(id int) error
}

type IFeedService interface {
	GetGlobalFeed(limit int, offset int) ([]models.Post, error)
}
