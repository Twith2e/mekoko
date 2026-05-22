package waitlist

import (
	"context"
	"mekoko/internal/domain"
)

type EmailSender interface {
	SendEmail(ctx context.Context, payload *domain.Email) (string, error)
}
