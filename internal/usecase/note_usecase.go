package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type NoteService interface {
	GetByID(ctx context.Context, id string) (*domain.Note, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error)
	CreateNote(ctx context.Context, req CreateNoteRequest) error
	UpdateNote(ctx context.Context, req UpdateNoteRequest) error
	DeleteNote(ctx context.Context, id string) error
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

type NoteUseCase struct {
	notesRepo domain.NoteRepository
}

func NewNoteService(notes domain.NoteRepository) *NoteUseCase {
	return &NoteUseCase{
		notesRepo: notes,
	}
}

func (s *NoteUseCase) GetByID(ctx context.Context, id string) (*domain.Note, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.notesRepo.GetByID(ctx, id)
}

func (s *NoteUseCase) GetByUserID(ctx context.Context, userID string) ([]*domain.Note, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.notesRepo.GetByUserID(ctx, userID)
}

func (s *NoteUseCase) CreateNote(ctx context.Context, req CreateNoteRequest) error {
	id := strings.TrimSpace(req.ID)
	userID := strings.TrimSpace(req.UserID)
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if id == "" || userID == "" || title == "" || content == "" {
		return domain.ErrInvalidInput
	}

	note := &domain.Note{
		ID:        id,
		UserID:    userID,
		Title:     title,
		Content:   content,
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.notesRepo.Create(ctx, note); err != nil {
		return err
	}

	return nil
}

func (s *NoteUseCase) UpdateNote(ctx context.Context, req UpdateNoteRequest) error {
	if strings.TrimSpace(req.ID) == "" {
		return domain.ErrInvalidInput
	}

	name := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if name == "" || content == "" {
		return domain.ErrInvalidInput
	}

	note, err := s.notesRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	note.Title = name
	note.Content = content
	note.UpdatedAt = time.Now().UTC()

	return s.notesRepo.Update(ctx, note)
}

func (s *NoteUseCase) DeleteNote(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return domain.ErrInvalidInput
	}

	return s.notesRepo.Delete(ctx, id)
}

var _ NoteService = (*NoteUseCase)(nil)
