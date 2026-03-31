package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type createNoteRequest struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Handler) registerNoteRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.With().Post("/", h.createNote)
		r.With().Get("/{noteID}", h.getNote)
		r.With().Patch("/{noteID}", h.updateNote)
		r.With().Delete("/{noteID}", h.deleteNote)
	})
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var request createNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := h.services.Notes.CreateNote(r.Context(), usecase.CreateNoteRequest{
		ID:      request.ID,
		UserID:  request.UserID,
		Title:   request.Title,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	note, err := h.services.Notes.GetByID(r.Context(), noteID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponse(note))
}

func (h *Handler) listNotesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	notes, err := h.services.Notes.GetByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponses(notes))
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	var request updateNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err = h.services.Notes.UpdateNote(r.Context(), usecase.UpdateNoteRequest{
		ID:      noteID,
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
	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.services.Notes.DeleteNote(r.Context(), noteID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
