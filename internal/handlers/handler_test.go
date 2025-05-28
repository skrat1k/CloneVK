package handlers

import (
	"bytes"
	"CloneVK/internal/dto/posts"
	"CloneVK/internal/models"
	"CloneVK/internal/services/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
)

func TestCreatePost_Success(t *testing.T) {
	mockService := new(mocks.IPostService)
	handler := &postHandler{PostService: mockService, Log: slog.Default()}

	post := posts.CreatePostDTO{Text: "test post", UserID: 1}
	postJSON, _ := json.Marshal(post)

	mockService.On("CreatePost", &post).Return(123, nil)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(postJSON))
	w := httptest.NewRecorder()

	handler.CreatePost(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp int
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, 123, resp)
}

func TestCreatePost_InvalidJSON(t *testing.T) {
	handler := &postHandler{Log: slog.Default()}

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader([]byte("{invalid-json")))
	w := httptest.NewRecorder()

	handler.CreatePost(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindPostByID_Success(t *testing.T) {
	mockService := new(mocks.IPostService)
	handler := &postHandler{PostService: mockService, Log: slog.Default()}

	expectedPost := &models.Post{ID: 1, Text: "test"}
	mockService.On("FindPostByID", 1).Return(expectedPost, nil)

	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(chi.NewRouteContextContext(req.Context(), rctx))

	w := httptest.NewRecorder()
	handler.FindPostByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Post
	_ = json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, expectedPost.Text, resp.Text)
}

func TestFindPostByID_NotFound(t *testing.T) {
	mockService := new(mocks.IPostService)
	handler := &postHandler{PostService: mockService, Log: slog.Default()}

	mockService.On("FindPostByID", 1).Return(nil, errors.New("sql: no rows in result set"))

	req := httptest.NewRequest(http.MethodGet, "/posts/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(1))
	req = req.WithContext(chi.NewRouteContextContext(req.Context(), rctx))

	w := httptest.NewRecorder()
	handler.FindPostByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
