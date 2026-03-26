package usecase

import (
	"context"
	"strings"
	"time"

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

type UserUseCase struct {
	users          domain.UserRepository
	passwordHasher PasswordHasher
}

func NewUserService(users domain.UserRepository, passwordHasher PasswordHasher) *UserUseCase {
	return &UserUseCase{
		users:          users,
		passwordHasher: passwordHasher,
	}
}

func (s *UserUseCase) GetById(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return s.users.GetById(ctx, id)
}

func (s *UserUseCase) CreateUser(ctx context.Context, request CreateUserRequest) error {
	username := strings.TrimSpace(request.Username)
	password := strings.TrimSpace(request.Password)
	if username == "" || password == "" {
		return domain.ErrInvalidInput
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	return s.users.Create(ctx, user)
}

func (s *UserUseCase) UpdateUser(ctx context.Context, request UpdateUserRequest) error {
	if request.Id <= 0 {
		return domain.ErrInvalidInput
	}

	username := strings.TrimSpace(request.Username)
	if username == "" {
		return domain.ErrInvalidInput
	}

	user, err := s.users.GetById(ctx, request.Id)
	if err != nil {
		return err
	}

	user.Username = username

	return s.users.Update(ctx, user)
}

func (s *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	return s.users.Delete(ctx, id)
}

var _ UserService = (*UserUseCase)(nil)
