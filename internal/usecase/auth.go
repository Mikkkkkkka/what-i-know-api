package usecase

import (
	"context"

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
