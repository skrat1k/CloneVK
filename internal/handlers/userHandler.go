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
)

type userHandler struct {
	UserService services.IUserService
}

func NewUserHandler(userService services.IUserService) IHandler {
	return &userHandler{userService}
}

func (uh *userHandler) Register(router *chi.Mux) {
	router.Post(createUserURL, uh.CreateUser)
	router.Get(getUserURL, uh.FindUserByID)
	router.Get(getAllUserURL, uh.FindAllUsers)
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
