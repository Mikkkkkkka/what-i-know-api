package usecase

import (
	"context"
	"time"

	"what-i-know-api/internal/domain"
)

type MarkService interface {
	GetById(ctx context.Context, id int64) (*domain.Mark, error)
	GetByUserId(ctx context.Context, userId int64) ([]*domain.Mark, error)
	CreateMark(ctx context.Context, req CreateMarkRequest) error
	UpdateMark(ctx context.Context, req UpdateMarkRequest) error
	DeleteMark(ctx context.Context, id int64) error
}

type CreateMarkRequest struct {
	UserId  int64
	Date    time.Time
	Content string
}

type UpdateMarkRequest struct {
	Id      int64
	Content string
}
