package usecase

import (
	"context"
	"strings"

	"what-i-know-api/internal/domain"
)

type AuthService interface {
	Authenticate(ctx context.Context, req AuthRequest) error
	CreateSession(ctx context.Context, request CreateSessionRequest) (domain.Session, error)
}

type CreateSessionRequest struct {
	Username string
	Password string
}

type AuthRequest struct {
	Token string
}

type AuthUseCase struct {
	users          domain.UserRepository
	sessions       domain.SessionRepository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

func NewAuthService(
	users domain.UserRepository,
	sessions domain.SessionRepository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
) *AuthUseCase {
	return &AuthUseCase{
		users:          users,
		sessions:       sessions,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (s *AuthUseCase) Authenticate(ctx context.Context, req AuthRequest) error {
	token := strings.TrimSpace(req.Token)
	if token == "" {
		return domain.ErrInvalidInput
	}

	_, err := s.sessions.GetByToken(ctx, token)
	return err
}

func (s *AuthUseCase) CreateSession(ctx context.Context, request CreateSessionRequest) (domain.Session, error) {
	username := strings.TrimSpace(request.Username)
	password := strings.TrimSpace(request.Password)
	if username == "" || password == "" {
		return domain.Session{}, domain.ErrInvalidInput
	}

	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return domain.Session{}, err
	}

	if err := s.passwordHasher.Compare(user.Password, password); err != nil {
		return domain.Session{}, domain.ErrInvalidInput
	}

	token, err := s.tokenGenerator.Generate()
	if err != nil {
		return domain.Session{}, err
	}

	session := domain.Session{
		UserId:   user.Id,
		Username: user.Username,
		Token:    token,
	}

	if err := s.sessions.Create(ctx, &session); err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

var _ AuthService = (*AuthUseCase)(nil)
