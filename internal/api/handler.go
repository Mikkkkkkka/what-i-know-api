package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type Services struct {
	Users *usecase.UserUseCase
	Notes *usecase.NoteUseCase
	Marks *usecase.MarkUseCase
}

type Handler struct {
	services Services
}

func NewHandler(services Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	h.registerUserRoutes(r)
	h.registerNoteRoutes(r)
	h.registerMarkRoutes(r)
}
