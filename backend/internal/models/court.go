package models

import (
	"gorm.io/gorm"
)

// Court represents a Pilates studio/court
type Court struct {
	gorm.Model
	Name         string        `json:"name" gorm:"not null"`
	Capacity     int           `json:"capacity" gorm:"not null"`
	Description  string        `json:"description"`
	IsActive     bool          `json:"is_active" gorm:"default:true"`
	Reservations []Reservation `json:"reservations,omitempty" gorm:"foreignKey:CourtID"`
}

// TableName specifies the table name for Court model
func (Court) TableName() string {
	return "courts"
}