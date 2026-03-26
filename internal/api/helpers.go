package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"what-i-know-api/internal/domain"
)

type errorResponse struct {
	Error string `json:"error"`
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return domain.ErrInvalidInput
	}
	if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
		return domain.ErrInvalidInput
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
	case errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrAlreadyExists):
		status = http.StatusConflict
		message = err.Error()
	}

	writeJSON(w, status, errorResponse{Error: message})
}

func urlParamInt64(r *http.Request, key string) (int64, error) {
	value := chi.URLParam(r, key)
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, domain.ErrInvalidInput
	}

	return id, nil
}

func bearerToken(r *http.Request) (string, error) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return "", domain.ErrInvalidInput
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", domain.ErrInvalidInput
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", domain.ErrInvalidInput
	}

	return token, nil
}
