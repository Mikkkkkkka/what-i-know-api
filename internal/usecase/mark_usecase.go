package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type MarkService interface {
	GetById(ctx context.Context, id string) (*domain.Mark, error)
	GetByUserId(ctx context.Context, userId string) ([]*domain.Mark, error)
	CreateMark(ctx context.Context, req CreateMarkRequest) error
	UpdateMark(ctx context.Context, req UpdateMarkRequest) error
	DeleteMark(ctx context.Context, id string) error
}

type CreateMarkRequest struct {
	Id      string
	UserId  string
	Date    time.Time
	Content string
}

type UpdateMarkRequest struct {
	Id      string
	Content string
}

type MarkUseCase struct {
	marksRepo domain.MarkRepository
}

func NewMarkService(marks domain.MarkRepository) *MarkUseCase {
	return &MarkUseCase{marksRepo: marks}
}

func (s *MarkUseCase) GetById(ctx context.Context, id string) (*domain.Mark, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.marksRepo.GetById(ctx, id)
}

func (s *MarkUseCase) GetByUserId(ctx context.Context, userId string) ([]*domain.Mark, error) {
	if strings.TrimSpace(userId) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.marksRepo.GetByUserId(ctx, userId)
}

func (s *MarkUseCase) CreateMark(ctx context.Context, req CreateMarkRequest) error {
	userID := strings.TrimSpace(req.UserId)
	content := strings.TrimSpace(req.Content)
	if userID == "" || req.Date.IsZero() || content == "" {
		return domain.ErrInvalidInput
	}

	mark := &domain.Mark{
		UserId:    userID,
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
	if strings.TrimSpace(req.Id) == "" {
		return domain.ErrInvalidInput
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		return domain.ErrInvalidInput
	}

	mark, err := s.marksRepo.GetById(ctx, req.Id)
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

var _ MarkService = (*MarkUseCase)(nil)
