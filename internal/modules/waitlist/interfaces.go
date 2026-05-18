package waitlist

import "context"

type EmailSender interface {
	SendEmail(ctx context.Context, recipient, subject, htmlBody string) error
}
