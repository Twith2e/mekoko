package email

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"mekoko/internal/modules/auth"
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

//go:embed templates/password_reset.html
var passwordReset string

func (r *Resend) SendEmail(ctx context.Context, recipient, subject string, data auth.ResetEmailData) error {
	client := resend.NewClient(r.ApiKey)

	tmpl, err := template.New("reset").Parse(passwordReset)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("%s", err)
		return err
	}

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{recipient},
		Html:    buf.String(),
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
