package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"gorm.io/gorm"
)

type authUserRepoStub struct {
	gotUsername string
	gotUser     *domain.User
	getErr      error
	created     *domain.User
	createErr   error
}

func (s *authUserRepoStub) GetByID(_ context.Context, _ string) (*domain.User, error) {
	return nil, errors.New("not implemented in authUserRepoStub")
}

func (s *authUserRepoStub) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	s.gotUsername = username
	if s.getErr != nil {
		return nil, s.getErr
	}

	return s.gotUser, nil
}

func (s *authUserRepoStub) Create(_ context.Context, user *domain.User) error {
	s.created = user
	return s.createErr
}

func (s *authUserRepoStub) Update(_ context.Context, _ *domain.User) error {
	return nil
}

func (s *authUserRepoStub) Delete(_ context.Context, _ string) error {
	return nil
}

type authIDGeneratorStub struct {
	id  string
	err error
}

func (s authIDGeneratorStub) Generate() (string, error) {
	if s.err != nil {
		return "", s.err
	}

	return s.id, nil
}

type authPasswordHasherStub struct {
	hashed           string
	hashErr          error
	compareErr       error
	lastHashPassword string
	lastHashed       string
	lastCompared     string
}

func (s *authPasswordHasherStub) Hash(password string) (string, error) {
	s.lastHashPassword = password
	if s.hashErr != nil {
		return "", s.hashErr
	}

	return s.hashed, nil
}

func (s *authPasswordHasherStub) Compare(hashedPassword, password string) error {
	s.lastHashed = hashedPassword
	s.lastCompared = password
	return s.compareErr
}

type tokenManagerStub struct {
	token      string
	err        error
	lastUserID string
}

func (s *tokenManagerStub) Generate(userID string) (string, error) {
	s.lastUserID = userID
	if s.err != nil {
		return "", s.err
	}

	return s.token, nil
}

func TestAuthServiceLoginRejectsInvalidUsername(t *testing.T) {
	repo := &authUserRepoStub{}
	hasher := &authPasswordHasherStub{}
	tokenGenerator := &tokenManagerStub{}
	svc := &AuthService{
		userRepo:       repo,
		jwtManager:     tokenGenerator,
		passwordHasher: hasher,
	}

	_, err := svc.Login(context.Background(), LoginRequest{
		Username: " ",
		Password: "secret",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.gotUsername != "" {
		t.Fatalf("expected repository lookup not to be called, got %q", repo.gotUsername)
	}
	if hasher.lastCompared != "" {
		t.Fatal("expected password comparison not to be called")
	}
	if tokenGenerator.lastUserID != "" {
		t.Fatal("expected token generation not to be called")
	}
}

func TestAuthServiceLoginMapsUserNotFound(t *testing.T) {
	repo := &authUserRepoStub{getErr: gorm.ErrRecordNotFound}
	svc := &AuthService{
		userRepo:       repo,
		jwtManager:     &tokenManagerStub{},
		passwordHasher: &authPasswordHasherStub{},
	}

	_, err := svc.Login(context.Background(), LoginRequest{
		Username: " user ",
		Password: "secret",
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if repo.gotUsername != "user" {
		t.Fatalf("expected normalized username lookup, got %q", repo.gotUsername)
	}
}

func TestAuthServiceLoginRejectsInvalidPassword(t *testing.T) {
	repo := &authUserRepoStub{
		gotUser: &domain.User{ID: "user-id", Username: "user", Password: "hashed"},
	}
	hasher := &authPasswordHasherStub{}
	tokenGenerator := &tokenManagerStub{}
	svc := &AuthService{
		userRepo:       repo,
		jwtManager:     tokenGenerator,
		passwordHasher: hasher,
	}

	_, err := svc.Login(context.Background(), LoginRequest{
		Username: "user",
		Password: " ",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if hasher.lastCompared != "" {
		t.Fatal("expected password comparison not to be called")
	}
	if tokenGenerator.lastUserID != "" {
		t.Fatal("expected token generation not to be called")
	}
}

func TestAuthServiceLoginReturnsIncorrectCredentialsOnCompareError(t *testing.T) {
	repo := &authUserRepoStub{
		gotUser: &domain.User{ID: "user-id", Username: "user", Password: "hashed"},
	}
	hasher := &authPasswordHasherStub{compareErr: errors.New("bad password")}
	tokenGenerator := &tokenManagerStub{}
	svc := &AuthService{
		userRepo:       repo,
		jwtManager:     tokenGenerator,
		passwordHasher: hasher,
	}

	_, err := svc.Login(context.Background(), LoginRequest{
		Username: " user ",
		Password: " secret ",
	})
	if !errors.Is(err, domain.ErrIncorrectCredentials) {
		t.Fatalf("expected ErrIncorrectCredentials, got %v", err)
	}
	if hasher.lastHashed != "hashed" || hasher.lastCompared != "secret" {
		t.Fatalf("unexpected compare arguments: hashed=%q password=%q", hasher.lastHashed, hasher.lastCompared)
	}
	if tokenGenerator.lastUserID != "" {
		t.Fatal("expected token generation not to be called")
	}
}

func TestAuthServiceLoginReturnsJWTToken(t *testing.T) {
	repo := &authUserRepoStub{
		gotUser: &domain.User{ID: "user-id", Username: "user", Password: "hashed"},
	}
	hasher := &authPasswordHasherStub{}
	tokenGenerator := &tokenManagerStub{token: "jwt-token"}
	svc := &AuthService{
		userRepo:       repo,
		jwtManager:     tokenGenerator,
		passwordHasher: hasher,
	}

	token, err := svc.Login(context.Background(), LoginRequest{
		Username: " user ",
		Password: " secret ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token to be generated")
	}
	if token != "jwt-token" {
		t.Fatalf("expected stub token, got %q", token)
	}
	if repo.gotUsername != "user" {
		t.Fatalf("expected normalized username lookup, got %q", repo.gotUsername)
	}
	if hasher.lastCompared != "secret" {
		t.Fatalf("expected normalized password comparison, got %q", hasher.lastCompared)
	}
	if tokenGenerator.lastUserID != "user-id" {
		t.Fatalf("expected token generation for user-id, got %q", tokenGenerator.lastUserID)
	}
}

func TestAuthServiceRegisterRejectsInvalidInput(t *testing.T) {
	repo := &authUserRepoStub{}
	hasher := &authPasswordHasherStub{}
	svc := &AuthService{
		userRepo:       repo,
		idGenerator:    authIDGeneratorStub{id: "user-id"},
		passwordHasher: hasher,
	}

	err := svc.Register(context.Background(), RegisterRequest{
		Username: " ",
		Password: "secret",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository create not to be called")
	}
	if hasher.lastHashPassword != "" {
		t.Fatal("expected password hashing not to be called")
	}
}

func TestAuthServiceRegisterNormalizesStringsAndCreatesUser(t *testing.T) {
	repo := &authUserRepoStub{}
	hasher := &authPasswordHasherStub{hashed: "hashed"}
	svc := &AuthService{
		userRepo:       repo,
		idGenerator:    authIDGeneratorStub{id: "user-id"},
		passwordHasher: hasher,
	}

	err := svc.Register(context.Background(), RegisterRequest{
		Username: " user ",
		Password: " secret ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasher.lastHashPassword != "secret" {
		t.Fatalf("expected trimmed password to be hashed, got %q", hasher.lastHashPassword)
	}
	if repo.created == nil {
		t.Fatal("expected repository create to be called")
	}
	if repo.created.ID != "user-id" || repo.created.Username != "user" || repo.created.Password != "hashed" {
		t.Fatalf("unexpected created user: %+v", repo.created)
	}
	if repo.created.CreatedAt.Location() != time.UTC {
		t.Fatalf("expected UTC timestamp, got %v", repo.created.CreatedAt.Location())
	}
}

func TestAuthServiceRegisterReturnsUsernameAlreadyExists(t *testing.T) {
	repo := &authUserRepoStub{createErr: gorm.ErrDuplicatedKey}
	svc := &AuthService{
		userRepo:       repo,
		idGenerator:    authIDGeneratorStub{id: "user-id"},
		passwordHasher: &authPasswordHasherStub{hashed: "hashed"},
	}

	err := svc.Register(context.Background(), RegisterRequest{
		Username: "user",
		Password: "secret",
	})
	if !errors.Is(err, domain.ErrUsernameAlreadyExists) {
		t.Fatalf("expected ErrUsernameAlreadyExists, got %v", err)
	}
}

func TestAuthServiceRegisterPropagatesIDGenerationError(t *testing.T) {
	expectedErr := errors.New("generate id")
	repo := &authUserRepoStub{}
	hasher := &authPasswordHasherStub{}
	svc := &AuthService{
		userRepo:       repo,
		idGenerator:    authIDGeneratorStub{err: expectedErr},
		passwordHasher: hasher,
	}

	err := svc.Register(context.Background(), RegisterRequest{
		Username: "user",
		Password: "secret",
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected id generation error, got %v", err)
	}
	if hasher.lastHashPassword != "" {
		t.Fatal("expected password hashing not to be called")
	}
	if repo.created != nil {
		t.Fatal("expected repository create not to be called")
	}
}

func TestAuthServiceRegisterPropagatesHashError(t *testing.T) {
	expectedErr := errors.New("hash password")
	repo := &authUserRepoStub{}
	hasher := &authPasswordHasherStub{hashErr: expectedErr}
	svc := &AuthService{
		userRepo:       repo,
		idGenerator:    authIDGeneratorStub{id: "user-id"},
		passwordHasher: hasher,
	}

	err := svc.Register(context.Background(), RegisterRequest{
		Username: "user",
		Password: "secret",
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected hash error, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository create not to be called")
	}
}
