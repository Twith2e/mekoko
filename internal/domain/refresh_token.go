package domain

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	JTI       string
	SID       string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
