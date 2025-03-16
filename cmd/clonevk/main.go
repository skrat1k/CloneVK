package main

import (
	"CloneVK/internal/handlers"
	"CloneVK/internal/repositories"
	"CloneVK/internal/services"
	"CloneVK/internal/storage"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// TODO: add logger, add error handler, add config and env variables
// TODO: Реализовать всю структуру для постов и лайков
// TODO: Сделать ленту
// TODO: Добавить реддис
// TODO: Сделать чат
// TODO: Добавить таймаута, а значит и прокидывать контекты во все репозитории
// TODO: Написать тесты
// Потом сделать подгрузку из файла окружения
const (
	UsernameDB = "postgres"
	PasswordDB = "admin"
	HostDB     = "localhost"
	PortDB     = "5432"
	NameDB     = "clonevk"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO", 0)

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

	infoLog.Printf("Success connect to database: %s", NameDB)

	router := chi.NewRouter()

	ur := repositories.NewUserRepositories(conn)

	us := services.NewUserService(ur)

	uh := handlers.NewUserHandler(us)

	uh.Register(router)

	infoLog.Printf("Server succesfully started at port: %s", "8082")

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		log.Fatal(err)
	}
}
