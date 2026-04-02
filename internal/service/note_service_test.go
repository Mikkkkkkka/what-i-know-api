package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type noteRepoStub struct {
	created   *domain.Note
	gotID     string
	gotNote   *domain.Note
	deletedID string
}

func (s *noteRepoStub) GetByID(_ context.Context, id string) (*domain.Note, error) {
	s.gotID = id
	if s.gotNote == nil {
		return nil, errors.New("note not found in stub")
	}

	return s.gotNote, nil
}

func (s *noteRepoStub) GetByUserID(_ context.Context, _ string) ([]*domain.Note, error) {
	return nil, nil
}

func (s *noteRepoStub) Create(_ context.Context, note *domain.Note) error {
	s.created = note
	return nil
}

func (s *noteRepoStub) Update(_ context.Context, note *domain.Note) error {
	s.created = note
	return nil
}

func (s *noteRepoStub) Delete(_ context.Context, id string) error {
	s.deletedID = id
	return nil
}

func TestNoteServiceCreateNoteRejectsInvalidInput(t *testing.T) {
	repo := &noteRepoStub{}
	svc := NewNoteService(repo)

	err := svc.CreateNote(context.Background(), CreateNoteRequest{
		ID:      " ",
		UserID:  "user",
		Title:   "title",
		Content: "content",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository create not to be called")
	}
}

func TestNoteServiceCreateNoteNormalizesStrings(t *testing.T) {
	repo := &noteRepoStub{}
	svc := NewNoteService(repo)

	err := svc.CreateNote(context.Background(), CreateNoteRequest{
		ID:      " note-id ",
		UserID:  " user-id ",
		Title:   " title ",
		Content: "  content  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created == nil {
		t.Fatal("expected repository create to be called")
	}
	if repo.created.ID != "note-id" || repo.created.UserID != "user-id" || repo.created.Title != "title" || repo.created.Content != "content" {
		t.Fatalf("unexpected normalized note: %+v", repo.created)
	}
}

func TestNoteServiceUpdateNoteRejectsWhitespaceOnlyFields(t *testing.T) {
	repo := &noteRepoStub{
		gotNote: &domain.Note{ID: "note-id", Title: "old", Content: "old"},
	}
	svc := NewNoteService(repo)

	err := svc.UpdateNote(context.Background(), UpdateNoteRequest{
		ID:      " note-id ",
		UserID:  "user-id",
		Title:   " ",
		Content: "content",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository update not to be called")
	}
}

func TestNoteServiceUpdateNoteRejectsDifferentOwner(t *testing.T) {
	repo := &noteRepoStub{
		gotNote: &domain.Note{ID: "note-id", UserID: "owner-id", Title: "old", Content: "old"},
	}
	svc := NewNoteService(repo)

	err := svc.UpdateNote(context.Background(), UpdateNoteRequest{
		ID:      "note-id",
		UserID:  "other-user",
		Title:   "title",
		Content: "content",
	})
	if !errors.Is(err, domain.ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository update not to be called")
	}
}

func TestNoteServiceDeleteNoteRejectsDifferentOwner(t *testing.T) {
	repo := &noteRepoStub{
		gotNote: &domain.Note{ID: "note-id", UserID: "owner-id"},
	}
	svc := NewNoteService(repo)

	err := svc.DeleteNote(context.Background(), "note-id", "other-user")
	if !errors.Is(err, domain.ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got %v", err)
	}
	if repo.deletedID != "" {
		t.Fatal("expected repository delete not to be called")
	}
}
