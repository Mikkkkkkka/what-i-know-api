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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	id, err := h.users.CreateUser(r.Context(), service.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newUserResponse(user))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	var request updateUserRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = h.users.UpdateUser(r.Context(), service.UpdateUserRequest{
		ID:       userID,
		Username: request.Username,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.users.DeleteUser(r.Context(), userID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
