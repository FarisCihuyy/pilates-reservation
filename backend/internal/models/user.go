package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Name         string        `json:"name" gorm:"not null"`
	Email        string        `json:"email" gorm:"uniqueIndex;not null"`
	Password     string        `json:"-" gorm:"not null"` // Never expose password in JSON
	Phone        string        `json:"phone"`
	IsActive     bool          `json:"is_active" gorm:"default:true"`
	Reservations []Reservation `json:"reservations,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Any pre-creation logic can go here
	// For example, email validation, default values, etc.
	return nil
}