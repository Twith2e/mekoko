package email

import (
	"context"
	_ "embed"
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

//go:embed templates/password_change.html
var test string

func (r *Resend) SendEmail(ctx context.Context, recipient, subject, url string) error {
	client := resend.NewClient(r.ApiKey)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{recipient},
		Html:    test,
		Subject: subject,
	}

	fmt.Println(test)

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(sent.Id)
	return nil
}
