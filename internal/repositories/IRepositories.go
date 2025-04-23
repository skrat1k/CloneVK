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

	//DeleteRefreshToken(userID int) error
}
