package email

import (
	"context"
	"fmt"
	"io"
	"log"
	"mekoko/internal/domain"

	sendinblue "github.com/sendinblue/APIv3-go-library/v2/lib"
)

type Brevo struct {
	Sib *sendinblue.APIClient
}

func NewBrevo(apiKey string) *Brevo {
	log.Printf("brevo api key loaded: %v", apiKey != "")
	cfg := sendinblue.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	cfg.AddDefaultHeader("partner-key", apiKey)
	return &Brevo{
		Sib: sendinblue.NewAPIClient(cfg),
	}
}

func (b *Brevo) SendEmail(ctx context.Context, payload *domain.Email) (string, error) {
	sendSmtpEmailSender := sendinblue.SendSmtpEmailSender{
		Name:  payload.SenderName,
		Email: payload.Sender,
	}

	sendSmtpEmailTo := []sendinblue.SendSmtpEmailTo{}

	recipient := sendinblue.SendSmtpEmailTo{
		Name:  payload.RecipientName,
		Email: payload.Recipient,
	}

	sendSmtpEmailTo = append(sendSmtpEmailTo, recipient)

	sendSmtpEmail := sendinblue.SendSmtpEmail{
		Sender:      &sendSmtpEmailSender,
		To:          sendSmtpEmailTo,
		HtmlContent: payload.HtmlContent,
		Subject:     payload.Subject,
	}

	createEmail, resp, err := b.Sib.TransactionalEmailsApi.SendTransacEmail(ctx, sendSmtpEmail)

	if err != nil {
		if apiErr, ok := err.(sendinblue.GenericSwaggerError); ok {
			log.Printf("BREVO error body: %s", string(apiErr.Body()))
		}
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("BREVO error: %s", string(body))

		return "", fmt.Errorf("failed to send email with Brevo")
	}

	return createEmail.MessageId, nil
}
