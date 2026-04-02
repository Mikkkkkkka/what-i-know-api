package api

import (
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

type MarkHandler struct {
	marks *service.MarkService
}

func NewMarkHandler(marks *service.MarkService) *MarkHandler {
	return &MarkHandler{marks: marks}
}

func (h *MarkHandler) CreateMark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	var request createMarkRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	err := h.marks.CreateMark(r.Context(), service.CreateMarkRequest{
		ID:      request.ID,
		UserID:  userID,
		Date:    request.Date,
		Content: request.Content,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *MarkHandler) GetMark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	markID, err := urlParamString(r, "markID")
	if err != nil {
		WriteError(w, err)
		return
	}

	mark, err := h.marks.GetByID(r.Context(), markID)
	if err != nil {
		WriteError(w, err)
		return
	}

	if userID != mark.UserID {
		WriteError(w, domain.ErrMarkNotFound)
		return
	}

	writeJSON(w, http.StatusOK, newMarkResponse(mark))
}

func (h *MarkHandler) ListMarksByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	marks, err := h.marks.GetByUserID(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newMarkResponses(marks))
}

func (h *MarkHandler) UpdateMark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	markID, err := urlParamString(r, "markID")
	if err != nil {
		WriteError(w, err)
		return
	}

	var request updateMarkRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	err = h.marks.UpdateMark(r.Context(), service.UpdateMarkRequest{
		ID:      markID,
		UserID:  userID,
		Content: request.Content,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *MarkHandler) DeleteMark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	markID, err := urlParamString(r, "markID")
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := h.marks.DeleteMark(r.Context(), markID, userID); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
