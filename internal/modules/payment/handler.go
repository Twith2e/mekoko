package payment

import "github.com/gin-gonic/gin"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitializeTransaction(c *gin.Context) {}

func (h *Handler) HandlePaymentWebhook(c *gin.Context) {}
