package order

type CreateOrderRequest struct {
	Order []Order `json:"order" binding:"required"`
}

type Order struct {
	VariantID string `json:"variant_id" binding:"required"`
	Quantity  int64  `json:"quantity" binding:"required,gt=0,lte=20"`
}

type CreateOrderResponse struct {
	OrderID string `json:"order_id"`
}
