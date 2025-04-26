package handlers

import (
	"CloneVK/internal/services"
	logger "CloneVK/pkg/Logger"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	_ "CloneVK/docs"

	"github.com/go-chi/chi/v5"
)

const (
	//createUserURL = "/user"
	getUserURL    = "/users/{id}"
	getAllUserURL = "/users"
	registerURL   = "/auth/register"
	loginURL      = "/auth/login"
)

type userHandler struct {
	UserService services.IUserService
	JWTService  services.JWTService
	Log         *slog.Logger
}

func NewUserHandler(userService services.IUserService, jwtService services.JWTService, log *slog.Logger) IHandler {
	lg := logger.WithHandler(log, "UserHandler")
	return &userHandler{userService, jwtService, lg}
}

func (uh *userHandler) Register(router *chi.Mux) {
	//router.Post(createUserURL, uh.CreateUser)
	//uh.Log.Info("Successfully created http route", slog.String("route", createUserURL))
	router.Get(getUserURL, uh.FindUserByID)
	uh.Log.Info("Successfully created http route", slog.String("route", getUserURL))
	router.Get(getAllUserURL, uh.FindAllUsers)
	uh.Log.Info("Successfully created http route", slog.String("route", getAllUserURL))

	router.Post(registerURL, uh.RegisterUser)
	uh.Log.Info("Successfully created http route", slog.String("route", registerURL))
	router.Post(loginURL, uh.LoginUser)
	uh.Log.Info("Successfully created http route", slog.String("route", loginURL))
}

// @Summary Получить пользователя по ID
// @Description Получает информацию о пользователе
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User Not Found"
// @Failure 500 {string} string "Failed to find user"
// @Router /users/{id} [get]
func (uh *userHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {
	log := logger.WithMethod(uh.Log, "FindUserByID")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Error("Failed convert to int", slog.String("error", err.Error()))
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	log.Debug("Getting id", slog.Int("userId", id))

	user, err := uh.UserService.FindUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			log.Error("User not found", slog.Int("userId", id))
			return
		}

		log.Error("Failed to find user by id", slog.Int("id", id), slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to find user by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

// @Summary Получить всех пользователей
// @Description Получает информацию о всех пользователях
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Failed to get users"
// @Router /users [get]
func (uh *userHandler) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	log := logger.WithMethod(uh.Log, "FindAllUsers")
	users, err := uh.UserService.FindAllUsers()
	if err != nil {
		log.Error("Failed to get users", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to get users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// @Summary Регистрация пользователя
// @Description Создать пользователя (потом у этого метода будет другой функционал, но пока так)
// @Tags users
// @Accept json
// @Produce json
// @Param userInfo body dto.RegisterDTO true "Пользователь"}
// @Success 200
// @Router /auth/register [post]
func (uh *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log := logger.WithMethod(uh.Log, "RegisterUser")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Invalid JSON payload", slog.String("error", err.Error()))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Debug("Attempting to register user", slog.String("email", req.Email))

	err := uh.UserService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		log.Error("Failed to register user", slog.String("error", err.Error()), slog.String("email", req.Email))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := uh.UserService.Login(req.Email, req.Password)
	if err != nil {
		log.Error("Auto-login failed after registration", slog.String("error", err.Error()))
		http.Error(w, "Auto-login failed", http.StatusInternalServerError)
		return
	}

	tokenPair, err := uh.JWTService.GenerateTokens(user.ID)
	if err != nil {
		log.Error("Token generation failed", slog.Int("userID", user.ID), slog.String("error", err.Error()))
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	err = uh.UserService.SaveRefreshToken(user.ID, tokenPair.RefreshToken, tokenPair.RefreshTokenExpiresAt)

	if err != nil {
		log.Error("Failed to save refresh token", slog.String("error", err.Error()))
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	log.Info("User registered and logged in", slog.Int("userID", user.ID), slog.String("email", user.Email))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

// @Summary Логин пользователя
// @Description Логин (потом у этого метода будет другой функционал, но пока так)
// @Tags users
// @Accept json
// @Produce json
// @Param userInfo body dto.LoginDTO true "Пользователь"}
// @Success 200 {object} map[string]string "Токен"
// @Router /auth/login [post]
func (uh *userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	log := logger.WithMethod(uh.Log, "LoginUser")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Invalid JSON payload", slog.String("error", err.Error()))
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Debug("Attempting login", slog.String("email", req.Email))

	user, err := uh.UserService.Login(req.Email, req.Password)
	if err != nil {
		log.Warn("Unauthorized login attempt", slog.String("email", req.Email))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenPair, err := uh.JWTService.GenerateTokens(user.ID)
	if err != nil {
		log.Error("Token generation failed", slog.Int("userID", user.ID), slog.String("error", err.Error()))
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	err = uh.UserService.SaveRefreshToken(user.ID, tokenPair.RefreshToken, tokenPair.RefreshTokenExpiresAt)
	if err != nil {
		log.Error("Failed to save refresh token", slog.String("error", err.Error()))
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	log.Info("User successfully logged in", slog.Int("userID", user.ID), slog.String("email", user.Email))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}
