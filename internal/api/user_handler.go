package api

import (
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		WriteError(w, err)
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newUserResponse(user))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		WriteError(w, err)
		return
	}

	var request updateUserRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	err = h.users.UpdateUser(r.Context(), service.UpdateUserRequest{
		ID:       userID,
		Username: request.Username,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := h.users.DeleteUser(r.Context(), userID); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
