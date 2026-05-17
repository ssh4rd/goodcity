package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"goodcity/backend/internal/domain"
)

type PracticeRepository struct {
	db *pgxpool.Pool
}

func NewPracticeRepository(db *pgxpool.Pool) *PracticeRepository {
	return &PracticeRepository{db: db}
}

func (r *PracticeRepository) Create(ctx context.Context, p *domain.Practice) (*domain.Practice, error) {
	var result domain.Practice
	err := r.db.QueryRow(ctx,
		`INSERT INTO practices (title, description, city, category, status, author_id)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, title, description, city, category, status, author_id, created_at`,
		p.Title, p.Description, p.City, p.Category, p.Status, p.AuthorID,
	).Scan(&result.ID, &result.Title, &result.Description, &result.City,
		&result.Category, &result.Status, &result.AuthorID, &result.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create practice: %w", err)
	}
	return &result, nil
}

func (r *PracticeRepository) GetByID(ctx context.Context, id int64) (*domain.Practice, error) {
	var p domain.Practice
	err := r.db.QueryRow(ctx,
		`SELECT p.id, p.title, p.description, p.city, p.category, p.status, p.author_id, p.created_at,
		        COUNT(v.user_id) as vote_count
		 FROM practices p
		 LEFT JOIN votes v ON v.practice_id = p.id
		 WHERE p.id = $1
		 GROUP BY p.id`,
		id,
	).Scan(&p.ID, &p.Title, &p.Description, &p.City, &p.Category,
		&p.Status, &p.AuthorID, &p.CreatedAt, &p.VoteCount)
	if err != nil {
		return nil, fmt.Errorf("get practice by id: %w", err)
	}
	return &p, nil
}

func (r *PracticeRepository) List(ctx context.Context, f domain.PracticeFilter) ([]*domain.Practice, int, error) {
	conditions := []string{}
	args := []any{}
	argIdx := 1

	if f.City != "" {
		conditions = append(conditions, fmt.Sprintf("p.city = $%d", argIdx))
		args = append(args, f.City)
		argIdx++
	}
	if f.Category != "" {
		conditions = append(conditions, fmt.Sprintf("p.category = $%d", argIdx))
		args = append(args, f.Category)
		argIdx++
	}
	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIdx))
		args = append(args, f.Status)
		argIdx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.db.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(*) FROM practices p %s`, where), args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count practices: %w", err)
	}

	perPage := f.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	page := f.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage

	args = append(args, perPage, offset)
	rows, err := r.db.Query(ctx, fmt.Sprintf(
		`SELECT p.id, p.title, p.description, p.city, p.category, p.status, p.author_id, p.created_at,
		        COUNT(v.user_id) as vote_count
		 FROM practices p
		 LEFT JOIN votes v ON v.practice_id = p.id
		 %s
		 GROUP BY p.id
		 ORDER BY p.created_at DESC
		 LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1),
		args...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list practices: %w", err)
	}
	defer rows.Close()

	var practices []*domain.Practice
	for rows.Next() {
		var p domain.Practice
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.City, &p.Category,
			&p.Status, &p.AuthorID, &p.CreatedAt, &p.VoteCount); err != nil {
			return nil, 0, err
		}
		practices = append(practices, &p)
	}
	return practices, total, rows.Err()
}

func (r *PracticeRepository) Update(ctx context.Context, p *domain.Practice) (*domain.Practice, error) {
	var result domain.Practice
	err := r.db.QueryRow(ctx,
		`UPDATE practices SET title=$1, description=$2, city=$3, category=$4, status=$5
		 WHERE id=$6
		 RETURNING id, title, description, city, category, status, author_id, created_at`,
		p.Title, p.Description, p.City, p.Category, p.Status, p.ID,
	).Scan(&result.ID, &result.Title, &result.Description, &result.City,
		&result.Category, &result.Status, &result.AuthorID, &result.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("update practice: %w", err)
	}
	return &result, nil
}

func (r *PracticeRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM practices WHERE id=$1`, id)
	return err
}

func (r *PracticeRepository) UpdateStatus(ctx context.Context, id int64, status domain.PracticeStatus) error {
	_, err := r.db.Exec(ctx, `UPDATE practices SET status=$1 WHERE id=$2`, status, id)
	return err
}
