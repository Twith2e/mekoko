package cart

import (
	"context"
	"database/sql"
	"log"
	"mekoko/internal/domain"
	appErr "mekoko/internal/errors"
)

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

func WithTx(db *sql.Tx) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByPublicID(ctx context.Context, userPublic string) (*domain.User, error) {
	query := `
		SELECT id, public_id, email, password_hash 
		FROM users 
		WHERE public_id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userPublic).Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Could not find user by email: %s\n", err)
			return nil, appErr.ErrFindingUser
		}
		log.Printf("Failed to find user row: %s\n", err)
		return nil, appErr.ErrFindingUser
	}

	return &user, nil
}

func (r *Repository) FindProductVariantByPublicID(ctx context.Context, variantPublicID string) (*domain.ProductVariant, error) {
	query := `
		SELECT id, public_id, product_id, color, size, stock_quantity
		FROM product_variants
		WHERE public_id = $1
	`

	var productVariant domain.ProductVariant
	if err := r.db.QueryRowContext(ctx, query, variantPublicID).Scan(&productVariant.ID, &productVariant.PublicID, &productVariant.ProductID, &productVariant.Color, &productVariant.Size, &productVariant.StockQuantity); err != nil {
		return nil, err
	}

	return &productVariant, nil
}

func (r *Repository) AddToCart(ctx context.Context, cartPublic_id string, user_id, variant_id, unit_price_at_selection, quantity int64) error {
	query := `
		INSERT INTO cart_items (public_id, user_id, variant_id, unit_price_at_selection, quantity)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, variant_id)
		DO UPDATE SET
			quantity = cart_items.quantity + EXCLUDED.quantity,
			updated_at = NOW()
	`

	if _, err := r.db.ExecContext(ctx, query, cartPublic_id, user_id, variant_id, unit_price_at_selection, quantity); err != nil {
		return err
	}

	return nil
}

func (r *Repository) FetchAllCartItems(ctx context.Context, user_id int64) ([]CartForUI, error) {
	query := `
		SELECT cart_items.public_id, product_variants.public_id, cart_items.quantity, cart_items.unit_price_at_selection, product_variants.image_url, product_variants.color
		FROM cart_items
		JOIN product_variants ON cart_items.variant_id = product_variants.id
		WHERE user_id = $1 
	`

	var carts []CartForUI

	rows, err := r.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item CartForUI
		if err := rows.Scan(&item.ID, &item.VariantID, &item.Quantity, &item.UnitPriceAtSelection, &item.ImageURL, &item.Color); err != nil {
			return nil, err
		}
		carts = append(carts, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carts, nil
}
