package middleware

import (
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString(UserRoleContextKey)
		if role == "" {
			mapped := response.MapError(appErr.ErrUnauthorized)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}

		if strings.ToLower(role) != "admin" {
			mapped := response.MapError(appErr.ErrForbidden)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}

		c.Next()
	}
}
