package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mikkkkkkka/what-i-know-api/internal/api"
	"github.com/mikkkkkkka/what-i-know-api/internal/config"
	"github.com/mikkkkkkka/what-i-know-api/internal/repository/gorm_postgres"
	"github.com/mikkkkkkka/what-i-know-api/internal/security"
	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

func Start() {
	cfg := config.Load()

	if cfg.DatabaseDSN == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	db, err := gorm_postgres.OpenPostgres(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}

	if err := gorm_postgres.AutoMigrate(db); err != nil {
		log.Fatalf("auto-migrate database: %v", err)
	}

	userRepository := gorm_postgres.NewUserRepository(db)
	noteRepository := gorm_postgres.NewNoteRepository(db)
	markRepository := gorm_postgres.NewMarkRepository(db)

	idGenerator := security.NewUUIDGenerator()
	passwordHasher := security.NewBcryptPasswordHasher(0)

	userService := usecase.NewUserUseCase(userRepository, idGenerator, passwordHasher)
	noteService := usecase.NewNoteUseCase(noteRepository)
	markService := usecase.NewMarkUseCase(markRepository)

	userHandler := api.NewUserHandler(userService)
	noteHandler := api.NewNoteHandler(noteService)
	markHandler := api.NewMarkHandler(markService)

	router := SetupRouter(userHandler, noteHandler, markHandler)

	if err := http.ListenAndServe(cfg.HTTPAddress, router); err != nil {
		log.Fatal(err)
	}
}
