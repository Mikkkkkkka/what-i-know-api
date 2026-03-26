package usecase

import (
	"context"
	"strings"
	"time"

	"what-i-know-api/internal/domain"
)

type NoteService interface {
	GetById(ctx context.Context, id int64) (*domain.Note, error)
	GetByUserId(ctx context.Context, userId int64) ([]*domain.Note, error)
	CreateNote(ctx context.Context, req CreateNoteRequest) error
	UpdateNote(ctx context.Context, req UpdateNoteRequest) error
	DeleteNote(ctx context.Context, id int64) error
}

type CreateNoteRequest struct {
	UserId  int64
	Name    string
	Content string
}

type UpdateNoteRequest struct {
	Id      int64
	Name    string
	Content string
}

type NoteUseCase struct {
	notes domain.NoteRepository
}

func NewNoteService(notes domain.NoteRepository) *NoteUseCase {
	return &NoteUseCase{notes: notes}
}

func (s *NoteUseCase) GetById(ctx context.Context, id int64) (*domain.Note, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return s.notes.GetById(ctx, id)
}

func (s *NoteUseCase) GetByUserId(ctx context.Context, userId int64) ([]*domain.Note, error) {
	if userId <= 0 {
		return nil, domain.ErrInvalidInput
	}

	return s.notes.GetByUserId(ctx, userId)
}

func (s *NoteUseCase) CreateNote(ctx context.Context, req CreateNoteRequest) error {
	name := strings.TrimSpace(req.Name)
	content := strings.TrimSpace(req.Content)
	if req.UserId <= 0 || name == "" || content == "" {
		return domain.ErrInvalidInput
	}

	note := &domain.Note{
		UserId:    req.UserId,
		Name:      name,
		Content:   content,
		Date:      time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return s.notes.Create(ctx, note)
}

func (s *NoteUseCase) UpdateNote(ctx context.Context, req UpdateNoteRequest) error {
	if req.Id <= 0 {
		return domain.ErrInvalidInput
	}

	name := strings.TrimSpace(req.Name)
	content := strings.TrimSpace(req.Content)
	if name == "" || content == "" {
		return domain.ErrInvalidInput
	}

	note, err := s.notes.GetById(ctx, req.Id)
	if err != nil {
		return err
	}

	note.Name = name
	note.Content = content
	note.UpdatedAt = time.Now().UTC()

	return s.notes.Update(ctx, note)
}

func (s *NoteUseCase) DeleteNote(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	return s.notes.Delete(ctx, id)
}

var _ NoteService = (*NoteUseCase)(nil)
