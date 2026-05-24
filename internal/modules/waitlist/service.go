package waitlist

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"
	"log"
	"mekoko/internal/domain"
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
	sender      string
}

func NewService(repo *Repository, emailSender EmailSender, appName, sender string) *Service {
	return &Service{repo: repo, emailSender: emailSender, appName: appName, sender: sender}
}

func (s *Service) JoinWaitlist(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)

	if err := s.repo.AddToWaitlist(ctx, email); err != nil {
		log.Printf("waitlist: failed to add %s to waitlist: %s", email, err)
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

	subject := "You're on the list!"

	payload := domain.Email{
		Recipient:     email,
		RecipientName: "",
		Sender:        s.sender,
		SenderName:    s.appName,
		HtmlContent:   buf.String(),
		Subject:       subject,
	}

	log.Printf("waitlist service, sender: %s", s.sender)

	if _, err = s.emailSender.SendEmail(ctx, &payload); err != nil {
		log.Printf("waitlist: failed to send confirmation email to %s: %s", email, err)
	}

	return nil
}

func (s *Service) FetchWaitlistedEmails(ctx context.Context) ([]WaitlistEntry, error) {
	return s.repo.FetchWaitlistedEmails(ctx)
}

func (s *Service) GetWaitlistCount(ctx context.Context) (int, error) {
	return s.repo.GetWaitlistCount(ctx)
}
