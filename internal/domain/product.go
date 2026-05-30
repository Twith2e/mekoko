package domain

import "time"

type Product struct {
	ID                 int64 `json:"id,omitempty"`
	PublicID           string
	Name               string
	DiscountPercentage int
	Description        string
	BasePrice          int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Variants           []ProductVariant
}
