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

const MaxRetry = 5

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

func (r *Repository) FetchWaitlistedEmails(ctx context.Context) ([]WaitlistEntry, error) {
	query := `SELECT email, created_at FROM waitlists`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []WaitlistEntry
	for rows.Next() {
		var entry WaitlistEntry
		if err := rows.Scan(&entry.Email, &entry.Date); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *Repository) GetWaitlistCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM waitlists`
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
