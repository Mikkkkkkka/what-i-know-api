package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID string) error
}

var ErrUserNotFound = errors.New("user not found")
var ErrUsernameAlreadyExists = errors.New("username already exists")
