package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, authGuard gin.HandlerFunc, handler *Handler) {
	authGroup := r.Group("/auth")

	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/refresh", handler.RefreshAccessToken)
		authGroup.POST("/password/forgot", handler.ForgotPassword)
		authGroup.PATCH("/password/reset", handler.ResetPassword)
	}

	{
		authGroup.POST("/logout", authGuard, handler.Logout)
		authGroup.PATCH("/password/change", authGuard, handler.ChangePassword)
	}
}
