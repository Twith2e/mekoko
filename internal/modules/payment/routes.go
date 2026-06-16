package payment

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, authGuard *gin.HandlerFunc, Handler *Handler) {
	payment := rg.Group("/payment")

	{
		payment.POST("/initialize", Handler.InitializeTransaction)
		payment.POST("/webhook", Handler.HandlePaymentWebhook)
	}
}
