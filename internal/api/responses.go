package api

import (
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type userResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
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

type markResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Date      time.Time `json:"date"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newMarkResponse(mark *domain.Mark) markResponse {
	return markResponse{
		ID:        mark.ID,
		UserID:    mark.UserID,
		Date:      mark.Date,
		Content:   mark.Content,
		UpdatedAt: mark.UpdatedAt,
	}
}

func newMarkResponses(marks []*domain.Mark) []markResponse {
	response := make([]markResponse, 0, len(marks))
	for _, mark := range marks {
		response = append(response, newMarkResponse(mark))
	}

	return response
}
