package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/mikkkkkkka/what-i-know-api/internal/api"
	"github.com/mikkkkkkka/what-i-know-api/internal/auth"
	"github.com/mikkkkkkka/what-i-know-api/internal/config"
	"github.com/mikkkkkkka/what-i-know-api/internal/middleware"
	"github.com/mikkkkkkka/what-i-know-api/internal/repository/gorm_postgres"
	"github.com/mikkkkkkka/what-i-know-api/internal/security"
	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

func Start() {
	cfg := config.Load()

	if missing := cfg.MissingRequiredDBEnv(); len(missing) > 0 {
		log.Fatalf("%s are required", strings.Join(missing, ", "))
	}

	db, err := gorm_postgres.OpenPostgres(cfg.DatabaseDSN())
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
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	userService := service.NewUserService(userRepository, idGenerator, passwordHasher)
	authService := service.NewAuthService(userRepository, jwtManager, idGenerator, passwordHasher)
	noteService := service.NewNoteService(noteRepository)
	markService := service.NewMarkService(markRepository)

	authHandler := api.NewAuthHandler(authService)
	userHandler := api.NewUserHandler(userService)
	noteHandler := api.NewNoteHandler(noteService)
	markHandler := api.NewMarkHandler(markService)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router := SetupRouter(cfg, authHandler, userHandler, noteHandler, markHandler, authMiddleware)

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      router,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	log.Printf("Server listening at %v\n", cfg.HTTPAddress)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
