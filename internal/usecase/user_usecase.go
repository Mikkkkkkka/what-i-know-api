package usecase

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

type UserUseCase struct {
	userRepo       domain.UserRepository
	idGenerator    IDGenerator
	passwordHasher PasswordHasher
}

func NewUserUseCase(users domain.UserRepository, idGenerator IDGenerator, passwordHasher PasswordHasher) *UserUseCase {
	return &UserUseCase{
		userRepo:       users,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
	}
}

func (s *UserUseCase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserUseCase) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (string, error) {
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

func (s *UserUseCase) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	user.Username = req.Username

	return s.userRepo.Update(ctx, user)
}

func (s *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
