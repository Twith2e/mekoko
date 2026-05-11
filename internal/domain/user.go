package domain

import "time"

type User struct {
	ID           int64
	UUID         string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string
	CreatedAt    time.Time
}
