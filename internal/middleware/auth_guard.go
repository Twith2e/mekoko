package middleware

import (
	"database/sql"
	"errors"
	"log"
	appErr "mekoko/internal/errors"
	"mekoko/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	PublicIDContextKey  = "public_id"
	SessionIDContextKey = "session_id"
	UserRoleContextKey  = "user_role"
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
		log.Printf("token validation: err=%v, claims=%+v", err, claims)

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

		log.Printf("pid from auth guard: %s", sub)
		log.Printf("sid from auth guard: %s", sid)

		refreshToken, err := sessionChecker.IsSessionActive(c.Request.Context(), sid)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			mapped := response.MapError(err)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}
		if refreshToken == nil || errors.Is(err, sql.ErrNoRows) {
			mapped := response.MapError(appErr.ErrUnauthorized)
			c.AbortWithStatusJSON(mapped.Status, response.APIResponse[any]{
				Status: "error",
				Error:  &mapped.Error,
			})
			return
		}

		c.Set(PublicIDContextKey, sub)
		c.Set(SessionIDContextKey, sid)
		c.Set(UserRoleContextKey, refreshToken.Role)
		c.Next()
	}
}
