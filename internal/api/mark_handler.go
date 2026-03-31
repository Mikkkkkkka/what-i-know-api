package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type createMarkRequest struct {
	UserID  string    `json:"user_id"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

type updateMarkRequest struct {
	Content string `json:"content"`
}

func (h *Handler) registerMarkRoutes(r chi.Router) {
	r.Route("/marks", func(r chi.Router) {
		r.With().Post("/", h.createMark)
		r.With().Get("/{markID}", h.getMark)
		r.With().Patch("/{markID}", h.updateMark)
		r.With().Delete("/{markID}", h.deleteMark)
	})
}

func (h *Handler) createMark(w http.ResponseWriter, r *http.Request) {
	var request createMarkRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := h.services.Marks.CreateMark(r.Context(), usecase.CreateMarkRequest{
		UserId:  request.UserID,
		Date:    request.Date,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *Handler) getMark(w http.ResponseWriter, r *http.Request) {
	markID, err := urlParamString(r, "markID")
	if err != nil {
		writeError(w, err)
		return
	}

	mark, err := h.services.Marks.GetById(r.Context(), markID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newMarkResponse(mark))
}

func (h *Handler) listMarksByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	marks, err := h.services.Marks.GetByUserId(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newMarkResponses(marks))
}

func (h *Handler) updateMark(w http.ResponseWriter, r *http.Request) {
	markID, err := urlParamString(r, "markID")
	if err != nil {
		writeError(w, err)
		return
	}

	_, err = h.services.Marks.GetById(r.Context(), markID)
	if err != nil {
		writeError(w, err)
		return
	}

	var request updateMarkRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = h.services.Marks.UpdateMark(r.Context(), usecase.UpdateMarkRequest{
		Id:      markID,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) deleteMark(w http.ResponseWriter, r *http.Request) {
	markID, err := urlParamString(r, "markID")
	if err != nil {
		writeError(w, err)
		return
	}

	_, err = h.services.Marks.GetById(r.Context(), markID)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.services.Marks.DeleteMark(r.Context(), markID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
