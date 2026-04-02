package service

import (
	"context"
	"errors"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"gorm.io/gorm"
)

type NoteRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Note, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error)
	Create(ctx context.Context, note *domain.Note) error
	Update(ctx context.Context, note *domain.Note) error
	Delete(ctx context.Context, id string) error
}

type CreateNoteRequest struct {
	ID      string
	UserID  string
	Title   string
	Content string
}

type UpdateNoteRequest struct {
	ID      string
	Title   string
	Content string
}

type NoteService struct {
	notesRepo NoteRepository
}

func NewNoteService(notes NoteRepository) *NoteService {
	return &NoteService{
		notesRepo: notes,
	}
}

func (s *NoteService) GetByID(ctx context.Context, id string) (*domain.Note, error) {
	note, err := s.notesRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNoteNotFound
		}

		return nil, err
	}

	return note, nil
}

func (s *NoteService) GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error) {
	return s.notesRepo.GetByUserID(ctx, userID)
}

func (s *NoteService) CreateNote(ctx context.Context, req CreateNoteRequest) error {
	id, err := normalizeRequiredString(req.ID)
	if err != nil {
		return err
	}

	userID, err := normalizeRequiredString(req.UserID)
	if err != nil {
		return err
	}

	title, err := normalizeRequiredString(req.Title)
	if err != nil {
		return err
	}

	content, err := normalizeRequiredString(req.Content)
	if err != nil {
		return err
	}

	note := &domain.Note{
		ID:        id,
		UserID:    userID,
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.notesRepo.Create(ctx, note); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrNoteAlreadyExists
		}

		return err
	}

	return nil
}

func (s *NoteService) UpdateNote(ctx context.Context, req UpdateNoteRequest) error {
	id, err := normalizeRequiredString(req.ID)
	if err != nil {
		return err
	}

	title, err := normalizeRequiredString(req.Title)
	if err != nil {
		return err
	}

	content, err := normalizeRequiredString(req.Content)
	if err != nil {
		return err
	}

	note, err := s.notesRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNoteNotFound
		}

		return err
	}

	note.Title = title
	note.Content = content
	note.UpdatedAt = time.Now().UTC()

	return s.notesRepo.Update(ctx, note)
}

func (s *NoteService) DeleteNote(ctx context.Context, id string) error {
	return s.notesRepo.Delete(ctx, id)
}
