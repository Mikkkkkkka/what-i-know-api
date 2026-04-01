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

	httpHandler := api.NewHandler(api.Services{
		Users: userService,
		Notes: noteService,
		Marks: markService,
	})

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      SetupRouter(cfg, httpHandler),
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	log.Printf("http server listening on %s", cfg.HTTPAddress)

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
	case sig := <-shutdownSignals:
		log.Printf("received signal %s, shutting down", sig)
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), cfg.HTTPShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		log.Fatalf("shutdown http server: %v", err)
	}
}
