package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type Services struct {
	Users usecase.UserService
	Notes usecase.NoteService
	Marks usecase.MarkService
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
