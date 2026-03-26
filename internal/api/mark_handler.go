package api

import (
	"net/http"
	"time"

	"what-i-know-api/internal/usecase"
)

type createMarkRequest struct {
	UserID  int64     `json:"user_id"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

type updateMarkRequest struct {
	Content string `json:"content"`
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
	markID, err := urlParamInt64(r, "markID")
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
	userID, err := urlParamInt64(r, "userID")
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
	markID, err := urlParamInt64(r, "markID")
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
	markID, err := urlParamInt64(r, "markID")
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
