package domain

import "time"

type CartItem struct {
	ID                   int64
	PublicID             string
	UserID               int64
	VariantID            int64
	Quantity             int64
	UnitPriceOnSelection int64
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
