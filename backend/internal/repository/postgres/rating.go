package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"goodcity/backend/internal/domain"
)

type RatingRepository struct {
	db *pgxpool.Pool
}

func NewRatingRepository(db *pgxpool.Pool) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) Upsert(ctx context.Context, rating *domain.Rating) (*domain.Rating, error) {
	var result domain.Rating
	err := r.db.QueryRow(ctx,
		`INSERT INTO ratings (user_id, practice_id, score_utility, score_aesthetics, score_ecology, comment, user_lat, user_lon)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (user_id, practice_id) DO UPDATE SET
		     score_utility    = EXCLUDED.score_utility,
		     score_aesthetics = EXCLUDED.score_aesthetics,
		     score_ecology    = EXCLUDED.score_ecology,
		     comment          = EXCLUDED.comment,
		     user_lat         = EXCLUDED.user_lat,
		     user_lon         = EXCLUDED.user_lon,
		     created_at       = NOW()
		 RETURNING id, user_id, practice_id, score_utility, score_aesthetics, score_ecology,
		           COALESCE(comment, ''), user_lat, user_lon, is_suspicious, created_at`,
		rating.UserID, rating.PracticeID,
		rating.ScoreUtility, rating.ScoreAesthetics, rating.ScoreEcology,
		rating.Comment, rating.UserLat, rating.UserLon,
	).Scan(
		&result.ID, &result.UserID, &result.PracticeID,
		&result.ScoreUtility, &result.ScoreAesthetics, &result.ScoreEcology,
		&result.Comment, &result.UserLat, &result.UserLon,
		&result.IsSuspicious, &result.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("upsert rating: %w", err)
	}
	return &result, nil
}

func (r *RatingRepository) GetByUserAndPractice(ctx context.Context, userID, practiceID int64) (*domain.Rating, error) {
	var result domain.Rating
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, practice_id, score_utility, score_aesthetics, score_ecology,
		        COALESCE(comment, ''), user_lat, user_lon, is_suspicious, created_at
		 FROM ratings WHERE user_id=$1 AND practice_id=$2`,
		userID, practiceID,
	).Scan(
		&result.ID, &result.UserID, &result.PracticeID,
		&result.ScoreUtility, &result.ScoreAesthetics, &result.ScoreEcology,
		&result.Comment, &result.UserLat, &result.UserLon,
		&result.IsSuspicious, &result.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get rating: %w", err)
	}
	return &result, nil
}

func (r *RatingRepository) GetStats(ctx context.Context, practiceID int64) (*domain.RatingStats, error) {
	var s domain.RatingStats
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*),
		        COALESCE(AVG(score_utility), 0),
		        COALESCE(AVG(score_aesthetics), 0),
		        COALESCE(AVG(score_ecology), 0)
		 FROM ratings WHERE practice_id=$1 AND is_suspicious=FALSE`,
		practiceID,
	).Scan(&s.Count, &s.AvgUtility, &s.AvgAesthetics, &s.AvgEcology)
	if err != nil {
		return nil, fmt.Errorf("get rating stats: %w", err)
	}
	s.Rating = 0.4*s.AvgUtility + 0.35*s.AvgAesthetics + 0.25*s.AvgEcology
	return &s, nil
}

func (r *RatingRepository) GetGlobalMean(ctx context.Context) (float64, error) {
	var mean float64
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(AVG(0.4*score_utility + 0.35*score_aesthetics + 0.25*score_ecology), 3.0)
		 FROM ratings WHERE is_suspicious=FALSE`,
	).Scan(&mean)
	if err != nil {
		return 3.0, fmt.Errorf("get global mean: %w", err)
	}
	return mean, nil
}
