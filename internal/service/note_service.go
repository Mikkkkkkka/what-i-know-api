package service

import (
	"context"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
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
	return s.notesRepo.GetByID(ctx, id)
}

func (s *NoteService) GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error) {
	return s.notesRepo.GetByUserID(ctx, userID)
}

func (s *NoteService) CreateNote(ctx context.Context, req CreateNoteRequest) error {
	note := &domain.Note{
		ID:        req.ID,
		UserID:    req.UserID,
		Title:     req.Title,
		Content:   req.Content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.notesRepo.Create(ctx, note); err != nil {
		return err
	}

	return nil
}

func (s *NoteService) UpdateNote(ctx context.Context, req UpdateNoteRequest) error {
	note, err := s.notesRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	note.Title = req.Title
	note.Content = req.Content
	note.UpdatedAt = time.Now().UTC()

	return s.notesRepo.Update(ctx, note)
}

func (s *NoteService) DeleteNote(ctx context.Context, id string) error {
	return s.notesRepo.Delete(ctx, id)
}
