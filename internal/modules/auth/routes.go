package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler) {
	authGroup := r.Group("/auth")

	{
		authGroup.POST("/register", nil) // TODO: implement register handler
	}
}
