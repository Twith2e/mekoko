package payment

import (
	"context"
	"database/sql"
	"mekoko/internal/domain"
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

func (r *Repository) WithTx(db *sql.Tx) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByPublicID(ctx context.Context, userPublicID string) (*domain.User, error) {
	query := `
		SELECT id, public_id, email, password_hash
		FROM users
		WHERE public_id = $1
	`
	var user domain.User

	if err := r.db.QueryRowContext(ctx, query, userPublicID).Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindOrderByPublicID(ctx context.Context, orderPublicID string) (*domain.Order, error) {
	query := `
		SELECT id, public_id, user_id, total_amount, payment_status, delivery_status, ordered_at, delivered_at
		FROM orders
		WHERE public_id = $1
	`
	var order domain.Order

	if err := r.db.QueryRowContext(ctx, query, orderPublicID).Scan(&order.ID, &order.PublicID, &order.UserID, &order.TotalAmount, &order.PaymentStatus, &order.DeliveryStatus, &order.OrderedAt, &order.DeliveredAt); err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) CreatePayment(ctx context.Context, userID int64, amount int64, orderID int64, providerReference, provider string) error {
	query := `
		INSERT INTO payments (user_id, amount, order_id, provider_reference, provider)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query, userID, amount, orderID, providerReference, provider)
	return err
}
