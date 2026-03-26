package api

import (
	"net/http"

	"what-i-know-api/internal/usecase"
)

type createNoteRequest struct {
	UserID  int64  `json:"user_id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	var request createNoteRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, err)
		return
	}

	err := h.services.Notes.CreateNote(r.Context(), usecase.CreateNoteRequest{
		UserId:  request.UserID,
		Name:    request.Name,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (h *Handler) getNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := urlParamInt64(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	note, err := h.services.Notes.GetById(r.Context(), noteID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newNoteResponse(note))
}

func (h *Handler) listNotesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := urlParamInt64(r, "userID")
	if err != nil {
		writeError(w, err)
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
	noteID, err := urlParamInt64(r, "noteID")
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
		Id:      noteID,
		Name:    request.Name,
		Content: request.Content,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	noteID, err := urlParamInt64(r, "noteID")
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.services.Notes.DeleteNote(r.Context(), noteID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
