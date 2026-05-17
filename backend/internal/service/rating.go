package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/repository/postgres"
)

type RatingService struct {
	ratings   *postgres.RatingRepository
	practices *postgres.PracticeRepository
}

func NewRatingService(ratings *postgres.RatingRepository, practices *postgres.PracticeRepository) *RatingService {
	return &RatingService{ratings: ratings, practices: practices}
}

func (s *RatingService) Rate(ctx context.Context, rating *domain.Rating) (*domain.Rating, error) {
	result, err := s.ratings.Upsert(ctx, rating)
	if err != nil {
		return nil, err
	}
	if err := s.recomputePracticeRating(ctx, rating.PracticeID); err != nil {
		fmt.Printf("recompute rating for practice %d: %v\n", rating.PracticeID, err)
	}
	return result, nil
}

func (s *RatingService) GetMyRating(ctx context.Context, userID, practiceID int64) (*domain.Rating, error) {
	rating, err := s.ratings.GetByUserAndPractice(ctx, userID, practiceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return rating, nil
}

func (s *RatingService) GetStats(ctx context.Context, practiceID int64) (*domain.RatingStats, error) {
	stats, err := s.ratings.GetStats(ctx, practiceID)
	if err != nil {
		return nil, err
	}
	if stats.Count == 0 {
		return stats, nil
	}
	globalMean, _ := s.ratings.GetGlobalMean(ctx)
	const C = 10.0
	n := float64(stats.Count)
	stats.Rating = (n*stats.Rating + C*globalMean) / (n + C)
	return stats, nil
}

func (s *RatingService) recomputePracticeRating(ctx context.Context, practiceID int64) error {
	stats, err := s.ratings.GetStats(ctx, practiceID)
	if err != nil {
		return err
	}

	var bayesRating float64
	if stats.Count > 0 {
		globalMean, _ := s.ratings.GetGlobalMean(ctx)
		const C = 10.0
		n := float64(stats.Count)
		bayesRating = (n*stats.Rating + C*globalMean) / (n + C)
	}

	return s.practices.UpdateRating(ctx, practiceID, bayesRating, stats.Count)
}
