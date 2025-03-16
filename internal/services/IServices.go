package services

import "CloneVK/internal/models"

type IUserService interface {
	CreateUser(user *models.User) error

	FindUserByID(id int) (*models.User, error)

	FindAllUsers() (*[]models.User, error)
}
