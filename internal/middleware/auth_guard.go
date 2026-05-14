package middleware

import (
	"log"
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	PublicIDContextKey  = "public_id"
	SessionIDContextKey = "session_id"
)

func AuthGuard(signer Signer, sessionChecker SessionChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			mapped := response.MapError(appErr.ErrUnauthorized)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}

		token := strings.TrimSpace(parts[1])
		claims, err := signer.ValidateAccessToken(token)
		if token == "" || err != nil {
			mapped := response.MapError(appErr.ErrUnauthorized)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}

		sub := claims.Subject
		sid := claims.SID

		log.Printf("sid from auth guard: %s", sid)

		if !sessionChecker.IsSessionActive(c.Request.Context(), sid) {
			mapped := response.MapError(appErr.ErrUnauthorized)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}

		c.Set(PublicIDContextKey, sub)
		c.Set(SessionIDContextKey, sid)
		c.Next()
	}
}
