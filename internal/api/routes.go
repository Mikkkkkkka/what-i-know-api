package api

import "github.com/go-chi/chi/v5"

func (h *Handler) registerUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/{userID}", h.getUser)
		r.Patch("/{userID}", h.updateUser)
		r.Delete("/{userID}", h.deleteUser)

		r.Get("/{userID}/notes", h.listNotesByUser)
		r.Get("/{userID}/marks", h.listMarksByUser)
	})
}

func (h *Handler) registerNoteRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Post("/", h.createNote)
		r.Get("/{noteID}", h.getNote)
		r.Patch("/{noteID}", h.updateNote)
		r.Delete("/{noteID}", h.deleteNote)
	})
}

func (h *Handler) registerMarkRoutes(r chi.Router) {
	r.Route("/marks", func(r chi.Router) {
		r.Post("/", h.createMark)
		r.Get("/{markID}", h.getMark)
		r.Patch("/{markID}", h.updateMark)
		r.Delete("/{markID}", h.deleteMark)
	})
}

func (h *Handler) registerAuthRoutes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sessions", h.createSession)
		r.Get("/validate", h.validateSession)
	})
}
