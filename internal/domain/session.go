package domain

import (
	"context"
)

type Session struct {
	Id       int64
	UserId   int64
	Username string
	Token    string
}

type SessionRepository interface {
	GetByID(ctx context.Context, id int64) (*Session, error)
	GetByToken(ctx context.Context, token string) (*Session, error)
	Create(ctx context.Context, session *Session) error
	Delete(ctx context.Context, session *Session) error
}
