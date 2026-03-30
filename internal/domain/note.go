package domain

import (
	"context"
	"time"
)

type Note struct {
	Id        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserId    string
	Title     string
	Content   string // Markdown string
	UpdatedAt time.Time
}

type NoteRepository interface {
	GetById(ctx context.Context, id string) (*Note, error)
	GetByUserId(ctx context.Context, userId string) ([]*Note, error)
	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id string) error
}
