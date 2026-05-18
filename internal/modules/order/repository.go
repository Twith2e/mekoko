package order

import (
	"context"
	"database/sql"
	"log"
	"mekoko/internal/domain"
	appErr "mekoko/internal/errors"
)

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTX(db *sql.Tx) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByPublicID(ctx context.Context, publicID string) (*domain.User, error) {
	query := `
		SELECT id, public_id, email, password_hash 
		FROM users 
		WHERE public_id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash)
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

func (r *Repository) FetchProducts(ctx context.Context, variantIDs []string) ([]NewOrderDetails, error) {
	query := `
		SELECT p.id, p.public_id, p.name, p.description, p.discount_percentage, p.base_price,
		       v.id, v.public_id, v.product_id, v.color, v.size, v.image_url, v.stock_quantity
		FROM products p
		JOIN product_variants v ON p.id = v.product_id
		WHERE v.public_id = ANY($1)
	`
	rows, err := r.db.QueryContext(ctx, query, variantIDs)
	if err != nil {
		return nil, err
	}

	orderProducts := make([]NewOrderDetails, 0)
	for rows.Next() {
		var o NewOrderDetails
		if err := rows.Scan(&o.ProductID, &o.ProductName, &o.DiscountPercentage, &o.BasePrice, &o.VariantID, &o.VariantPublicID, &o.StockQuantity); err != nil {
			return nil, err
		}
		orderProducts = append(orderProducts, o)
	}

	return orderProducts, nil
}

func (r *Repository) CreateOrder(ctx context.Context, order CreateOrder) (*domain.Order, error) {
	query := `
		INSERT INTO orders (public_id, user_id, subtotal, total_amount, delivery_fee, delivery_status, payment_status, currency, discount_amount)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, public_id, user_id, subtotal, total_amount, delivery_fee, delivery_status, payment_status, currency, discount_amount
	`
	var createdOrder domain.Order

	err := r.db.QueryRowContext(ctx, query, order.OrderPublicID, order.UserID, order.Subtotal, order.TotalAmount, order.DeliveryFee, order.DeliveryStatus, order.PaymentStatus, order.Currency, order.DiscountAmount).Scan(&createdOrder.ID, &createdOrder.PublicID, &createdOrder.UserID, &createdOrder.Subtotal, &createdOrder.TotalAmount, &createdOrder.DeliveryFee, &createdOrder.DeliveryStatus, &createdOrder.PaymentStatus, &createdOrder.Currency, &createdOrder.DiscountAmount)
	if err != nil {
		return nil, err
	}

	return &createdOrder, nil
}

func (r *Repository) CreateOrderItem(ctx context.Context, orderID int64, unitPrice, totalPrice, variantID, productID, quantity int64, orderItemPublicID, productName string) error {
	query := `
		INSERT INTO order_items (public_id, order_id, product_id, variant_id, quantity, unit_price, total_price, product_name_snapshot)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if _, err := r.db.ExecContext(ctx, query, orderItemPublicID, orderID, productID, variantID, quantity, unitPrice, totalPrice, productName); err != nil {
		return err
	}

	return nil
}
