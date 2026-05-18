package waitlist

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddToWaitlist(ctx context.Context, email string) error {
	query := `
		INSERT INTO waitlists (email)
		VALUES ($1)
		ON CONFLICT DO NOTHING
	`

	if _, err := r.db.ExecContext(ctx, query, email); err != nil {
		return err
	}

	return nil
}
