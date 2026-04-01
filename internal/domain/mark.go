package domain

import (
	"context"
	"errors"
	"time"
)

type Mark struct {
	ID        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserID    string
	Date      time.Time
	Content   string // Markdown string
	UpdatedAt time.Time
}

type MarkRepository interface {
	GetByID(ctx context.Context, id string) (*Mark, error)
	GetByUserID(ctx context.Context, userID string) ([]*Mark, error)
	Create(ctx context.Context, mark *Mark) error
	Update(ctx context.Context, mark *Mark) error
	Delete(ctx context.Context, id string) error
}

var ErrMarkNotFound = errors.New("mark not found")
var ErrMarkAlreadyExists = errors.New("mark already exists")
