package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, authGuard, adminGuard gin.HandlerFunc, handler, adminHandler *Handler) {
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

	adminAuth := r.Group("/admin/auth")
	{
		adminAuth.POST("/register", adminGuard, adminHandler.Register)
		adminAuth.POST("/login", adminHandler.Login)
		adminAuth.POST("/refresh", adminHandler.RefreshAccessToken)
		adminAuth.POST("/logout", adminGuard, adminHandler.Logout)
	}
}
