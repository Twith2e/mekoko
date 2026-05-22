package payment

type CreatePaymentRequest struct {
	OrderPublicID string `json:"order_public_id" validate:"required,uuid4"`
	Provider      string `json:"provider" validate:"required"`
}
