package usecase

import (
	"context"
	"strings"
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
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.marksRepo.GetByID(ctx, id)
}

func (s *MarkUseCase) GetByUserID(ctx context.Context, userID string) ([]*domain.Mark, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.marksRepo.GetByUserID(ctx, userID)
}

func (s *MarkUseCase) CreateMark(ctx context.Context, req CreateMarkRequest) error {
	id := strings.TrimSpace(req.ID)
	userID := strings.TrimSpace(req.UserID)
	content := strings.TrimSpace(req.Content)
	if id == "" || userID == "" || req.Date.IsZero() || content == "" {
		return domain.ErrInvalidInput
	}

	mark := &domain.Mark{
		ID:        id,
		UserID:    userID,
		Date:      req.Date.UTC(),
		Content:   content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.marksRepo.Create(ctx, mark); err != nil {
		return err
	}

	return nil
}

func (s *MarkUseCase) UpdateMark(ctx context.Context, req UpdateMarkRequest) error {
	if strings.TrimSpace(req.ID) == "" {
		return domain.ErrInvalidInput
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return domain.ErrInvalidInput
	}

	mark, err := s.marksRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	mark.Content = content
	mark.UpdatedAt = time.Now().UTC()

	return s.marksRepo.Update(ctx, mark)
}

func (s *MarkUseCase) DeleteMark(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.ErrInvalidInput
	}

	return s.marksRepo.Delete(ctx, id)
}
