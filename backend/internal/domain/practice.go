package domain

import "time"

type PracticeStatus string

const (
	StatusPending  PracticeStatus = "pending"
	StatusApproved PracticeStatus = "approved"
	StatusRejected PracticeStatus = "rejected"
)

type Practice struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	City        string         `json:"city"`
	Category    string         `json:"category"`
	Status      PracticeStatus `json:"status"`
	AuthorID    int64          `json:"author_id"`
	VoteCount   int            `json:"vote_count"`
	CreatedAt   time.Time      `json:"created_at"`
}

type PracticeFilter struct {
	City     string
	Category string
	Status   PracticeStatus
	Page     int
	PerPage  int
}
