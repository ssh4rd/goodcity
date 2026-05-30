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
	name, phone, city, district,
	COALESCE(income_bracket, 'middle'), COALESCE(social_role, 'resident'),
	COALESCE(reputation_score, 0), COALESCE(profile_weight, 1),
	COALESCE(is_role_verified, false), verification_doc_url,
	created_at`

func scanUser(row interface{ Scan(...any) error }, u *domain.User) error {
	return row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role,
		&u.Name, &u.Phone, &u.City, &u.District,
		&u.IncomeBracket, &u.SocialRole,
		&u.ReputationScore, &u.ProfileWeight,
		&u.IsRoleVerified, &u.VerificationDocURL,
		&u.CreatedAt,
	)
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string, role domain.Role, socialRole domain.SocialRole, name, city, district *string) (*domain.User, error) {
	var u domain.User
	err := scanUser(r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, role, social_role, name, city, district)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING `+userCols,
		email, passwordHash, role, socialRole, name, city, district,
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

func (r *UserRepository) Update(ctx context.Context, id int64, name, phone, city, district *string, socialRole domain.SocialRole, incomeBracket domain.IncomeBracket) (*domain.User, error) {
	var u domain.User
	err := scanUser(r.db.QueryRow(ctx,
		`UPDATE users SET name=$2, phone=$3, city=$4, district=$5, social_role=$6, income_bracket=$7
		 WHERE id=$1 RETURNING `+userCols,
		id, name, phone, city, district, socialRole, incomeBracket,
	), &u)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) SetVerificationDoc(ctx context.Context, userID int64, docURL string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET verification_doc_url = $1, is_role_verified = false WHERE id = $2`,
		docURL, userID,
	)
	return err
}

func (r *UserRepository) SetRoleVerified(ctx context.Context, userID int64, verified bool) error {
	var err error
	if verified {
		_, err = r.db.Exec(ctx,
			`UPDATE users SET is_role_verified = true WHERE id = $1`,
			userID,
		)
	} else {
		_, err = r.db.Exec(ctx,
			`UPDATE users SET is_role_verified = false, verification_doc_url = NULL WHERE id = $1`,
			userID,
		)
	}
	return err
}

func (r *UserRepository) ListPendingVerifications(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+userCols+` FROM users
		 WHERE verification_doc_url IS NOT NULL
		   AND is_role_verified = false
		   AND social_role != 'resident'
		 ORDER BY created_at`,
	)
	if err != nil {
		return nil, fmt.Errorf("list pending verifications: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := scanUser(rows, &u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}
