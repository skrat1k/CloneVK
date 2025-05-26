package handlers

import (
	"CloneVK/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	getGlobalFeed   = "/feed/global"
	getPersonalFeed = "/feed/personal"
)

type feedHandler struct {
	FeedService services.IFeedService
}

func NewFeedHandler(feedSrvice services.IFeedService) IHandler {
	return &feedHandler{feedSrvice}
}

func (h *feedHandler) Register(router *chi.Mux) {
	router.Get(getGlobalFeed, h.GetGlobalFeed)
	router.Get(getPersonalFeed, h.GetPersonalFeed)
}

// @Summary Получить глобальную ленту
// @Description Получает глобальную ленту
// @Tags feed
// @Produce json
// @Param limit query int true "ограничение количества постов"
// @Param offset query int true "пропуск первых n постов"
// @Success 200 {array} models.Follow
// @Failure 400 {string} string "Invalid limit value"
// @Failure 400 {string} string "Invalid offset value"
// @Failure 404 {string} string "No posts in feed"
// @Failure 500 {string} string "Failed to get global feed"
// @Router /feed/global [get]
func (h *feedHandler) GetGlobalFeed(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
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

	posts, err := h.FeedService.GetGlobalFeed(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get posts: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "No posts in feed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// @Summary Получить глобальную ленту
// @Description Получает персональную ленту
// @Tags feed
// @Produce json
// @Param userid query int true "id получателя ленты"
// @Param limit query int true "ограничение количества постов"
// @Param offset query int true "пропуск первых n постов"
// @Success 200 {array} models.Follow
// @Failure 400 {string} string "Invalid user id value"
// @Failure 400 {string} string "Invalid limit value"
// @Failure 400 {string} string "Invalid offset value"
// @Failure 404 {string} string "No posts in feed"
// @Failure 500 {string} string "Failed to get personal feed"
// @Router /feed/personal [get]
func (h *feedHandler) GetPersonalFeed(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.FeedService.GetPersonalFeed(userid, limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get personal feed for user %d: %s", userid, err.Error()), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "No posts in feed", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
