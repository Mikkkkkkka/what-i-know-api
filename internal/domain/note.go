package domain

import (
	"context"
	"strings"
	"time"
)

type Note struct {
	ID        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserID    string
	Title     string
	Content   string // Markdown string
	UpdatedAt time.Time
}

func NewNote(id, userID, title, content string, updatedAt time.Time) (*Note, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if id == "" || userID == "" || title == "" || content == "" || updatedAt.IsZero() {
		return nil, ErrInvalidInput
	}

	return &Note{
		ID:        id,
		UserID:    userID,
		Title:     title,
		Content:   content,
		UpdatedAt: updatedAt.UTC(),
	}, nil
}

func (n *Note) Update(title, content string, updatedAt time.Time) error {
	if n == nil {
		return ErrInvalidInput
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" || content == "" || updatedAt.IsZero() {
		return ErrInvalidInput
	}

	n.Title = title
	n.Content = content
	n.UpdatedAt = updatedAt.UTC()
	return nil
}

type NoteRepository interface {
	GetByID(ctx context.Context, id string) (*Note, error)
	GetByUserID(ctx context.Context, userID string) ([]*Note, error)
	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id string) error
}
