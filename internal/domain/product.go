package domain

import "time"

type Product struct {
	ID                 int64
	PublicID           string
	Name               string
	DiscountPercentage int
	Description        string
	BasePrice          int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
