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
	"github.com/go-chi/cors"
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

// @title CloneVK
// @version dev
// @description Социальная сеть на golang
// @host localhost:8082
// @BasePath /
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

	jwtService := services.NewJWTService()

	ur := repositories.NewUserRepositories(conn)

	us := services.NewUserService(ur)

	uh := handlers.NewUserHandler(us, jwtService)

	uh.Register(router)
	// кусок кода сгенеренный гпт, потом надо бы разобраться с ним
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // или "*", если тестируешь локально
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Max value = 600
	}))
	//
	infoLog.Printf("Server succesfully started at port: %s", "8082")

	err = http.ListenAndServe(":8082", router)
	if err != nil {
		log.Fatal(err)
	}
}
