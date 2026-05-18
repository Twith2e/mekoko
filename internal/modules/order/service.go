package order

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mekoko/internal/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
	db   *sql.DB
}

func NewService(repo *Repository, db *sql.DB) *Service {
	return &Service{
		repo: repo,
		db:   db,
	}
}

func (s *Service) CreateOrder(ctx context.Context, userPublicID string, payload CreateOrderRequest) (*domain.Order, error) {
	user, err := s.repo.FindUserByPublicID(ctx, userPublicID)
	if err != nil {
		return nil, err
	}

	orderedItemsMap := make(map[string]*Order)
	orderedVariants := make([]string, 0, len(payload.Order))
	subtotal := int64(0)
	deliveryFee := int64(5000) // Assuming a fixed delivery fee for simplicity

	for _, p := range payload.Order {
		if orderedItemsMap[p.VariantID] == nil {
			orderedItemsMap[p.VariantID] = &Order{
				VariantID: p.VariantID,
				Quantity:  p.Quantity,
			}
			orderedVariants = append(orderedVariants, p.VariantID)

		} else {
			orderedItemsMap[p.VariantID].Quantity += p.Quantity
		}
	}

	products, err := s.repo.FetchProducts(ctx, orderedVariants)
	if err != nil {
		log.Printf("Error fetching products with variant id array, %s\n", err)
		return nil, err
	}

	for _, p := range products {
		if orderedItemsMap[p.VariantPublicID].Quantity > p.StockQuantity {
			return nil, fmt.Errorf("%s is low on stock", p.ProductName)
		}
		subtotal += p.BasePrice * orderedItemsMap[p.VariantPublicID].Quantity
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error starting transaction, %s\n", err)
		return nil, err
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTX(tx)

	order := CreateOrder{
		OrderPublicID:  uuid.NewString(),
		UserID:         user.ID,
		Subtotal:       subtotal,
		TotalAmount:    subtotal + deliveryFee, // Assuming no discounts or delivery fee for simplicity
		DeliveryFee:    deliveryFee,
		DeliveryStatus: DeliveryStatusStitching,
		PaymentStatus:  "pending",
		Currency:       "NGN",
		DiscountAmount: 0,
	}

	createdOrder, err := txRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	for _, p := range products {
		orderItemPublicID := uuid.NewString()
		if err := txRepo.CreateOrderItem(ctx, createdOrder.ID, p.BasePrice, p.BasePrice*orderedItemsMap[p.VariantPublicID].Quantity, p.VariantID, p.ProductID, orderedItemsMap[p.VariantPublicID].Quantity, orderItemPublicID, p.ProductName); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return createdOrder, nil
}
