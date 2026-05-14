package auth

import (
	"mekoko/internal/domain"
	"time"
)

const CookieName = "mekoko_refresh_token"

type CreateUserInput struct {
	PublicID     string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
}

type Tokens struct {
	RefreshToken string
	AccessToken  string
	ExpiresAt    time.Time
}

type UserAndTokens struct {
	User   domain.User
	Tokens Tokens
}

type ResetEmailData struct {
	AppName       string
	ResetURL      string
	Name          string
	ExpiryMinutes int
	Year          int
	Email         string
}
