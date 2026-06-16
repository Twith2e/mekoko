package payment

import (
	"database/sql"
	appErr "mekoko/internal/errors"
)

type Service struct {
	repo *Repository
	db   *sql.DB
}

func NewService(repo *Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) InitializeTransaction(payload InitializeTransactionRequest) (*InitializeTransactionResponse, error) {
	amount := payload.Amount
	// email := payload.Email

	if amount <= 0 {
		return nil, appErr.ErrInvalidAmount
	}

	return nil, nil
}

func (s *Service) ProcessWebhookEvent() {

}
