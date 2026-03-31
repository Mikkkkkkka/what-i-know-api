package domain

import (
	"context"
	"strings"
	"time"
)

type Mark struct {
	ID        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserID    string
	Date      time.Time
	Content   string // Markdown string
	UpdatedAt time.Time
}

func NewMark(id, userID string, date time.Time, content string, updatedAt time.Time) (*Mark, error) {
	id = strings.TrimSpace(id)
	userID = strings.TrimSpace(userID)
	content = strings.TrimSpace(content)
	if id == "" || userID == "" || date.IsZero() || content == "" || updatedAt.IsZero() {
		return nil, ErrInvalidInput
	}

	return &Mark{
		ID:        id,
		UserID:    userID,
		Date:      date.UTC(),
		Content:   content,
		UpdatedAt: updatedAt.UTC(),
	}, nil
}

func (m *Mark) UpdateContent(content string, updatedAt time.Time) error {
	if m == nil {
		return ErrInvalidInput
	}

	content = strings.TrimSpace(content)
	if content == "" || updatedAt.IsZero() {
		return ErrInvalidInput
	}

	m.Content = content
	m.UpdatedAt = updatedAt.UTC()
	return nil
}

type MarkRepository interface {
	GetByID(ctx context.Context, id string) (*Mark, error)
	GetByUserID(ctx context.Context, userID string) ([]*Mark, error)
	Create(ctx context.Context, mark *Mark) error
	Update(ctx context.Context, mark *Mark) error
	Delete(ctx context.Context, id string) error
}
