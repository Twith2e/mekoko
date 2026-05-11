package auth

import (
	"context"
	"database/sql"
	"log"
	"mekoko/internal/domain"
	appErr "mekoko/internal/errors"
	"time"
)

type DBTX interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

type Repository struct {
	db DBTX
}

func NewRepository(db DBTX) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	query := `
		INSERT INTO users (public_id, first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, public_id, first_name, last_name, email, created_at
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query,
		input.PublicID,
		input.FirstName,
		input.LastName,
		input.Email,
		input.PasswordHash,
	).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.ErrRegisteringUser
		}
		log.Printf("Failed to create user row: %s\n", err)
		return nil, appErr.ErrRegisteringUser
	}

	return &user, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, public_id, email, password_hash FROM users WHERE email = $1
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash)

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

func (r *Repository) StoreRefreshToken(ctx context.Context, userID int64, tokenHash, jti string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, jti, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, jti, expiresAt)
	return err
}
