package handler

import (
	"net/http"

	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/service"
)

type ModerationHandler struct {
	svc *service.PracticeService
}

func NewModerationHandler(svc *service.PracticeService) *ModerationHandler {
	return &ModerationHandler{svc: svc}
}

func (h *ModerationHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	practices, _, err := h.svc.ListPending(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list pending practices")
		return
	}
	if practices == nil {
		practices = []*domain.Practice{}
	}
	writeJSON(w, http.StatusOK, practices)
}

func (h *ModerationHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Approve(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to approve practice")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "approved"})
}

func (h *ModerationHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Reject(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to reject practice")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}
