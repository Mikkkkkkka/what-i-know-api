package api

import (
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/domain"
	"github.com/mikkkkkkka/what-i-know-api/internal/service"
)

type NoteHandler struct {
	notes *service.NoteService
}

func NewNoteHandler(notes *service.NoteService) *NoteHandler {
	return &NoteHandler{notes: notes}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	var request createNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	err := h.notes.CreateNote(r.Context(), service.CreateNoteRequest{
		ID:      request.ID,
		UserID:  userID,
		Title:   request.Title,
		Content: request.Content,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		WriteError(w, err)
		return
	}

	note, err := h.notes.GetByID(r.Context(), noteID)
	if err != nil {
		WriteError(w, err)
		return
	}

	if userID != note.UserID {
		WriteError(w, domain.ErrNoteNotFound)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponse(note))
}

func (h *NoteHandler) ListNotesByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	notes, err := h.notes.GetByUserID(r.Context(), userID)
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponses(notes))
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		WriteError(w, err)
		return
	}

	var request updateNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		WriteError(w, err)
		return
	}

	err = h.notes.UpdateNote(r.Context(), service.UpdateNoteRequest{
		ID:      noteID,
		UserID:  userID,
		Title:   request.Title,
		Content: request.Content,
	})
	if err != nil {
		WriteError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		WriteError(w, ErrInternal)
		return
	}

	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		WriteError(w, err)
		return
	}

	if err := h.notes.DeleteNote(r.Context(), noteID, userID); err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
