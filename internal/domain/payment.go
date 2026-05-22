package domain

import "time"

type Payment struct {
	ID                int64
	PaymentPublicID   string
	OrderID           int64
	Provider          string
	ProviderReference string
	Amount            int64
	Currency          string
	Status            string
	PaidAt            time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
