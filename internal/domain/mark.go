package domain

import (
	"context"
	"time"
)

type Mark struct {
	Id        int64
	UserId    int64
	Date      time.Time
	Content   string // Markdown string
	UpdatedAt time.Time
}

type MarkRepository interface {
	GetById(ctx context.Context, id int64) (*Mark, error)
	GetByUserId(ctx context.Context, userId int64) ([]*Mark, error)
	Create(ctx context.Context, mark *Mark) error
	Update(ctx context.Context, mark *Mark) error
	Delete(ctx context.Context, id int64) error
}
