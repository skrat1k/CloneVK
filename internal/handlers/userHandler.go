package handlers

import (
	"CloneVK/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "CloneVK/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	createUserURL = "/user"
	getUserURL    = "/user/{id}"
	getAllUserURL = "/users"
	registerURL   = "/auth/register"
	loginURL      = "/auth/login"
)

type userHandler struct {
	UserService services.IUserService
	JWTService  services.JWTService
}

func NewUserHandler(userService services.IUserService, jwtService services.JWTService) IHandler {
	return &userHandler{userService, jwtService}
}

func (uh *userHandler) Register(router *chi.Mux) {
	router.Post(createUserURL, uh.CreateUser)
	router.Get(getUserURL, uh.FindUserByID)
	router.Get(getAllUserURL, uh.FindAllUsers)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Post(registerURL, uh.RegisterUser)
	router.Post(loginURL, uh.LoginUser)
}

func (uh *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// user := models.User{}
	// c.ShouldBindJSON(&user)
	// err := uh.UserService.CreateUser(&user)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.JSON(http.StatusOK, user.ID)
}

// @Summary Получить пользователя по ID
// @Description Получает информацию о пользователе
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /user/{id} [get]
func (uh *userHandler) FindUserByID(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Fatal(err)
	}

	user, err := uh.UserService.FindUserByID(id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(user)

}

// @Summary Получить всех пользователей
// @Description Получает информацию о всех пользователях
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Router /users [get]
func (uh *userHandler) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := uh.UserService.FindAllUsers()
	if err != nil {
		log.Fatal(err)
	}

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
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := uh.UserService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := uh.UserService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := uh.JWTService.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
