package cart

type AddToCartRequest struct {
	VariantID            string `json:"variant_id" binding:"required"`
	UnitPriceAtSelection int64  `json:"unit_price_at_selection" binding:"required"`
	Quantity             int64  `json:"quantity" binding:"required"`
}
