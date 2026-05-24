package waitlist

import "time"

type WaitlistEntry struct {
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}
