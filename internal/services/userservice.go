package services

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

// Рега и логин через почту и пароль, не судите строго

func (us *userService) Register(username, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{Username: username, Email: email, PasswordHash: string(hash)}
	return us.CreateUser(user)
}

func (us *userService) Login(email, password string) (*models.User, error) {
	user, err := us.UserRepository.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}
