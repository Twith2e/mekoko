package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, authGuard, adminGuard gin.HandlerFunc, handler *Handler, adminHandler *AdminHandler) {
	product := rg.Group("/product")

	{
		product.GET("", handler.GetProducts)
		product.GET("/:public_id", handler.GetProductByPublicID)
	}

	adminProduct := rg.Group("/admin/product", authGuard, adminGuard)

	{
		adminProduct.POST("/add", adminHandler.AddProducts)
	}
}
