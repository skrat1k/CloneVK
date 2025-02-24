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
	user := models.User{}
	c.ShouldBindJSON(&user)
	err := uh.UserService.CreateUser(&user)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, user.ID)
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

func (uh *UserHandler) FindAllUsers(c *gin.Context) {
	users, err := uh.UserService.FindAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, users)
}
