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
	mark, err := domain.NewMark(req.ID, req.UserID, req.Date, req.Content, time.Now().UTC())
	if err != nil {
		return err
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

	if err := mark.UpdateContent(content, time.Now().UTC()); err != nil {
		return err
	}

	return s.marksRepo.Update(ctx, mark)
}

func (s *MarkUseCase) DeleteMark(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.ErrInvalidInput
	}

	return s.marksRepo.Delete(ctx, id)
}
