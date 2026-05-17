package domain

import "time"

type Rating struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	PracticeID      int64     `json:"practice_id"`
	ScoreUtility    int       `json:"score_utility"`
	ScoreAesthetics int       `json:"score_aesthetics"`
	ScoreEcology    int       `json:"score_ecology"`
	Comment         string    `json:"comment,omitempty"`
	UserLat         *float64  `json:"user_lat,omitempty"`
	UserLon         *float64  `json:"user_lon,omitempty"`
	IsSuspicious    bool      `json:"is_suspicious"`
	CreatedAt       time.Time `json:"created_at"`
}

type RatingStats struct {
	Count           int     `json:"count"`
	AvgUtility      float64 `json:"avg_utility"`
	AvgAesthetics   float64 `json:"avg_aesthetics"`
	AvgEcology      float64 `json:"avg_ecology"`
	Rating          float64 `json:"rating"`
}
