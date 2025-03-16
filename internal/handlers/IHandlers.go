package handlers

import "github.com/go-chi/chi/v5"

type IHandler interface {
	Register(router *chi.Mux)
}
