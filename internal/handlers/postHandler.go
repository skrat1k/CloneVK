package handlers

import (
	dto "CloneVK/internal/dto/posts"
	"CloneVK/internal/services"
	logger "CloneVK/pkg/Logger"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

const (
	getPostURL       = "/posts/{id}"
	getAllPostURL    = "/posts"
	createPostURL    = "/posts"
	getPostsFromUser = "/posts/user/{id}"
)

type postHandler struct {
	PostService services.IPostService
	Log         *slog.Logger
}

func NewPostHandler(postService services.IPostService, log *slog.Logger) IHandler {
	lg := logger.WithHandler(log, "PostHandler")
	return &postHandler{postService, lg}
}

func (ph *postHandler) Register(router *chi.Mux) {
	router.Post(createPostURL, ph.CreatePost)
}

func (ph *postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var dto dto.CreatePostDTO

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(dto)
	if err != nil {
		http.Error(w, "Validate error", http.StatusBadRequest)
		return
	}

	id, err := ph.PostService.CreatePost(&dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create person: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
