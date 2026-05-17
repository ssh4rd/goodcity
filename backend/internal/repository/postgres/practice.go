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
		`INSERT INTO practices (title, description, city, category, address, latitude, longitude, budget_rub, implemented_at, status, author_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, title, description, city, category, address, latitude, longitude, budget_rub, implemented_at, status, author_id, rating, created_at`,
		p.Title, p.Description, p.City, p.Category, p.Address, p.Latitude, p.Longitude, p.BudgetRub, p.ImplementedAt, p.Status, p.AuthorID,
	).Scan(
		&result.ID, &result.Title, &result.Description, &result.City, &result.Category,
		&result.Address, &result.Latitude, &result.Longitude, &result.BudgetRub, &result.ImplementedAt,
		&result.Status, &result.AuthorID, &result.Rating, &result.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create practice: %w", err)
	}
	return &result, nil
}

func (r *PracticeRepository) GetByID(ctx context.Context, id int64) (*domain.Practice, error) {
	var p domain.Practice
	err := r.db.QueryRow(ctx,
		`SELECT p.id, p.title, p.description, p.city, p.category,
		        p.address, p.latitude, p.longitude, p.budget_rub, p.implemented_at,
		        p.status, p.author_id, p.rating, p.created_at,
		        COUNT(rt.id) AS vote_count
		 FROM practices p
		 LEFT JOIN ratings rt ON rt.practice_id = p.id AND rt.is_suspicious = FALSE
		 WHERE p.id = $1
		 GROUP BY p.id`,
		id,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.City, &p.Category,
		&p.Address, &p.Latitude, &p.Longitude, &p.BudgetRub, &p.ImplementedAt,
		&p.Status, &p.AuthorID, &p.Rating, &p.CreatedAt, &p.VoteCount,
	)
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
	if f.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(p.title ILIKE $%d OR p.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+f.Search+"%")
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
		`SELECT p.id, p.title, p.description, p.city, p.category,
		        p.address, p.latitude, p.longitude, p.budget_rub, p.implemented_at,
		        p.status, p.author_id, p.rating, p.created_at,
		        COUNT(rt.id) AS vote_count
		 FROM practices p
		 LEFT JOIN ratings rt ON rt.practice_id = p.id AND rt.is_suspicious = FALSE
		 %s
		 GROUP BY p.id
		 ORDER BY p.rating DESC, p.created_at DESC
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
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.City, &p.Category,
			&p.Address, &p.Latitude, &p.Longitude, &p.BudgetRub, &p.ImplementedAt,
			&p.Status, &p.AuthorID, &p.Rating, &p.CreatedAt, &p.VoteCount,
		); err != nil {
			return nil, 0, err
		}
		practices = append(practices, &p)
	}
	return practices, total, rows.Err()
}

func (r *PracticeRepository) Update(ctx context.Context, p *domain.Practice) (*domain.Practice, error) {
	var result domain.Practice
	err := r.db.QueryRow(ctx,
		`UPDATE practices SET title=$1, description=$2, city=$3, category=$4, status=$5,
		        address=$6, latitude=$7, longitude=$8, budget_rub=$9, implemented_at=$10
		 WHERE id=$11
		 RETURNING id, title, description, city, category, address, latitude, longitude, budget_rub, implemented_at, status, author_id, rating, created_at`,
		p.Title, p.Description, p.City, p.Category, p.Status,
		p.Address, p.Latitude, p.Longitude, p.BudgetRub, p.ImplementedAt, p.ID,
	).Scan(
		&result.ID, &result.Title, &result.Description, &result.City, &result.Category,
		&result.Address, &result.Latitude, &result.Longitude, &result.BudgetRub, &result.ImplementedAt,
		&result.Status, &result.AuthorID, &result.Rating, &result.CreatedAt,
	)
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

func (r *PracticeRepository) UpdateRating(ctx context.Context, id int64, rating float64, count int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE practices SET rating=$1, rating_count=$2 WHERE id=$3`,
		rating, count, id,
	)
	return err
}
