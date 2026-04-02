package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type userRepoStub struct {
	created *domain.User
	gotUser *domain.User
}

func (s *userRepoStub) GetByID(_ context.Context, _ string) (*domain.User, error) {
	if s.gotUser == nil {
		return nil, errors.New("user not found in stub")
	}

	return s.gotUser, nil
}

func (s *userRepoStub) GetByUsername(_ context.Context, _ string) (*domain.User, error) {
	return nil, nil
}

func (s *userRepoStub) Create(_ context.Context, user *domain.User) error {
	s.created = user
	return nil
}

func (s *userRepoStub) Update(_ context.Context, user *domain.User) error {
	s.created = user
	return nil
}

func (s *userRepoStub) Delete(_ context.Context, _ string) error {
	return nil
}

type idGeneratorStub struct {
	id string
}

func (s idGeneratorStub) Generate() (string, error) {
	return s.id, nil
}

type passwordHasherStub struct {
	hashed       string
	lastPassword string
}

func (s *passwordHasherStub) Hash(password string) (string, error) {
	s.lastPassword = password
	return s.hashed, nil
}

func (s *passwordHasherStub) Compare(_, _ string) error {
	return nil
}

func TestUserServiceCreateUserRejectsInvalidInput(t *testing.T) {
	repo := &userRepoStub{}
	hasher := &passwordHasherStub{hashed: "hashed"}
	svc := NewUserService(repo, idGeneratorStub{id: "user-id"}, hasher)

	_, err := svc.CreateUser(context.Background(), CreateUserRequest{
		Username: " ",
		Password: "secret",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository create not to be called")
	}
	if hasher.lastPassword != "" {
		t.Fatal("expected password hashing not to be called")
	}
}

func TestUserServiceCreateUserNormalizesUsernameAndPassword(t *testing.T) {
	repo := &userRepoStub{}
	hasher := &passwordHasherStub{hashed: "hashed"}
	svc := NewUserService(repo, idGeneratorStub{id: "user-id"}, hasher)

	id, err := svc.CreateUser(context.Background(), CreateUserRequest{
		Username: " user ",
		Password: " secret ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "user-id" {
		t.Fatalf("expected generated id, got %q", id)
	}
	if hasher.lastPassword != "secret" {
		t.Fatalf("expected original password to be hashed, got %q", hasher.lastPassword)
	}
	if repo.created == nil {
		t.Fatal("expected repository create to be called")
	}
	if repo.created.Username != "user" {
		t.Fatalf("expected normalized username, got %q", repo.created.Username)
	}
}

func TestUserServiceUpdateUserRejectsWhitespaceOnlyUsername(t *testing.T) {
	repo := &userRepoStub{
		gotUser: &domain.User{ID: "user-id", Username: "old"},
	}
	svc := NewUserService(repo, idGeneratorStub{id: "user-id"}, &passwordHasherStub{hashed: "hashed"})

	err := svc.UpdateUser(context.Background(), UpdateUserRequest{
		ID:       " user-id ",
		Username: " ",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository update not to be called")
	}
}
