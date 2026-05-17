package domain

import "time"

type Role string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
