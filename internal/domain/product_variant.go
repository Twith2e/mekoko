package domain

import "time"

type ProductVariant struct {
	ID            int64 `json:"id,omitempty"`
	PublicID      string
	ProductID     int64
	Color         string
	Size          string
	ImageURL      string
	StockQuantity int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
