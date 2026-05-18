package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, authGuard gin.HandlerFunc, handler *Handler) {
	order := rg.Group("/order", authGuard)

	{
		order.POST("/create", handler.CreateOrder)
	}
}
