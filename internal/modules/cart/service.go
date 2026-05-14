package cart

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	db   *sql.DB
}

func NewService(repo *Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) AddToCart(ctx context.Context, userPublicID string, payload AddToCartRequest) error {
	user, err := s.repo.FindUserByPublicID(ctx, userPublicID)
	if err != nil {
		return err
	}

	productVariant, err := s.repo.FindProductVariantByPublicID(ctx, payload.VariantID)
	if err != nil {
		return err
	}

	cartPublicID := uuid.NewString()

	if err := s.repo.AddToCart(ctx, cartPublicID, user.ID, productVariant.ID, payload.UnitPriceAtSelection, payload.Quantity); err != nil {
		return err
	}
	return nil
}
