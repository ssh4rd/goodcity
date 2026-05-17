package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/service"
)

type PracticeHandler struct {
	svc *service.PracticeService
}

func NewPracticeHandler(svc *service.PracticeService) *PracticeHandler {
	return &PracticeHandler{svc: svc}
}

func (h *PracticeHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	perPage, _ := strconv.Atoi(q.Get("per_page"))

	filter := domain.PracticeFilter{
		City:     q.Get("city"),
		Category: q.Get("category"),
		Search:   q.Get("search"),
		Status:   domain.StatusApproved,
		Page:     page,
		PerPage:  perPage,
	}

	practices, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list practices")
		return
	}
	if practices == nil {
		practices = []*domain.Practice{}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data":  practices,
		"total": total,
	})
}

func (h *PracticeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	practice, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "practice not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get practice")
		return
	}

	writeJSON(w, http.StatusOK, practice)
}

func (h *PracticeHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		City        string `json:"city"`
		Category    string `json:"category"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" || req.Description == "" || req.City == "" || req.Category == "" {
		writeError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	practice, err := h.svc.Create(r.Context(), &domain.Practice{
		Title:       req.Title,
		Description: req.Description,
		City:        req.City,
		Category:    req.Category,
		AuthorID:    user.ID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create practice")
		return
	}

	writeJSON(w, http.StatusCreated, practice)
}

func (h *PracticeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		City        string `json:"city"`
		Category    string `json:"category"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	practice, err := h.svc.Update(r.Context(), &domain.Practice{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		City:        req.City,
		Category:    req.Category,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update practice")
		return
	}

	writeJSON(w, http.StatusOK, practice)
}

func (h *PracticeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete practice")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *PracticeHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	comments, err := h.svc.ListComments(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list comments")
		return
	}
	if comments == nil {
		comments = []*domain.Comment{}
	}

	writeJSON(w, http.StatusOK, comments)
}

func (h *PracticeHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	id, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "text is required")
		return
	}

	comment, err := h.svc.AddComment(r.Context(), &domain.Comment{
		PracticeID: id,
		UserID:     user.ID,
		Text:       req.Text,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add comment")
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

func parseID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}
