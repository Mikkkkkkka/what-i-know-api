package domain

import (
	"context"
	"time"
)

type Mark struct {
	Id        string // Generated on client via UUID v4. Slim (1/2^122) chance of collision!
	UserId    string
	Date      time.Time
	Content   string // Markdown string
	UpdatedAt time.Time
}

type MarkRepository interface {
	GetById(ctx context.Context, id string) (*Mark, error)
	GetByUserId(ctx context.Context, userId string) ([]*Mark, error)
	Create(ctx context.Context, mark *Mark) error
	Update(ctx context.Context, mark *Mark) error
	Delete(ctx context.Context, id string) error
}
