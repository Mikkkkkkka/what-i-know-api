package api

import (
	"net/http"

	"what-i-know-api/internal/usecase"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Username string `json:"username"`
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := h.services.Users.CreateUser(r.Context(), usecase.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamInt64(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := h.services.Users.GetById(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newUserResponse(user))
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamInt64(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	var request updateUserRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = h.services.Users.UpdateUser(r.Context(), usecase.UpdateUserRequest{
		Id:       userID,
		Username: request.Username,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamInt64(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.services.Users.DeleteUser(r.Context(), userID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
