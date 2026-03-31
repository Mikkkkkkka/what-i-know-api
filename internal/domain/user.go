package domain

import (
	"context"
	"strings"
	"time"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

func NewUser(id, username, hashedPassword string, createdAt time.Time) (*User, error) {
	id = strings.TrimSpace(id)
	username = strings.TrimSpace(username)
	hashedPassword = strings.TrimSpace(hashedPassword)
	if id == "" || username == "" || hashedPassword == "" || createdAt.IsZero() {
		return nil, ErrInvalidInput
	}

	return &User{
		ID:        id,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: createdAt.UTC(),
	}, nil
}

func (u *User) Rename(username string) error {
	if u == nil {
		return ErrInvalidInput
	}

	username = strings.TrimSpace(username)
	if username == "" {
		return ErrInvalidInput
	}

	u.Username = username
	return nil
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
}
