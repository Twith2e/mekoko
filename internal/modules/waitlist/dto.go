package waitlist

type WaitlistRequest struct {
	Email string `json:"email" binding:"required,email"`
}
