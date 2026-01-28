package models

import (
	"time"

	"gorm.io/gorm"
)

// ReservationStatus defines the status of a reservation
type ReservationStatus string

const (
	StatusPending   ReservationStatus = "pending"
	StatusConfirmed ReservationStatus = "confirmed"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
)

// Reservation represents a booking made by a user
type Reservation struct {
	gorm.Model
	UserID     uint              `json:"user_id" gorm:"not null"`
	CourtID    uint              `json:"court_id" gorm:"not null"`
	TimeslotID uint              `json:"timeslot_id" gorm:"not null"`
	Date       time.Time         `json:"date" gorm:"not null;index"`
	Status     ReservationStatus `json:"status" gorm:"default:'pending'"`
	Notes      string            `json:"notes,omitempty"`

	// Relations
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Court    Court    `json:"court,omitempty" gorm:"foreignKey:CourtID"`
	Timeslot Timeslot `json:"timeslot,omitempty" gorm:"foreignKey:TimeslotID"`
	Payment  *Payment `json:"payment,omitempty" gorm:"foreignKey:ReservationID"`
}

// TableName specifies the table name for Reservation model
func (Reservation) TableName() string {
	return "reservations"
}

// IsConfirmed checks if reservation is confirmed
func (r *Reservation) IsConfirmed() bool {
	return r.Status == StatusConfirmed
}

// IsPending checks if reservation is pending
func (r *Reservation) IsPending() bool {
	return r.Status == StatusPending
}

// IsCancelled checks if reservation is cancelled
func (r *Reservation) IsCancelled() bool {
	return r.Status == StatusCancelled
}

// CanBeCancelled checks if reservation can be cancelled
func (r *Reservation) CanBeCancelled() bool {
	// Can only cancel if not already cancelled or completed
	return r.Status != StatusCancelled && r.Status != StatusCompleted
}