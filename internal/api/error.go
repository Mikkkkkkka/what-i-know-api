package api

import (
	"errors"
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

var ErrInvalidInput = errors.New("invalid input")

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal server error"

	switch {
	case errors.Is(err, ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, service.ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		status = http.StatusForbidden
		message = err.Error()
	case errors.Is(err, domain.ErrIncorrectCredentials):
		status = http.StatusUnauthorized
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
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrNoteAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	case errors.Is(err, domain.ErrMarkAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	}

	writeJSON(w, status, errorResponse{Error: message})
}
