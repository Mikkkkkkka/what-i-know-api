package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"what-i-know-api/internal/api"
	"what-i-know-api/internal/config"
	"what-i-know-api/internal/repository"
	"what-i-know-api/internal/security"
	"what-i-know-api/internal/server"
	"what-i-know-api/internal/usecase"
)

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if appConfig.Database.DSN == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	db, err := repository.OpenPostgres(appConfig.Database.DSN)
	if err != nil {
		log.Fatalf("open postgres: %v", err)
	}

	if err := repository.AutoMigrate(db); err != nil {
		log.Fatalf("auto-migrate database: %v", err)
	}

	userRepository := repository.NewUserRepository(db)
	noteRepository := repository.NewNoteRepository(db)
	markRepository := repository.NewMarkRepository(db)
	sessionRepository := repository.NewSessionRepository(db)

	idGenerator := security.NewUUIDGenerator()
	passwordHasher := security.NewBcryptPasswordHasher(0)
	tokenGenerator := security.NewRandomTokenGenerator(0)

	userService := usecase.NewUserService(userRepository, idGenerator, passwordHasher)
	noteService := usecase.NewNoteService(noteRepository, idGenerator)
	markService := usecase.NewMarkService(markRepository)
	authService := usecase.NewAuthService(userRepository, sessionRepository, passwordHasher, tokenGenerator)

	apiHandler := api.NewHandler(api.Services{
		Users: userService,
		Notes: noteService,
		Marks: markService,
		Auth:  authService,
	})

	httpServer := server.New(appConfig.HTTP, apiHandler)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- httpServer.ListenAndServe()
	}()

	log.Printf("http server listening on %s", appConfig.HTTP.Address)

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

	shutdownContext, cancel := context.WithTimeout(context.Background(), appConfig.HTTP.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		log.Fatalf("shutdown http server: %v", err)
	}
}
