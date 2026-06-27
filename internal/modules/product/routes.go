package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, authGuard, adminGuard gin.HandlerFunc, handler *Handler, adminHandler *AdminHandler) {
	product := rg.Group("/product")

	{
		product.GET("", handler.GetProducts)
		product.GET("/:public_id", handler.GetProductByPublicID)
		product.GET("/slug/:slug", handler.GetProductBySlug)
	}

	adminProduct := rg.Group("/admin/products", authGuard, adminGuard)

	{
		adminProduct.POST("/add", adminHandler.AddProducts)
		adminProduct.PATCH("/:public_id", adminHandler.UpdateProduct)
		adminProduct.DELETE("/:public_id", adminHandler.DeleteProduct)
		adminProduct.GET("", handler.GetProducts)
	}
}
