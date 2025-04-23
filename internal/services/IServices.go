package services

import (
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
