package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type NoteService interface {
	GetById(ctx context.Context, id string) (*domain.Note, error)
	GetByUserId(ctx context.Context, userId string) ([]*domain.Note, error)
	CreateNote(ctx context.Context, req CreateNoteRequest) error
	UpdateNote(ctx context.Context, req UpdateNoteRequest) error
	DeleteNote(ctx context.Context, id string) error
}

type CreateNoteRequest struct {
	Id      string
	UserId  string
	Title   string
	Content string
}

type UpdateNoteRequest struct {
	Id      string
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

func (s *NoteUseCase) GetById(ctx context.Context, id string) (*domain.Note, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.notesRepo.GetById(ctx, id)
}

func (s *NoteUseCase) GetByUserId(ctx context.Context, userId string) ([]*domain.Note, error) {
	if strings.TrimSpace(userId) == "" {
		return nil, domain.ErrInvalidInput
	}

	return s.notesRepo.GetByUserId(ctx, userId)
}

func (s *NoteUseCase) CreateNote(ctx context.Context, req CreateNoteRequest) error {
	id := strings.TrimSpace(req.Id)
	userID := strings.TrimSpace(req.UserId)
	title := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if id == "" || userID == "" || title == "" || content == "" {
		return domain.ErrInvalidInput
	}

	note := &domain.Note{
		Id:        id,
		UserId:    userID,
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
	if strings.TrimSpace(req.Id) == "" {
		return domain.ErrInvalidInput
	}

	name := strings.TrimSpace(req.Title)
	content := strings.TrimSpace(req.Content)
	if name == "" || content == "" {
		return domain.ErrInvalidInput
	}

	note, err := s.notesRepo.GetById(ctx, req.Id)
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
