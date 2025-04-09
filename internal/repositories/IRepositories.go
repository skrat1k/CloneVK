package repositories

import "CloneVK/internal/models"

type IUserRepositories interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)

	FindUserByEmail(email string) (*models.User, error)
}
