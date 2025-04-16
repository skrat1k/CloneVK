package services

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
	logger "CloneVK/pkg/Logger"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	UserRepository repositories.IUserRepositories
	Log            *slog.Logger
}

func NewUserService(userRepository repositories.IUserRepositories, log *slog.Logger) IUserService {
	lg := logger.WithService(log, "UserService")
	return &userService{UserRepository: userRepository, Log: lg}
}

func (us *userService) CreateUser(user *models.User) error {
	log := logger.WithMethod(us.Log, "CreateUser")
	log.Debug("Attempting to create user", slog.String("email", user.Email), slog.String("username", user.Username))

	err := us.UserRepository.CreateUser(user)
	if err != nil {
		log.Error("Failed to create user", slog.String("error", err.Error()))
		return err
	}

	log.Info("User successfully created", slog.Int("userID", user.ID))
	return nil
}

func (us *userService) FindUserByID(id int) (*models.User, error) {
	log := logger.WithMethod(us.Log, "FindUserByID")
	log.Debug("Searching user by ID", slog.Int("userID", id))

	user, err := us.UserRepository.FindUserByID(id)
	if err != nil {
		log.Error("User not found or error occurred", slog.String("error", err.Error()), slog.Int("userID", id))
		return nil, err
	}

	log.Info("User found", slog.Int("userID", user.ID))
	return user, nil
}

func (us *userService) FindAllUsers() (*[]models.User, error) {
	log := logger.WithMethod(us.Log, "FindAllUsers")
	log.Debug("Fetching all users")

	users, err := us.UserRepository.FindAllUsers()
	if err != nil {
		log.Error("Failed to fetch users", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("Users fetched successfully", slog.Int("count", len(*users)))
	return users, nil
}

func (us *userService) Register(username, email, password string) error {
	log := logger.WithMethod(us.Log, "Register")
	log.Debug("Attempting to register user", slog.String("email", email), slog.String("username", username))

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Password hashing failed", slog.String("error", err.Error()))
		return err
	}

	user := &models.User{Username: username, Email: email, PasswordHash: string(hash)}

	err = us.CreateUser(user)
	if err != nil {
		log.Error("User registration failed", slog.String("error", err.Error()))
		return err
	}

	log.Info("User registered successfully", slog.String("email", email))
	return nil
}

func (us *userService) Login(email, password string) (*models.User, error) {
	log := logger.WithMethod(us.Log, "Login")
	log.Debug("Attempting login", slog.String("email", email))

	user, err := us.UserRepository.FindUserByEmail(email)
	if err != nil {
		log.Warn("Email not found", slog.String("email", email))
		return nil, fmt.Errorf("email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Warn("Invalid credentials", slog.String("email", email))
		return nil, fmt.Errorf("invalid credentials")
	}

	log.Info("User login successful", slog.Int("userID", user.ID))
	return user, nil
}
