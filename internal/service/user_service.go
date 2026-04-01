package service

import (
	"context"
	"errors"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"gorm.io/gorm"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

type IDGenerator interface {
	Generate() (string, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID string) error
}

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
	userRepo       UserRepository
	idGenerator    IDGenerator
	passwordHasher PasswordHasher
}

func NewUserService(users UserRepository, idGenerator IDGenerator, passwordHasher PasswordHasher) *UserService {
	return &UserService{
		userRepo:       users,
		idGenerator:    idGenerator,
		passwordHasher: passwordHasher,
	}
}

func (s *UserService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", domain.ErrUsernameAlreadyExists
		}

		return "", err
	}

	return user.ID, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}

		return err
	}

	user.Username = req.Username

	if err := s.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrUsernameAlreadyExists
		}

		return err
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
