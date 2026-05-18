package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, authGuard gin.HandlerFunc, handler *Handler) {
	auth := r.Group("/auth")

	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.RefreshAccessToken)
		auth.POST("/password/forgot", handler.ForgotPassword)
		auth.PATCH("/password/reset", handler.ResetPassword)
	}

	{
		auth.POST("/logout", authGuard, handler.Logout)
		auth.PATCH("/password/change", authGuard, handler.ChangePassword)
	}
}
