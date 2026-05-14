package auth

import (
	"context"
	"time"
)

type TokenGenerator interface {
	GenerateAccessToken(userID, sid string) (string, error)
	GenerateRefreshToken(userID, sid string) (string, string, time.Time, error)
}

type EmailSender interface {
	SendEmail(ctx context.Context, recipient, subject, url string) error
}
