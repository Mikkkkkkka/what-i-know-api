package api

import (
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type createMarkRequest struct {
	ID      string    `json:"id"`
	UserID  string    `json:"user_id"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

type updateMarkRequest struct {
	Content string `json:"content"`
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
