package domain

import (
	"context"
	"time"
)

type User struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt time.Time
}

type UserRepository interface {
	GetById(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userId int64) error
}
