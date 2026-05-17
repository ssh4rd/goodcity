package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"goodcity/backend/internal/domain"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, c *domain.Comment) (*domain.Comment, error) {
	var result domain.Comment
	err := r.db.QueryRow(ctx,
		`INSERT INTO comments (practice_id, user_id, text)
		 VALUES ($1, $2, $3)
		 RETURNING id, practice_id, user_id, text, created_at`,
		c.PracticeID, c.UserID, c.Text,
	).Scan(&result.ID, &result.PracticeID, &result.UserID, &result.Text, &result.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return &result, nil
}

func (r *CommentRepository) ListByPractice(ctx context.Context, practiceID int64) ([]*domain.Comment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, practice_id, user_id, text, created_at
		 FROM comments WHERE practice_id=$1 ORDER BY created_at ASC`,
		practiceID,
	)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.PracticeID, &c.UserID, &c.Text, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, rows.Err()
}
