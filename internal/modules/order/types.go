package order

import "time"

type DeliveryStatus string

const (
	DeliveryStatusStitching DeliveryStatus = "stitching"
	DeliveryStatusYarnReady DeliveryStatus = "yarn_ready"
	DeliveryStatusInTransit DeliveryStatus = "in_transit"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
	DeliveryStatusUnraveled DeliveryStatus = "cancelled"
	DeliveryStatusTangled   DeliveryStatus = "tangled"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusSuccess PaymentStatus = "success"
)

type NewOrderDetails struct {
	ProductID          int64
	VariantID          int64
	VariantPublicID    string
	Quantity           int64
	BasePrice          int64
	DiscountPercentage int
	StockQuantity      int64
	ProductName        string
}

type CreateOrder struct {
	OrderPublicID  string
	UserID         int64
	Subtotal       int64
	TotalAmount    int64
	DeliveryFee    int64
	Currency       string
	DeliveryStatus DeliveryStatus
	PaymentStatus  PaymentStatus
	DiscountAmount int64
	OrderedAt      time.Time
}
