package usecase

import (
	"context"

	"what-i-know-api/internal/domain"
)

type UserService interface {
	GetById(ctx context.Context, id int64) (*domain.User, error)
	CreateUser(ctx context.Context, request CreateUserRequest) error
	UpdateUser(ctx context.Context, request UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int64) error
}

type CreateUserRequest struct {
	Username string
	Password string
}

type UpdateUserRequest struct {
	Id       int64
	Username string
}
