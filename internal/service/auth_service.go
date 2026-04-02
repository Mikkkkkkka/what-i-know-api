package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"gorm.io/gorm"
)

type TokenManager interface {
	Generate(userID string) (string, error)
}

type AuthService struct {
	userRepo       UserRepository
	jwtManager     TokenManager
	idGenerator    IDGenerator
	passwordHasher PasswordHasher
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthService(userRepo UserRepository, jwtManager TokenManager, idGenerator IDGenerator, passwordHasher PasswordHasher) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
	}
}

func (as *AuthService) Login(ctx context.Context, req LoginRequest) (string, error) {
	username, err := normalizeRequiredString(req.Username)
	if err != nil {
		return "", err
	}

	user, err := as.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", domain.ErrUserNotFound
		}

		return "", err
	}

	password, err := normalizeRequiredString(req.Password)
	if err != nil {
		return "", err
	}

	err = as.passwordHasher.Compare(user.Password, password)
	if err != nil {
		return "", domain.ErrIncorrectCredentials
	}

	return as.jwtManager.Generate(user.ID)
}

func (as *AuthService) Register(ctx context.Context, req RegisterRequest) error {
	username, err := normalizeRequiredString(req.Username)
	if err != nil {
		return err
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		return ErrInvalidInput
	}

	id, err := as.idGenerator.Generate()
	if err != nil {
		return err
	}

	hashedPassword, err := as.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:        id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	if err := as.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrUsernameAlreadyExists
		}

		return err
	}

	return nil
}
