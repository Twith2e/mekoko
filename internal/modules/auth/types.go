package auth

import "mekoko/internal/domain"

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
}

type UserAndTokens struct {
	User   domain.User
	Tokens Tokens
}
