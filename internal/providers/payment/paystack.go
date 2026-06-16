package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type Paystack struct {
	SecretKey string
	BaseURL   string
	Client    *http.Client
}

func NewPaystack(secretKey string, baseURL string) *Paystack {
	return &Paystack{
		SecretKey: secretKey,
		BaseURL:   baseURL,
		Client:    &http.Client{Timeout: time.Second * 15},
	}
}

func (p *Paystack) InitializeTransaction(amount int64, email string) {
	if p.SecretKey == "" {
		panic("Paystack secret key is required")
	}
	if p.BaseURL == "" {
		panic("Paystack base URL is required")
	}

	url := p.BaseURL + "/transaction/initialize"

	payload := InitializeTransactionRequest{
		Amount: amount,
		Email:  email,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.SecretKey)

	resp, err := p.Client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response: %d", resp.StatusCode)
		return
	}

	log.Printf("Transaction initialized successfully for email: %s", email)
}

func (p *Paystack) ValidateSignature(ctx context.Context, signature string, body []byte) bool {
	signature = strings.TrimSpace(signature)
	if signature == "" {
		return false
	}

	mac := hmac.New(sha512.New, []byte(p.SecretKey))
	_, err := mac.Write(body)
	if err != nil {
		return false
	}

	result := mac.Sum(nil)
	computed := hex.EncodeToString(result)

	return hmac.Equal([]byte(computed), []byte(signature))
}
