package api

import (
	"net/http"

	"what-i-know-api/internal/usecase"
)

type createSessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) createSession(w http.ResponseWriter, r *http.Request) {
	var request createSessionRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	session, err := h.services.Auth.CreateSession(r.Context(), usecase.CreateSessionRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newSessionResponse(session))
}

func (h *Handler) validateSession(w http.ResponseWriter, r *http.Request) {
	token, err := bearerToken(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.services.Auth.Authenticate(r.Context(), usecase.AuthRequest{Token: token}); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
