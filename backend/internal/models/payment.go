package models

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus defines the status of a payment
type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
	PaymentExpired PaymentStatus = "expired"
)

// Payment represents a payment transaction
type Payment struct {
	gorm.Model
	ReservationID uint          `json:"reservation_id" gorm:"not null;uniqueIndex"`
	Amount        float64       `json:"amount" gorm:"not null"`
	Status        PaymentStatus `json:"status" gorm:"default:'pending'"`
	PaymentMethod string        `json:"payment_method,omitempty"`
	TransactionID string        `json:"transaction_id" gorm:"uniqueIndex"`

	// Midtrans specific fields
	MidtransToken string `json:"midtrans_token,omitempty"`
	MidtransURL   string `json:"midtrans_url,omitempty"`

	// Payment timestamps
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`

	// Relation
	Reservation Reservation `json:"reservation,omitempty" gorm:"foreignKey:ReservationID"`
}

// TableName specifies the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}

// IsPaid checks if payment is completed
func (p *Payment) IsPaid() bool {
	return p.Status == PaymentPaid
}

// IsPending checks if payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == PaymentPending
}

// IsFailed checks if payment has failed
func (p *Payment) IsFailed() bool {
	return p.Status == PaymentFailed
}

// IsExpired checks if payment has expired
func (p *Payment) IsExpired() bool {
	if p.Status == PaymentExpired {
		return true
	}
	if p.ExpiredAt != nil && time.Now().After(*p.ExpiredAt) {
		return true
	}
	return false
}