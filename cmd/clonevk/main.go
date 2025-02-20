package main

import (
	"CloneVK/internal/models"
	"CloneVK/internal/repositories"
	"CloneVK/internal/storage"
	"context"
	"log"
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
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close(context.Background())

	ur := repositories.UserRepository{DB: conn}
	err = ur.CreateUser(&models.User{Username: "g", PasswordHash: "222", Email: "g@e.com", AvatarURL: "555"})
	if err != nil {
		log.Fatal(err)
	}
}
