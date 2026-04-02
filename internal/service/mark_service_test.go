package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type markRepoStub struct {
	created *domain.Mark
	gotNote *domain.Mark
}

func (s *markRepoStub) GetByID(_ context.Context, _ string) (*domain.Mark, error) {
	if s.gotNote == nil {
		return nil, errors.New("mark not found in stub")
	}

	return s.gotNote, nil
}

func (s *markRepoStub) GetByUserID(_ context.Context, _ string) ([]*domain.Mark, error) {
	return nil, nil
}

func (s *markRepoStub) Create(_ context.Context, mark *domain.Mark) error {
	s.created = mark
	return nil
}

func (s *markRepoStub) Update(_ context.Context, mark *domain.Mark) error {
	s.created = mark
	return nil
}

func (s *markRepoStub) Delete(_ context.Context, _ string) error {
	return nil
}

func TestMarkServiceCreateMarkRejectsInvalidInput(t *testing.T) {
	repo := &markRepoStub{}
	svc := NewMarkService(repo)

	err := svc.CreateMark(context.Background(), CreateMarkRequest{
		ID:      "mark-id",
		UserID:  " ",
		Date:    time.Now(),
		Content: "content",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository create not to be called")
	}
}

func TestMarkServiceCreateMarkNormalizesStrings(t *testing.T) {
	repo := &markRepoStub{}
	svc := NewMarkService(repo)

	err := svc.CreateMark(context.Background(), CreateMarkRequest{
		ID:      " mark-id ",
		UserID:  " user-id ",
		Date:    time.Date(2026, 4, 2, 12, 0, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
		Content: "  content  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created == nil {
		t.Fatal("expected repository create to be called")
	}
	if repo.created.ID != "mark-id" || repo.created.UserID != "user-id" || repo.created.Content != "content" {
		t.Fatalf("unexpected normalized mark: %+v", repo.created)
	}
	if repo.created.Date.Location() != time.UTC {
		t.Fatalf("expected UTC date, got %v", repo.created.Date.Location())
	}
}

func TestMarkServiceUpdateMarkRejectsWhitespaceOnlyContent(t *testing.T) {
	repo := &markRepoStub{
		gotNote: &domain.Mark{ID: "mark-id", Content: "old"},
	}
	svc := NewMarkService(repo)

	err := svc.UpdateMark(context.Background(), UpdateMarkRequest{
		ID:      " mark-id ",
		Content: " ",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("expected repository update not to be called")
	}
}
