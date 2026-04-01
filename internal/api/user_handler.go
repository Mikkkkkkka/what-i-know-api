package api

import (
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type UserHandler struct {
	users *usecase.UserUseCase
}

func NewUserHandler(users *usecase.UserUseCase) *UserHandler {
	return &UserHandler{users: users}
}

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Username string `json:"username"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	id, err := h.users.CreateUser(r.Context(), usecase.CreateUserRequest{
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

	err = h.users.UpdateUser(r.Context(), usecase.UpdateUserRequest{
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
