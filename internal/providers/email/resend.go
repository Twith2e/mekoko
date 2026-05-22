package email

import (
	"context"
	"fmt"
	"mekoko/internal/domain"
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

func (r *Resend) SendEmail(ctx context.Context, payload *domain.Email) (string, error) {
	client := resend.NewClient(r.ApiKey)
	params := &resend.SendEmailRequest{
		From:    payload.Sender,
		To:      []string{payload.Recipient},
		Html:    payload.HtmlContent,
		Subject: payload.Subject,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Println(sent.Id)
	return sent.Id, nil
}
