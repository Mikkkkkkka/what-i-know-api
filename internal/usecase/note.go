package usecase

import (
	"context"

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
