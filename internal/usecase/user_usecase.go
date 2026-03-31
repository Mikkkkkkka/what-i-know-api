package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type UserService interface {
	GetById(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateUser(ctx context.Context, request CreateUserRequest) (string, error)
	UpdateUser(ctx context.Context, request UpdateUserRequest) error
	DeleteUser(ctx context.Context, id string) error
}

type CreateUserRequest struct {
	Username string
	Password string
}

type UpdateUserRequest struct {
	Id       string
	Username string
	// updating password isn't implemented yet
}

type UserUseCase struct {
	userRepo       domain.UserRepository
	idGenerator    IDGenerator
	passwordHasher PasswordHasher
}

func NewUserService(users domain.UserRepository, idGenerator IDGenerator, passwordHasher PasswordHasher) *UserUseCase {
	return &UserUseCase{
		userRepo:       users,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
	}
}

func (s *UserUseCase) GetById(ctx context.Context, id string) (*domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.userRepo.GetById(ctx, id)
}

func (s *UserUseCase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserUseCase) CreateUser(ctx context.Context, request CreateUserRequest) (string, error) {
	username := strings.TrimSpace(request.Username)
	password := strings.TrimSpace(request.Password)
	if username == "" || password == "" {
		return "", domain.ErrInvalidInput
	}

	id, err := s.idGenerator.Generate()
	if err != nil {
		return "", err
	}

	hashedPassword, err := s.passwordHasher.Hash(password)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Id:        id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	return user.Id, nil
}

func (s *UserUseCase) UpdateUser(ctx context.Context, request UpdateUserRequest) error {
	if strings.TrimSpace(request.Id) == "" {
		return domain.ErrInvalidInput
	}

	username := strings.TrimSpace(request.Username)
	if username == "" {
		return domain.ErrInvalidInput
	}

	user, err := s.userRepo.GetById(ctx, request.Id)
	if err != nil {
		return err
	}

	user.Username = username

	return s.userRepo.Update(ctx, user)
}

func (s *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.ErrInvalidInput
	}

	return s.userRepo.Delete(ctx, id)
}

var _ UserService = (*UserUseCase)(nil)
