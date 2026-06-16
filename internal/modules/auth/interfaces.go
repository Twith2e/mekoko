package auth

import (
	"context"
	"mekoko/internal/domain"
	"time"
)

type TokenGenerator interface {
	GenerateAccessToken(userID, sid, role string) (string, error)
	GenerateRefreshToken(userID, sid string) (string, string, time.Time, error)
}

type EmailSender interface {
	SendEmail(ctx context.Context, payload *domain.Email) (string, error)
}
