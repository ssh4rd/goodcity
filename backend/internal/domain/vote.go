package domain

type Vote struct {
	UserID     int64 `json:"user_id"`
	PracticeID int64 `json:"practice_id"`
}
