package api

import (
	"net/http"

	"github.com/mikkkkkkka/what-i-know-api/internal/usecase"
)

type NoteHandler struct {
	notes *usecase.NoteUseCase
}

func NewNoteHandler(notes *usecase.NoteUseCase) *NoteHandler {
	return &NoteHandler{notes: notes}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var request createNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := h.notes.CreateNote(r.Context(), usecase.CreateNoteRequest{
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

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	note, err := h.notes.GetByID(r.Context(), noteID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponse(note))
}

func (h *NoteHandler) ListNotesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamString(r, "userID")
	if err != nil {
		writeError(w, err)
		return
	}

	notes, err := h.notes.GetByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponses(notes))
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
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

	err = h.notes.UpdateNote(r.Context(), usecase.UpdateNoteRequest{
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

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := urlParamString(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.notes.DeleteNote(r.Context(), noteID); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
