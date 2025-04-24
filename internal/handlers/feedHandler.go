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
	getGlobalFeed = "/feed/global"
)

type feedHandler struct {
	FeedService services.IFeedService
}

func NewFeedHandler(feedSrvice services.IFeedService) IHandler {
	return &feedHandler{feedSrvice}
}

func (h *feedHandler) Register(router *chi.Mux) {
	router.Get(getGlobalFeed, h.GetGlobalFeed)
}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
