package services

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type UserService struct {
	UserRepository *repositories.UserRepository
}

func (us *UserService) CreateUser(user *models.User) error {
	return us.UserRepository.CreateUser(user)
}

func (us *UserService) FindUserByID(id int) (*models.User, error) {
	return us.UserRepository.FindUserByID(id)
}
