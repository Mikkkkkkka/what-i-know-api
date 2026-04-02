package service

import (
	"context"
	"errors"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"gorm.io/gorm"
)

type MarkRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Mark, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Mark, error)
	Create(ctx context.Context, mark *domain.Mark) error
	Update(ctx context.Context, mark *domain.Mark) error
	Delete(ctx context.Context, id string) error
}

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

type MarkService struct {
	marksRepo MarkRepository
}

func NewMarkService(marks MarkRepository) *MarkService {
	return &MarkService{marksRepo: marks}
}

func (s *MarkService) GetByID(ctx context.Context, id string) (*domain.Mark, error) {
	mark, err := s.marksRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrMarkNotFound
		}

		return nil, err
	}

	return mark, nil
}

func (s *MarkService) GetByUserID(ctx context.Context, userID string) ([]*domain.Mark, error) {
	return s.marksRepo.GetByUserID(ctx, userID)
}

func (s *MarkService) CreateMark(ctx context.Context, req CreateMarkRequest) error {
	mark := &domain.Mark{
		ID:        req.ID,
		UserID:    req.UserID,
		Date:      req.Date.UTC(),
		Content:   req.Content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.marksRepo.Create(ctx, mark); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrMarkAlreadyExists
		}

		return err
	}

	return nil
}

func (s *MarkService) UpdateMark(ctx context.Context, req UpdateMarkRequest) error {

	mark, err := s.marksRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrMarkNotFound
		}

		return err
	}

	mark.Content = req.Content
	mark.UpdatedAt = time.Now().UTC()

	return s.marksRepo.Update(ctx, mark)
}

func (s *MarkService) DeleteMark(ctx context.Context, id string) error {
	return s.marksRepo.Delete(ctx, id)
}
