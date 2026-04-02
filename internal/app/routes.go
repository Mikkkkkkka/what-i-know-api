package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mikkkkkkka/what-i-know-api/internal/api"
	"github.com/mikkkkkkka/what-i-know-api/internal/config"
	mw "github.com/mikkkkkkka/what-i-know-api/internal/middleware"
)

func SetupRouter(cfg config.Config, authHandler *api.AuthHandler, userHandler *api.UserHandler, noteHandler *api.NoteHandler, markHandler *api.MarkHandler, authMiddleware *mw.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Route(cfg.HTTPAPIBasePath, func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/register", authHandler.Register)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/{userID}", userHandler.GetUser)

			r.With(authMiddleware.AuthMiddleware).Patch("/{userID}", userHandler.UpdateUser)
			r.With(authMiddleware.AuthMiddleware).Delete("/{userID}", userHandler.DeleteUser)

			r.With(authMiddleware.AuthMiddleware).Get("/{userID}/notes", noteHandler.ListNotesByUser)
			r.With(authMiddleware.AuthMiddleware).Get("/{userID}/marks", markHandler.ListMarksByUser)
		})

		r.Route("/notes", func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware)
			r.Post("/", noteHandler.CreateNote)
			r.Get("/{noteID}", noteHandler.GetNote)
			r.Patch("/{noteID}", noteHandler.UpdateNote)
			r.Delete("/{noteID}", noteHandler.DeleteNote)
		})

		r.Route("/marks", func(r chi.Router) {
			r.Use(authMiddleware.AuthMiddleware)
			r.Post("/", markHandler.CreateMark)
			r.Get("/{markID}", markHandler.GetMark)
			r.Patch("/{markID}", markHandler.UpdateMark)
			r.Delete("/{markID}", markHandler.DeleteMark)
		})
	})

	return r
}
