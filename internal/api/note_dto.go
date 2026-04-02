package api

import (
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type createNoteRequest struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type noteResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newNoteResponse(note *domain.Note) noteResponse {
	return noteResponse{
		ID:        note.ID,
		UserID:    note.UserID,
		Title:     note.Title,
		Content:   note.Content,
		UpdatedAt: note.UpdatedAt,
	}
}

func newNoteResponses(notes []*domain.Note) []noteResponse {
	response := make([]noteResponse, 0, len(notes))
	for _, note := range notes {
		response = append(response, newNoteResponse(note))
	}

	return response
}
