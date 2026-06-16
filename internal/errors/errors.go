package errors

import "errors"

var (
	ErrInvalidRequestBody        = errors.New("Invalid request body")
	ErrInvalidPasswordLength     = errors.New("Password must be at least 8 characters long")
	ErrInvalidPassword           = errors.New("Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	ErrPasswordMismatch          = errors.New("Password and confirm password do not match")
	ErrRegisteringUser           = errors.New("Failed to register user")
	ErrFindingUser               = errors.New("User not found")
	ErrInvalidCredentials        = errors.New("Invalid credentials")
	ErrUserExists                = errors.New("Email is taken")
	ErrUnauthorized              = errors.New("Unauthorized")
	ErrForbidden                 = errors.New("Forbidden")
	ErrInvalidSession            = errors.New("Invalid or expired session")
	ErrRefreshingAccessToken     = errors.New("Failed to refresh access token")
	ErrTooManyRequests           = errors.New("Too many requests")
	ErrPasswordResetUserNotFound = errors.New("If that email exists, you'll receuve a reset link")
	ErrInvalidToken              = errors.New("Invalid or expired tokens")
	ErrInvalidRequestQuery       = errors.New("Invalid request query")
	ErrInvalidPriceRange         = errors.New("Max price must be greater than min price")
	ErrInvalidAmount             = errors.New("Amount must be greater than zero")
)
