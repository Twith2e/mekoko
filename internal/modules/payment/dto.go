package payment

type InitializeTransactionRequest struct {
	OrderPublicID  string `json:"order_public_id" validate:"required,uuid4"`
	Amount         int64  `json:"amount" validate:"required,gt=0"`
	Email          string `json:"email" validate:"required,email"`
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`
}

type InitializeTransactionResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}
