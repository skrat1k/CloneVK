package handlers

import (
	dto "CloneVK/internal/dto/follows"
	"CloneVK/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
)

type followHandler struct {
	FollowService services.IFollowService
}

const (
	createFollowURL        = "/follows"
	getAllFollowsURL       = "/follows"
	getAllUserFollowsURL   = "/follow/{id}"
	getAllUserFollowersURL = "/followers/{id}"
	deleteFollowURL        = "/follows"
)

func NewFollowHandler(followService services.IFollowService) IHandler {
	return &followHandler{followService}
}

func (h *followHandler) Register(router *chi.Mux) {
	router.Post(createFollowURL, h.CreateFollow)
	router.Get(getAllFollowsURL, h.GetAllFollows)
	router.Get(getAllUserFollowsURL, h.GetAllUserFollows)
	router.Get(getAllUserFollowersURL, h.GetAllUserFollowers)
	router.Delete(deleteFollowURL, h.DeleteFollow)
}

// @Summary Создать фоллов
// @Description Фолловит человека на другого
// @Tags follows
// @Accept json
// @Produce json
// @Param followInfo body dto.FollowDTO true "Фоллов"
// @Success 204
// @Failure 400 {string} string "Invalid JSON"
// @Failure 400 {string} string "Validate error"
// @Failure 500 {string} string "Failed to create follow"
// @Router /follows [post]
func (h *followHandler) CreateFollow(w http.ResponseWriter, r *http.Request) {
	var dto dto.FollowDTO

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

	err = h.FollowService.CreateFollow(dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create follow: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Получить список всех фолловов
// @Description Получить все фолловы
// @Tags follows
// @Produce json
// @Success 200 {array} models.Follow
// @Failure 404 {string} string "Follows notfound"
// @Failure 500 {string} string "Failed to get follows"
// @Router /follows [get]
func (h *followHandler) GetAllFollows(w http.ResponseWriter, r *http.Request) {
	follows, err := h.FollowService.GetAllFollows()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get follows: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(follows) == 0 {
		http.Error(w, "Follows not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(follows)
}

// @Summary Получить всех фолловеров пользователя
// @Description Получает всех фолловеров пользователя
// @Tags follows
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.Follow
// @Failure 400 {string} string "Invalid JSON"
// @Failure 404 {string} string "Followers not found"
// @Failure 500 {string} string "Failed to get followers"
// @Router /followers/{id} [get]
func (h *followHandler) GetAllUserFollowers(w http.ResponseWriter, r *http.Request) {
	followedID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	follows, err := h.FollowService.GetAllUserFollowers(followedID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get followers by user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(follows) == 0 {
		http.Error(w, "Followers not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(follows)
}

// @Summary Получить все фолловы пользователя
// @Description Получает все фолловы пользователя
// @Tags follows
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} models.Follow
// @Failure 400 {string} string "Invalid JSON"
// @Failure 404 {string} string "Follows not found"
// @Failure 500 {string} string "Failed to get follows"
// @Router /follow/{id} [get]
func (h *followHandler) GetAllUserFollows(w http.ResponseWriter, r *http.Request) {
	followerID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	follows, err := h.FollowService.GetAllUserFollows(followerID)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get follows by user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(follows) == 0 {
		http.Error(w, "Follows not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(follows)
}

// @Summary Удалить фоллов
// @Description Удаляет фоллов
// @Tags follows
// @Accept json
// @Produce json
// @Param followInfo body dto.FollowDTO true "Фоллов"
// @Success 204
// @Failure 400 {string} string "Invalid JSON"
// @Failure 400 {string} string "Validate error"
// @Failure 500 {string} string "Failed to delete follow"
// @Router /follows [delete]
func (h *followHandler) DeleteFollow(w http.ResponseWriter, r *http.Request) {
	var dto dto.FollowDTO

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

	err = h.FollowService.DeleteFollow(dto)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete follow: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
