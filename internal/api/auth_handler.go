package api

import (
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	token, err := h.auth.Login(r.Context(), service.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request registerRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	err := h.auth.Register(r.Context(), service.RegisterRequest{
		Username: request.Username,
		Password: request.Password,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "registered"})
}
