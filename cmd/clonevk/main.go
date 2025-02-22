package main

import (
	"CloneVK/internal/handlers"
	"CloneVK/internal/repositories"
	"CloneVK/internal/services"
	"CloneVK/internal/storage"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

// Потом сделать подгрузку из файла окружения
const (
	UsernameDB = "postgres"
	PasswordDB = "admin"
	HostDB     = "localhost"
	PortDB     = "5432"
	NameDB     = "clonevk"
)

func main() {
	conn, err := storage.CreatePostgresConnection(storage.ConnectionInfo{
		Username: UsernameDB,
		Password: PasswordDB,
		Host:     HostDB,
		Port:     PortDB,
		DBName:   NameDB,
	})

	router := gin.Default()

	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close(context.Background())
	ur := &repositories.UserRepository{DB: conn}

	us := &services.UserService{UserRepository: ur}

	uh := handlers.UserHandler{UserService: us}
	// err = us.CreateUser(&models.User{Username: "g", PasswordHash: "222", Email: "g@e.com", AvatarURL: "555"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	router.POST("/user", uh.CreateUser)
	router.GET("/user/:id", uh.FindUserByID)

	err = router.Run("localhost:8081")
	if err != nil {
		log.Fatal(err)
	}
}
