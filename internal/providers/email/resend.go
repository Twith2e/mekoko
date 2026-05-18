package email

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/resend/resend-go/v3"
)

type Resend struct {
	ApiKey string
	Client *http.Client
}

func NewResend(apiKey string) *Resend {
	return &Resend{ApiKey: apiKey, Client: &http.Client{Timeout: time.Second * 15}}
}

func (r *Resend) SendEmail(ctx context.Context, recipient, subject string, htmlBody string) error {
	client := resend.NewClient(r.ApiKey)
	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{recipient},
		Html:    htmlBody,
		Subject: subject,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sent.Id)
	return nil
}
