package domain

import "time"

type PracticeStatus string

const (
	StatusPending      PracticeStatus = "pending"
	StatusApproved     PracticeStatus = "approved"
	StatusRejected     PracticeStatus = "rejected"
	StatusBestPractice PracticeStatus = "best_practice"
)

type Practice struct {
	ID            int64          `json:"id"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	City          string         `json:"city"`
	Category      string         `json:"category"`
	Address       string         `json:"address,omitempty"`
	Latitude      *float64       `json:"latitude,omitempty"`
	Longitude     *float64       `json:"longitude,omitempty"`
	BudgetRub     *float64       `json:"budget_rub,omitempty"`
	ImplementedAt *time.Time     `json:"implemented_at,omitempty"`
	Status        PracticeStatus `json:"status"`
	AuthorID      int64          `json:"author_id"`
	VoteCount     int            `json:"vote_count"`
	Rating        float64        `json:"rating"`
	CreatedAt     time.Time      `json:"created_at"`
}

type PracticeFilter struct {
	City     string
	Category string
	Search   string
	Status   PracticeStatus
	Page     int
	PerPage  int
}
