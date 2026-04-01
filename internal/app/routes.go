package app

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mikkkkkkka/what-i-know-api/internal/api"
	"github.com/mikkkkkkka/what-i-know-api/internal/config"
)

func SetupRouter(cfg config.Config, httpHandler *api.Handler) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/health"))

	basePath := strings.TrimSpace(cfg.HTTPAPIBasePath)
	if basePath == "" || basePath == "/" {
		httpHandler.RegisterRoutes(router)
		return router
	}

	router.Route(basePath, func(r chi.Router) {
		httpHandler.RegisterRoutes(r)
	})

	return router
}
