package domain

import "time"

type Role string
type IncomeBracket string
type SocialRole string

const (
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

const (
	IncomeLow    IncomeBracket = "low"
	IncomeMiddle IncomeBracket = "middle"
	IncomeHigh   IncomeBracket = "high"
)

const (
	SocialResident   SocialRole = "resident"
	SocialActivist   SocialRole = "activist"
	SocialSpecialist SocialRole = "specialist"
	SocialOfficial   SocialRole = "official"
)

type User struct {
	ID              int64         `json:"id"`
	Email           string        `json:"email"`
	PasswordHash    string        `json:"-"`
	Role            Role          `json:"role"`
	IncomeBracket   IncomeBracket `json:"income_bracket"`
	SocialRole      SocialRole    `json:"social_role"`
	ReputationScore float64       `json:"reputation_score"`
	ProfileWeight   float64       `json:"profile_weight"`
	CreatedAt       time.Time     `json:"created_at"`
}
