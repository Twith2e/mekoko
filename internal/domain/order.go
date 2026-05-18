package domain

import (
	"time"
)

type DeliveryStatus string

const (
	DeliveryStatusStitching DeliveryStatus = "stitching"
	DeliveryStatusYarnReady DeliveryStatus = "yarn_ready"
	DeliveryStatusInTransit DeliveryStatus = "in_transit"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusUnraveled DeliveryStatus = "cancelled"
	DeliveryStatusTangled   DeliveryStatus = "tangled"
)

type Order struct {
	ID             int64
	PublicID       string
	UserID         int64
	Subtotal       int64
	TotalAmount    int64
	DeliveryFee    int64
	DeliveryStatus DeliveryStatus
	PaymentStatus  string
	Currency       string
	DiscountAmount int64
	OrderedAt      *time.Time
	DeliveredAt    *time.Time
}
