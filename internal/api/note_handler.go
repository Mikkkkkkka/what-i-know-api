package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"what-i-know-api/internal/domain"
	"what-i-know-api/internal/usecase"
)

type createNoteRequest struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Handler) registerNoteRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.With(h.requireAuth).Post("/", h.createNote)
		r.With(h.requireAuth).Get("/{noteID}", h.getNote)
		r.With(h.requireAuth).Patch("/{noteID}", h.updateNote)
		r.With(h.requireAuth).Delete("/{noteID}", h.deleteNote)
	})
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var request createNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	if request.UserID != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	id, err := h.services.Notes.CreateNote(r.Context(), usecase.CreateNoteRequest{
		UserId:  request.UserID,
		Title:   request.Name,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	note, err := h.services.Notes.GetById(r.Context(), noteID)
	if err != nil {
		writeError(w, err)
		return
	}
	if note.UserId != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponse(note))
}

func (h *Handler) listNotesByUser(w http.ResponseWriter, r *http.Request) {
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

	notes, err := h.services.Notes.GetByUserId(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponses(notes))
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	note, err := h.services.Notes.GetById(r.Context(), noteID)
	if err != nil {
		writeError(w, err)
		return
	}
	if note.UserId != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	var request updateNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = h.services.Notes.UpdateNote(r.Context(), usecase.UpdateNoteRequest{
		Id:      noteID,
		Title:   request.Title,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	session, err := currentSession(r)
	if err != nil {
		writeError(w, err)
		return
	}

	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	note, err := h.services.Notes.GetById(r.Context(), noteID)
	if err != nil {
		writeError(w, err)
		return
	}
	if note.UserId != session.UserId {
		writeError(w, domain.ErrForbidden)
		return
	}

	if err := h.services.Notes.DeleteNote(r.Context(), noteID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
