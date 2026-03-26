package usecase

import (
	"context"
	"strings"
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

type MarkUseCase struct {
	marks domain.MarkRepository
}

func NewMarkService(marks domain.MarkRepository) *MarkUseCase {
	return &MarkUseCase{marks: marks}
}

func (s *MarkUseCase) GetById(ctx context.Context, id int64) (*domain.Mark, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return s.marks.GetById(ctx, id)
}

func (s *MarkUseCase) GetByUserId(ctx context.Context, userId int64) ([]*domain.Mark, error) {
	if userId <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return s.marks.GetByUserId(ctx, userId)
}

func (s *MarkUseCase) CreateMark(ctx context.Context, req CreateMarkRequest) error {
	content := strings.TrimSpace(req.Content)
	if req.UserId <= 0 || req.Date.IsZero() || content == "" {
		return domain.ErrInvalidInput
	}

	mark := &domain.Mark{
		UserId:    req.UserId,
		Date:      req.Date.UTC(),
		Content:   content,
		UpdatedAt: time.Now().UTC(),
	}

	return s.marks.Create(ctx, mark)
}

func (s *MarkUseCase) UpdateMark(ctx context.Context, req UpdateMarkRequest) error {
	if req.Id <= 0 {
		return domain.ErrInvalidInput
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return domain.ErrInvalidInput
	}

	mark, err := s.marks.GetById(ctx, req.Id)
	if err != nil {
		return err
	}

	mark.Content = content
	mark.UpdatedAt = time.Now().UTC()

	return s.marks.Update(ctx, mark)
}

func (s *MarkUseCase) DeleteMark(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	return s.marks.Delete(ctx, id)
}

var _ MarkService = (*MarkUseCase)(nil)
