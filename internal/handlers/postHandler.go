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
	getPostURL     = "/posts/{id}"
	getAllPostURL  = "/posts"
	createPostURL  = "/posts"
	getPostsByUser = "/posts/user"
	deletePost     = "/posts/{id}"
	updatePostURL  = "/posts/update"
)

type postHandler struct {
	PostService services.IPostService
	Log         *slog.Logger
}

func NewPostHandler(postService services.IPostService, log *slog.Logger) IHandler {
	lg := logger.WithHandler(log, "PostHandler")
	return &postHandler{postService, lg}
}

func (h *postHandler) Register(router *chi.Mux) {
	router.Post(createPostURL, h.CreatePost)
	router.Get(getPostURL, h.FindPostByID)
	router.Get(getPostsByUser, h.GetAllPostsByUser)
	router.Delete(deletePost, h.DeletePost)
	router.Put(updatePostURL, h.UpdatePost)
}

// @Summary Создание поста
// @Description Создаёт пост и добавялет его в базу данных
// @Tags posts
// @Accept json
// @Produce json
// @Param postInfo body dto.CreatePostDTO true "Пост"
// @Success 200
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to create posts"
// @Router /posts [post]
func (h *postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.PostService.CreatePost(&dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create person: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Post not found"
// @Failure 500 {string} string "Failed to find posts"
// @Router /posts/{id} [get]
func (h *postHandler) FindPostByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	post, err := h.PostService.FindPostByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Failed to find post by id: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// @Summary Получить все посты от пользователя
// @Description Получает информацию о всех постах пользователя по его идентификатору
// @Tags posts
// @Produce json
// @Param userid query int true "userid"
// @Param limit query int true "limit"
// @Param offset query int true "offset"
// @Success 200 {array} models.Post
// @Failure 400 {string} string "Invalid userID value"
// @Failure 400 {string} string "Invalid limit value"
// @Failure 400 {string} string "Invalid offset value"
// @Failure 404 {string} string "Posts not found"
// @Failure 500 {string} string "Failed to find posts"
// @Router /posts/user [get]
func (h *postHandler) GetAllPostsByUser(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	userid, err := strconv.Atoi(queryParams.Get("userid"))
	if err != nil {
		http.Error(w, "Invalid user id value", http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		http.Error(w, "Invalid offset value", http.StatusBadRequest)
		return
	}

	posts, err := h.PostService.GetAllPostsByUser(userid, limit, offset)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to find posts by userID: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "Posts not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// @Summary Удаление поста
// @Description Удаляет пост
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to delete posts"
// @Router /posts/{id} [delete]
func (h *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.PostService.DeletePost(postId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete post: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Обновление поста
// @Description Обновляет пост
// @Tags posts
// @Produce json
// @Param post body dto.UpdatePostDTO true "Новые данные поста"
// @Success 204
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to update person"
// @Router /posts/update [put]
func (h *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	newPost := dto.UpdatePostDTO{}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	err = h.PostService.UpdatePost(&newPost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update post: %s", err.Error()), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
