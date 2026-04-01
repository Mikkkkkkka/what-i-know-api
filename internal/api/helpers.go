package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
)

var ErrInvalidInput = errors.New("invalid input")

type errorResponse struct {
	Error string `json:"error"`
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return ErrInvalidInput
	}
	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return ErrInvalidInput
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		message = err.Error()
	case errors.Is(err, domain.ErrUserNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrNoteNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrMarkNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrMarkAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	}

	writeJSON(w, status, errorResponse{Error: message})
}

func urlParamString(r *http.Request, key string) (string, error) {
	value := strings.TrimSpace(chi.URLParam(r, key))
	if value == "" {
		return "", ErrInvalidInput
	}

	return value, nil
}
