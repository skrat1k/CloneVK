package handlers

import (
	"CloneVK/internal/models"
	"CloneVK/internal/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *services.UserService
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	err := uh.UserService.CreateUser(&models.User{Username: "ttt", PasswordHash: "222", Email: "ttt@e.com", AvatarURL: "555"})
	if err != nil {
		log.Fatal(err)
	}
}

func (uh *UserHandler) FindUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Fatal(err)
	}
	user, err := uh.UserService.FindUserByID(id)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, user)
}
