package models

import (
	"gorm.io/gorm"
)

// Timeslot represents a time slot for reservations
type Timeslot struct {
	gorm.Model
	Time         string        `json:"time" gorm:"not null"` // Format: "HH:MM"
	Duration     int           `json:"duration" gorm:"not null"` // Duration in minutes
	IsActive     bool          `json:"is_active" gorm:"default:true"`
	Reservations []Reservation `json:"reservations,omitempty" gorm:"foreignKey:TimeslotID"`
}

// TableName specifies the table name for Timeslot model
func (Timeslot) TableName() string {
	return "timeslots"
}