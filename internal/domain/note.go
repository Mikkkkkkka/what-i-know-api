package domain

import (
	"context"
	"errors"
	"time"
)

type Note struct {
	ID        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserID    string
	Title     string
	Content   string // Markdown string
	UpdatedAt time.Time
}

type NoteRepository interface {
	GetByID(ctx context.Context, id string) (*Note, error)
	GetByUserID(ctx context.Context, userID string) ([]*Note, error)
	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id string) error
}

var ErrNoteNotFound = errors.New("mark not found")
