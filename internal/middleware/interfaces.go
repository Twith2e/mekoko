package middleware

import (
	"context"
	"mekoko/internal/providers/tokens"
)

type Signer interface {
	ValidateAccessToken(tokenString string) (*tokens.AccessTokenClaims, error)
}

type SessionChecker interface {
	IsSessionActive(ctx context.Context, sid string) bool
}
