package service

import (
	"context"
	"errors"
	"fmt"

	"goodcity/backend/internal/domain"
	"goodcity/backend/internal/repository/postgres"
)

type PracticeService struct {
	practices *postgres.PracticeRepository
	votes     *postgres.VoteRepository
	comments  *postgres.CommentRepository
}

func NewPracticeService(
	practices *postgres.PracticeRepository,
	votes *postgres.VoteRepository,
	comments *postgres.CommentRepository,
) *PracticeService {
	return &PracticeService{practices: practices, votes: votes, comments: comments}
}

var ErrNotFound = errors.New("not found")

func (s *PracticeService) Create(ctx context.Context, p *domain.Practice) (*domain.Practice, error) {
	p.Status = domain.StatusPending
	return s.practices.Create(ctx, p)
}

func (s *PracticeService) GetByID(ctx context.Context, id int64) (*domain.Practice, error) {
	p, err := s.practices.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: practice %d", ErrNotFound, id)
	}
	return p, nil
}

func (s *PracticeService) List(ctx context.Context, f domain.PracticeFilter) ([]*domain.Practice, int, error) {
	return s.practices.List(ctx, f)
}

func (s *PracticeService) Update(ctx context.Context, p *domain.Practice) (*domain.Practice, error) {
	return s.practices.Update(ctx, p)
}

func (s *PracticeService) Delete(ctx context.Context, id int64) error {
	return s.practices.Delete(ctx, id)
}

func (s *PracticeService) Approve(ctx context.Context, id int64) error {
	return s.practices.UpdateStatus(ctx, id, domain.StatusApproved)
}

func (s *PracticeService) Reject(ctx context.Context, id int64) error {
	return s.practices.UpdateStatus(ctx, id, domain.StatusRejected)
}

func (s *PracticeService) ListPending(ctx context.Context) ([]*domain.Practice, int, error) {
	return s.practices.List(ctx, domain.PracticeFilter{Status: domain.StatusPending, Page: 1, PerPage: 100})
}

func (s *PracticeService) Vote(ctx context.Context, userID, practiceID int64) (bool, error) {
	return s.votes.Toggle(ctx, userID, practiceID)
}

func (s *PracticeService) AddComment(ctx context.Context, c *domain.Comment) (*domain.Comment, error) {
	return s.comments.Create(ctx, c)
}

func (s *PracticeService) ListComments(ctx context.Context, practiceID int64) ([]*domain.Comment, error) {
	return s.comments.ListByPractice(ctx, practiceID)
}
