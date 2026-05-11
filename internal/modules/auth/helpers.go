package auth

import (
	appErr "mekoko/internal/errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pw string) (string, error) {
	hashedPasswordByte, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", appErr.ErrRegisteringUser
	}

	return string(hashedPasswordByte), nil
}

func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return appErr.ErrInvalidPasswordLength
	}

	var hasUpper, hasLower, hasSpecialChar, hasNum bool

	for _, r := range pw {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
		}
		if r >= 'a' && r <= 'z' {
			hasLower = true
		}
		if r >= 0 && r <= 9 {
			hasNum = true
		}
		if r == '!' || r == '@' || r == '?' || r == '$' || r == '%' || r == '#' || r == '^' || r == '*' {
			hasSpecialChar = true
		}

		if hasUpper && hasLower && hasNum && hasSpecialChar {
			return nil
		}
	}

	return appErr.ErrInvalidPassword
}

func ComparePassword(pwHash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(pw))
}
