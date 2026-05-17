package handler

import (
	"errors"
	"net/http"

	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/service"
)

type RatingHandler struct {
	svc *service.RatingService
}

func NewRatingHandler(svc *service.RatingService) *RatingHandler {
	return &RatingHandler{svc: svc}
}

func (h *RatingHandler) Rate(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	practiceID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req struct {
		ScoreUtility    int      `json:"score_utility"`
		ScoreAesthetics int      `json:"score_aesthetics"`
		ScoreEcology    int      `json:"score_ecology"`
		Comment         string   `json:"comment"`
		UserLat         *float64 `json:"user_lat"`
		UserLon         *float64 `json:"user_lon"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ScoreUtility < 1 || req.ScoreUtility > 5 ||
		req.ScoreAesthetics < 1 || req.ScoreAesthetics > 5 ||
		req.ScoreEcology < 1 || req.ScoreEcology > 5 {
		writeError(w, http.StatusBadRequest, "scores must be between 1 and 5")
		return
	}

	rating, err := h.svc.Rate(r.Context(), &domain.Rating{
		UserID:          user.ID,
		PracticeID:      practiceID,
		ScoreUtility:    req.ScoreUtility,
		ScoreAesthetics: req.ScoreAesthetics,
		ScoreEcology:    req.ScoreEcology,
		Comment:         req.Comment,
		UserLat:         req.UserLat,
		UserLon:         req.UserLon,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to save rating")
		return
	}

	writeJSON(w, http.StatusOK, rating)
}

func (h *RatingHandler) GetMyRating(w http.ResponseWriter, r *http.Request) {
	user := UserFromContext(r.Context())
	practiceID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	rating, err := h.svc.GetMyRating(r.Context(), user.ID, practiceID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeJSON(w, http.StatusOK, nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get rating")
		return
	}

	writeJSON(w, http.StatusOK, rating)
}

func (h *RatingHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	practiceID, err := parseID(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	stats, err := h.svc.GetStats(r.Context(), practiceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get rating stats")
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
