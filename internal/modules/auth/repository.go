package auth

import (
	"context"
	"database/sql"
	"log"
	"mekoko/internal/domain"
	appErr "mekoko/internal/errors"
	"time"
)

const MaxRetry = 5

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

func (r *Repository) WithTx(tx *sql.Tx) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) CreateUser(ctx context.Context, input CreateUserInput) (*domain.User, error) {
	query := `
		INSERT INTO users (public_id, first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, public_id, first_name, last_name, email, role, created_at
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query,
		input.PublicID,
		input.FirstName,
		input.LastName,
		input.Email,
		input.PasswordHash,
		input.Role,
	).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.CreatedAt)

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
		SELECT id, public_id, email, password_hash, first_name, role
		FROM users
		WHERE email = $1
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(
			&user.ID,
			&user.UUID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.Role,
		)
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

func (r *Repository) FindUserByTokenHash(ctx context.Context, tokenHash string) (*domain.User, error) {
	query := `
		SELECT users.id, users.public_id, users.email
		FROM password_reset_attempts
		JOIN users ON users.id = password_reset_attempts.user_id
		WHERE token_hash = $1 AND token_expires_at > NOW() AND token_used_at IS NULL
	`
	var user domain.User
	if err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&user.ID, &user.UUID, &user.Email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, public_id, email, password_hash, role
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&user.ID,
			&user.UUID,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
		)
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

func (r *Repository) StoreRefreshToken(ctx context.Context, userID int64, sid, tokenHash, jti, role string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, sid, jti, role, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, sid, jti, role, expiresAt)
	return err
}

func (r *Repository) UpdateUserPasswordHash(ctx context.Context, pwHash string, userID int64) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, pwHash, userID)
	return err
}

func (r *Repository) IsSessionActive(ctx context.Context, sid string) (*domain.RefreshToken, error) {
	query := `
		SELECT role
		FROM refresh_tokens
		WHERE sid = $1 AND revoked_at IS NULL AND expires_at > NOW()
	`
	var refreshToken domain.RefreshToken
	if err := r.db.QueryRowContext(ctx, query, sid).Scan(&refreshToken.Role); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *Repository) FindRefreshTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT user_id, token_hash, sid, expires_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > NOW() AND revoked_at IS NULL
	`
	var refreshToken domain.RefreshToken
	if err := r.db.QueryRowContext(ctx, query, tokenHash).
		Scan(&refreshToken.UserID, &refreshToken.TokenHash, &refreshToken.SID, &refreshToken.ExpiresAt, &refreshToken.RevokedAt); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *Repository) RevokeCurrentSession(ctx context.Context, sid string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE sid = $1
	`

	_, err := r.db.ExecContext(ctx, query, sid)

	return err
}

func (r *Repository) RevokeAllSessions(ctx context.Context, userID int64) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1
	`

	_, err := r.db.ExecContext(ctx, query, userID)

	return err
}

func (r *Repository) RecordPasswordResetAttempt(ctx context.Context, tokenHash string, userID int64, token_expires_at time.Time) error {
	query := `
		INSERT INTO password_reset_attempts (user_id, token_hash, token_expires_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, userID, tokenHash, token_expires_at)
	return err
}

func (r *Repository) FindPasswordResetAttempt(ctx context.Context, tokenHash string, userID int64) (*domain.PasswordResetAttempt, error) {
	query := `
		SELECT user_id, token_hash, token_expires_at, token_used_at
		FROM password_reset_attempts
		WHERE token_hash = $1 AND user_id = $2 AND token_used_at IS NULL AND token_expires_at > NOW()
	`
	var passwordResetAttempt domain.PasswordResetAttempt
	if err := r.db.QueryRowContext(ctx, query, tokenHash, userID).Scan(&passwordResetAttempt.UserID, passwordResetAttempt.TokenHash, passwordResetAttempt.TokenExpiresAt, passwordResetAttempt.TokenUsedAt); err != nil {
		return nil, err
	}

	return &passwordResetAttempt, nil
}

func (r *Repository) MarkTokenUsed(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE password_reset_attempts
		SET token_used_at = NOW()
		WHERE token_hash = $1
	`

	_, err := r.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (r *Repository) GetPasswordResetAttemptCount(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT (count(*)) FROM password_reset_attempts
		WHERE user_id = $1 AND token_expires_at >= NOW() - INTERVAL '30 minutes'
	`
	var attemptCount int
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&attemptCount); err != nil {
		return 0, err
	}

	return attemptCount, nil
}

func (r *Repository) CreateOutgoingEmail(ctx context.Context, PublicID, subject, provider string, recipient int64, emailStruct domain.Email) (int64, error) {
	query := `
		INSERT INTO outgoing_emails (public_id, subject, provider, recipient, email_struct)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var outgoingEmailID int64

	err := r.db.QueryRowContext(ctx, query, PublicID, subject, provider, recipient, emailStruct).Scan(&outgoingEmailID)

	if err != nil {
		return 0, err
	}
	return outgoingEmailID, nil
}

func (r *Repository) UpdateOutgoingEmailOnResult(ctx context.Context, result, reasonForFailure, messageID string, ID int64) error {
	var query string
	if result == "failed" {
		query = `
			UPDATE outgoing_emails
			SET status = $1, reason_for_failure = $2, updated_at = NOW()
			WHERE id = $3
		`

		_, err := r.db.ExecContext(ctx, query, "failed", reasonForFailure, ID)

		return err

	} else {
		query = `
			UPDATE outgoing_emails
			SET status = $1, message_id = $2, delivered_at = NOW(), updated_at = NOW()
			WHERE id = $3
		`

		_, err := r.db.ExecContext(ctx, query, "successful", messageID, ID)

		return err
	}
}

func (r *Repository) FetchPendingOutgoingEmails(ctx context.Context) ([]domain.OutgoingEmail, error) {
	query := `
		SELECT id, public_id, subject, recipient, status, email_struct
		FROM outgoing_emails
		WHERE (status = 'pending' AND retry_count < $1) OR (status = 'failed' AND retry_count < $1)
	`

	outgoingEmails := make([]domain.OutgoingEmail, 0)

	rows, err := r.db.QueryContext(ctx, query, MaxRetry)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var outgoingEmail domain.OutgoingEmail
		if err := rows.Scan(&outgoingEmail.ID, &outgoingEmail.PublicID, &outgoingEmail.Subject, &outgoingEmail.Recipient, &outgoingEmail.Status, &outgoingEmail.EmailStruct); err != nil {
			return nil, err
		}

		outgoingEmails = append(outgoingEmails, outgoingEmail)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return outgoingEmails, nil
}

func (r *Repository) UpdateOutgoingEmailOnRetry(ctx context.Context, ID int64) error {
	query := `
		UPDATE outgoing_emails
		SET last_retry_at = NOW(), retry_count = outgoing_emails.retry_count + 1
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, ID)

	return err
}
