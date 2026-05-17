package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"goodcity/backend/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

const userCols = `id, email, password_hash, role,
	COALESCE(income_bracket, 'middle'), COALESCE(social_role, 'resident'),
	COALESCE(reputation_score, 0), COALESCE(profile_weight, 1), created_at`

func scanUser(row interface{ Scan(...any) error }, u *domain.User) error {
	return row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role,
		&u.IncomeBracket, &u.SocialRole,
		&u.ReputationScore, &u.ProfileWeight, &u.CreatedAt,
	)
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string, role domain.Role, socialRole domain.SocialRole) (*domain.User, error) {
	var u domain.User
	err := scanUser(r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, role, social_role) VALUES ($1, $2, $3, $4)
		 RETURNING `+userCols,
		email, passwordHash, role, socialRole,
	), &u)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := scanUser(r.db.QueryRow(ctx,
		`SELECT `+userCols+` FROM users WHERE email = $1`,
		email,
	), &u)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	err := scanUser(r.db.QueryRow(ctx,
		`SELECT `+userCols+` FROM users WHERE id = $1`,
		id,
	), &u)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}
