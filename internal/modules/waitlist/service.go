package waitlist

import (
	"context"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) JoinWaitlist(ctx context.Context, email string) error {
	return s.repo.AddToWaitlist(ctx, strings.TrimSpace(email))
}
