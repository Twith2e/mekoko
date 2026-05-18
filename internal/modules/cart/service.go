package cart

import (
	"context"
	"database/sql"
	"log"

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

	cartItemPublicID := uuid.NewString()

	if err := s.repo.AddToCart(ctx, cartItemPublicID, user.ID, productVariant.ID, payload.UnitPriceAtSelection, payload.Quantity); err != nil {
		return err
	}
	return nil
}

func (s *Service) FetchAllCartItems(ctx context.Context, userPublicID string) ([]CartForUI, error) {
	log.Printf("pid from service: %s", userPublicID)
	user, err := s.repo.FindUserByPublicID(ctx, userPublicID)
	if err != nil {
		return nil, err
	}

	cartItems, err := s.repo.FetchAllCartItems(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}
