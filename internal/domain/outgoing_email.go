package domain

import "time"

type OutgoingEmail struct {
	ID               int64
	PublicID         string
	MessageID        string
	Subject          string
	Recipient        int64
	ReasonForFailure string
	Status           string
	EmailStruct      interface{}
	LastRetryAt      time.Time
	DeliveredAt      time.Time
	RetryCount       int
}
