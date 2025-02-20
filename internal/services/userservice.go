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
