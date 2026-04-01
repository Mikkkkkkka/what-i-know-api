package service

import (
	"context"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type CreateUserRequest struct {
	Username string
	Password string
}

type UpdateUserRequest struct {
	ID       string
	Username string
	// updating password isn't implemented yet
}

type UserService struct {
	userRepo       domain.UserRepository
	idGenerator    IDGenerator
	passwordHasher PasswordHasher
}

func NewUserService(users domain.UserRepository, idGenerator IDGenerator, passwordHasher PasswordHasher) *UserService {
	return &UserService{
		userRepo:       users,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
	}
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (string, error) {
	id, err := s.idGenerator.Generate()
	if err != nil {
		return "", err
	}

	hashedPassword, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		ID:        id,
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	user.Username = req.Username

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
