package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VoteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) *VoteRepository {
	return &VoteRepository{db: db}
}

// Toggle adds or removes a vote. Returns true if voted, false if unvoted.
func (r *VoteRepository) Toggle(ctx context.Context, userID, practiceID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM votes WHERE user_id=$1 AND practice_id=$2)`,
		userID, practiceID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check vote: %w", err)
	}

	if exists {
		_, err = r.db.Exec(ctx,
			`DELETE FROM votes WHERE user_id=$1 AND practice_id=$2`,
			userID, practiceID,
		)
		return false, err
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO votes (user_id, practice_id) VALUES ($1, $2)`,
		userID, practiceID,
	)
	return true, err
}

func (r *VoteRepository) HasVoted(ctx context.Context, userID, practiceID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM votes WHERE user_id=$1 AND practice_id=$2)`,
		userID, practiceID,
	).Scan(&exists)
	return exists, err
}
