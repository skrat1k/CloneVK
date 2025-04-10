package handlers

import (
	"CloneVK/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (uh *userHandler) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	// users, err := uh.UserService.FindAllUsers()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// c.JSON(http.StatusOK, users)
	w.Write([]byte("All users"))
}

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
