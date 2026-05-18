package waitlist

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"
	"log"
	"strings"
	"time"
)

//go:embed templates/waitlist_confirmation.html
var confirmationTemplate string

type confirmationEmailData struct {
	AppName string
	Email   string
	Year    int
}

type Service struct {
	repo        *Repository
	emailSender EmailSender
	appName     string
}

func NewService(repo *Repository, emailSender EmailSender, appName string) *Service {
	return &Service{repo: repo, emailSender: emailSender, appName: appName}
}

func (s *Service) JoinWaitlist(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)

	if err := s.repo.AddToWaitlist(ctx, email); err != nil {
		return err
	}

	tmpl, err := template.New("confirmation").Parse(confirmationTemplate)
	if err != nil {
		log.Printf("waitlist: failed to parse confirmation template: %s", err)
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, confirmationEmailData{
		AppName: s.appName,
		Email:   email,
		Year:    time.Now().Year(),
	}); err != nil {
		log.Printf("waitlist: failed to render confirmation template: %s", err)
		return err
	}

	if err := s.emailSender.SendEmail(ctx, email, "You're on the list!", buf.String()); err != nil {
		log.Printf("waitlist: failed to send confirmation email to %s: %s", email, err)
	}

	return nil
}
