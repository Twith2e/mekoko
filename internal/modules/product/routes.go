package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, authGuard gin.HandlerFunc, handler *Handler) {
	product := rg.Group("/product")

	{
		// product.POST("/add", handler.AddProducts)
		product.GET("", handler.GetProducts)
	}
}
