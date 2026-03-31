package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"what-i-know-api/internal/domain"
	"what-i-know-api/internal/usecase"
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
		r.With(h.requireAuth).Post("/", h.createMark)
		r.With(h.requireAuth).Get("/{markID}", h.getMark)
		r.With(h.requireAuth).Patch("/{markID}", h.updateMark)
		r.With(h.requireAuth).Delete("/{markID}", h.deleteMark)
	})
}

func (h *Handler) createMark(w http.ResponseWriter, r *http.Request) {
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request createMarkRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}
	if request.UserID != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	err = h.services.Marks.CreateMark(r.Context(), usecase.CreateMarkRequest{
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
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

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
	if mark.UserId != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	writeJSON(w, http.StatusOK, newMarkResponse(mark))
}

func (h *Handler) listMarksByUser(w http.ResponseWriter, r *http.Request) {
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}
	if userID != session.UserId {
		writeError(w, domain.ErrForbidden)
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
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

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
	if mark.UserId != session.UserId {
		writeError(w, domain.ErrForbidden)
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
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

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
	if mark.UserId != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	if err := h.services.Marks.DeleteMark(r.Context(), markID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
