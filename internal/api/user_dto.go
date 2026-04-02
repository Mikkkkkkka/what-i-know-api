package api

import (
	"time"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

type updateUserRequest struct {
	Username string `json:"username"`
}

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
