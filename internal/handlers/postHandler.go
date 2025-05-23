package handlers

import (
	dto "CloneVK/internal/dto/posts"
	"CloneVK/internal/services"
	logger "CloneVK/pkg/Logger"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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
	router.Get(getPostURL, ph.FindPostByID)
	router.Get(getPostsFromUser, ph.GetAllPostsByUser)
}

// @Summary Создание поста
// @Description Создаёт пост и добавялет его в базу данных
// @Tags posts
// @Accept json
// @Produce json
// @Param postInfo body dto.CreatePostDTO true "Пост"
// @Success 200
// @Router /posts [post]
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

// @Summary Получить пост по ID
// @Description Получает пост по идентификатору
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Router /posts/{id} [get]
func (ph *postHandler) FindPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := ph.PostService.FindPostByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to find post by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(post)
}

// @Summary Получить все посты от пользователя
// @Description Получает информацию о всех постах пользователя по его идентификатору
// @Tags posts
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Post
// @Router /posts/user/{id} [get]
func (ph *postHandler) GetAllPostsByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	posts, err := ph.PostService.GetAllPostsByUser(userId)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find posts by userID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "Posts not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(posts)
}
