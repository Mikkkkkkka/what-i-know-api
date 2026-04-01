package usecase

import (
	"context"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type CreateMarkRequest struct {
	ID      string
	UserID  string
	Date    time.Time
	Content string
}

type UpdateMarkRequest struct {
	ID      string
	Content string
}

type MarkUseCase struct {
	marksRepo domain.MarkRepository
}

func NewMarkUseCase(marks domain.MarkRepository) *MarkUseCase {
	return &MarkUseCase{marksRepo: marks}
}

func (s *MarkUseCase) GetByID(ctx context.Context, id string) (*domain.Mark, error) {
	return s.marksRepo.GetByID(ctx, id)
}

func (s *MarkUseCase) GetByUserID(ctx context.Context, userID string) ([]*domain.Mark, error) {
	return s.marksRepo.GetByUserID(ctx, userID)
}

func (s *MarkUseCase) CreateMark(ctx context.Context, req CreateMarkRequest) error {
	mark := &domain.Mark{
		ID:        req.ID,
		UserID:    req.UserID,
		Date:      req.Date.UTC(),
		Content:   req.Content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.marksRepo.Create(ctx, mark); err != nil {
		return err
	}

	return nil
}

func (s *MarkUseCase) UpdateMark(ctx context.Context, req UpdateMarkRequest) error {

	mark, err := s.marksRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	mark.Content = req.Content
	mark.UpdatedAt = time.Now().UTC()

	return s.marksRepo.Update(ctx, mark)
}

func (s *MarkUseCase) DeleteMark(ctx context.Context, id string) error {
	return s.marksRepo.Delete(ctx, id)
}
