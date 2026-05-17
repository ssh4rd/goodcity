package domain

import "time"

type Comment struct {
	ID         int64     `json:"id"`
	PracticeID int64     `json:"practice_id"`
	UserID     int64     `json:"user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}
