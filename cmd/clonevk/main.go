package main

import (
	"CloneVK/internal/config"
	"CloneVK/internal/handlers"
	"CloneVK/internal/repositories"
	"CloneVK/internal/services"
	"CloneVK/internal/storage"
	"CloneVK/internal/storage/migrations"
	logger "CloneVK/pkg/Logger"
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// TODO: add error handler
// TODO: Реализовать всю структуру для постов и лайков
// TODO: Сделать ленту
// TODO: Добавить реддис
// TODO: Сделать чат
// TODO: Добавить таймаута, а значит и прокидывать контекты во все репозитории
// TODO: Написать тесты

// @title CloneVK
// @version dev
// @description Социальная сеть на golang
// @host localhost:8083
// @BasePath /
func main() {

	cfg, err := config.MustLoad()
	if err != nil {
		panic(err)
	}

	log := logger.GetLogger(cfg.Env)

	psqlConnectionUrl := storage.MakeURL(storage.ConnectionInfo{
		Username: cfg.UsernameDB,
		Password: cfg.PasswordDB,
		Host:     cfg.HostDB,
		Port:     cfg.PortDB,
		DBName:   cfg.NameDB,
		SSLMode:  cfg.SSLModeDB,
	})

	if err := migrations.RunMigrations(psqlConnectionUrl); err != nil {
		if strings.Contains(err.Error(), "no change") {
			log.Debug("No new migrations found, continuing application startup...")
		} else {
			log.Error("Migration error", slog.String("error", err.Error()))
			panic(err)
		}
	}

	conn, err := storage.CreatePostgresConnection(psqlConnectionUrl)

	if err != nil {
		log.Error("Connection error", slog.String("error", err.Error()))
		os.Exit(1) // Мб тут корректнее прописать панику, а не выход, о пока хз не разобрался, работает и слава богу
	}

	defer conn.Close(context.Background())

	log.Info("Success connect to database")

	router := chi.NewRouter()
	// Вот это никому не трогать, я пока не совсем понимаю как оно работает, но без этих строк не работает фронт
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // или "*", если тестируешь локально
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Max value = 600
	}))

	jwtService := services.NewJWTService(log)

	userRepo := repositories.NewUserRepositories(conn)
	userService := services.NewUserService(userRepo, log)
	userHandler := handlers.NewUserHandler(userService, jwtService, log)
	userHandler.Register(router)

	postRepo := repositories.NewPostRepository(conn)
	postService := services.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService, log)
	postHandler.Register(router)

	feedRepo := repositories.NewFeedRepository(conn)
	feedService := services.NewFeedService(feedRepo)
	feedHandler := handlers.NewFeedHandler(feedService)
	feedHandler.Register(router)

	followRepo := repositories.NewFollowRepository(conn)
	followService := services.NewFollowService(followRepo)
	followHandler := handlers.NewFollowHandler(followService)
	followHandler.Register(router)

	router.Get("/swagger/*", httpSwagger.WrapHandler)
	log.Info("Swagger init")

	log.Info("Server succesfully started at port")

	serverPort := cfg.ServerPort

	err = http.ListenAndServe(serverPort, router)
	if err != nil {
		log.Error("Error", slog.String("error", err.Error()))
	}
}
