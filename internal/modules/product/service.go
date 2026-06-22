package product

import (
	"context"
	"database/sql"
	"log"
	"mekoko/internal/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	db   *sql.DB
}

func NewService(repo *Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) AddProducts(ctx context.Context, payload []AddProductsRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	for _, product := range payload {
		productPublicID := uuid.NewString()
		newProduct := NewProduct{
			PublicID:           productPublicID,
			Name:               product.Name,
			Description:        product.Description,
			BasePrice:          product.BasePrice * 100,
			DiscountPercentage: product.DiscountPercentage,
			Slug:               Slugify(product.Name),
		}
		storedProduct, err := txRepo.AddProduct(ctx, newProduct)
		if err != nil {
			return err
		}

		for _, variant := range product.Variants {
			variantPublicID := uuid.NewString()
			newVariant := NewVariant{
				PublicID:      variantPublicID,
				ProductID:     storedProduct.ID,
				Color:         variant.Color,
				StockQuantity: variant.StockQuantity,
				Size:          variant.Size,
				ImageURL:      variant.ImageURL,
			}

			if err := txRepo.AddProductVariant(ctx, newVariant); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProductsWithFilter(ctx context.Context, limit, offset int, filter Filter) (products []domain.Product, count int64, err error) {
	products, count, err = s.repo.GetProductsWithFilter(ctx, limit, offset, filter)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, 0, err
	}
	return products, count, nil
}

func (s *Service) GetProducts(ctx context.Context, limit, offset int) (products []domain.Product, count int64, err error) {
	products, count, err = s.repo.GetProducts(ctx, limit, offset)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, 0, err
	}
	return products, count, nil
}

func (s *Service) GetProductByPublicID(ctx context.Context, publicID string) (*domain.Product, error) {
	product, err := s.repo.GetProductByPublicID(ctx, publicID)
	if err != nil {
		log.Printf("Error fetching product by public ID: %v", err)
		return nil, err
	}
	return product, nil
}
