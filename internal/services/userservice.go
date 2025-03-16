package services

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
)

type userService struct {
	UserRepository repositories.IUserRepositories
}

func NewUserService(userRepository repositories.IUserRepositories) IUserService {
	return &userService{userRepository}
}

func (us *userService) CreateUser(user *models.User) error {
	return us.UserRepository.CreateUser(user)
}

func (us *userService) FindUserByID(id int) (*models.User, error) {
	return us.UserRepository.FindUserByID(id)
}

func (us *userService) FindAllUsers() (*[]models.User, error) {
	return us.UserRepository.FindAllUsers()
}
