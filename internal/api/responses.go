package api

import (
	"time"

	"what-i-know-api/internal/domain"
)

type userResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

type noteResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newNoteResponse(note *domain.Note) noteResponse {
	return noteResponse{
		ID:        note.Id,
		UserID:    note.UserId,
		Name:      note.Name,
		Content:   note.Content,
		Date:      note.Date,
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
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Date      time.Time `json:"date"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newMarkResponse(mark *domain.Mark) markResponse {
	return markResponse{
		ID:        mark.Id,
		UserID:    mark.UserId,
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

type sessionResponse struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

func newSessionResponse(session domain.Session) sessionResponse {
	return sessionResponse{
		ID:       session.Id,
		UserID:   session.UserId,
		Username: session.Username,
		Token:    session.Token,
	}
}
