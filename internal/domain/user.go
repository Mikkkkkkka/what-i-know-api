package domain

import (
	"context"
	"time"
)

type User struct {
	Id        string
	Username  string
	Password  string
	CreatedAt time.Time
}

type UserRepository interface {
	GetById(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userId string) error
}
