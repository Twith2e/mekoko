package domain

import "time"

type PasswordResetAttempt struct {
	ID             int64
	UserID         int64
	TokenHash      string
	TokenExpiresAt time.Time
	TokenUsedAt    time.Time
	CreatedAt      time.Time
}
