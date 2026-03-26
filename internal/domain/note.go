package domain

import (
	"context"
	"time"
)

type Note struct {
	Id        int64
	UserId    int64
	Name      string
	Content   string // Markdown string
	Date      time.Time
	UpdatedAt time.Time
}

type NoteRepository interface {
	GetById(ctx context.Context, id int64) (*Note, error)
	GetByUserId(ctx context.Context, userId int64) ([]*Note, error)
	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id int64) error
}
