package repositories

import (
	"CloneVK/internal/models"
	"time"
)

type IUserRepositories interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)

	FindUserByEmail(email string) (*models.User, error)

	SaveRefreshToken(userID int, token string, expiresAt time.Time) error

	//GetRefreshToken(userID int) (string, error)

	// DeleteRefreshTokensByUserID(userID int) error
}

type IPostRepositories interface {
	FindPostByID(id int) (*models.Post, error)

	GetAllPostsByUser(userId int) ([]models.Post, error)

	CreatePost(post *models.Post) error
}
