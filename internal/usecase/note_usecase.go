package usecase

import (
	"context"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

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

type NoteUseCase struct {
	notesRepo domain.NoteRepository
}

func NewNoteUseCase(notes domain.NoteRepository) *NoteUseCase {
	return &NoteUseCase{
		notesRepo: notes,
	}
}

func (s *NoteUseCase) GetByID(ctx context.Context, id string) (*domain.Note, error) {
	return s.notesRepo.GetByID(ctx, id)
}

func (s *NoteUseCase) GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error) {
	return s.notesRepo.GetByUserID(ctx, userID)
}

func (s *NoteUseCase) CreateNote(ctx context.Context, req CreateNoteRequest) error {
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

func (s *NoteUseCase) UpdateNote(ctx context.Context, req UpdateNoteRequest) error {
	note, err := s.notesRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	note.Title = req.Title
	note.Content = req.Content
	note.UpdatedAt = time.Now().UTC()

	return s.notesRepo.Update(ctx, note)
}

func (s *NoteUseCase) DeleteNote(ctx context.Context, id string) error {
	return s.notesRepo.Delete(ctx, id)
}
