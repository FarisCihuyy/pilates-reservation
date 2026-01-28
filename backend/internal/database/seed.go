package database

import (
	"log"
	"reservation-api/internal/models"

	"gorm.io/gorm"
)

// SeedData seeds initial data into database
func SeedData(db *gorm.DB) {
	log.Println("🌱 Seeding database...")

	// Check if data already exists
	var courtCount int64
	db.Model(&models.Court{}).Count(&courtCount)

	if courtCount > 0 {
		log.Println("ℹ️  Database already seeded, skipping...")
		return
	}

	// Seed courts
	seedCourts(db)

	// Seed timeslots
	seedTimeslots(db)

	log.Println("✅ Database seeding completed")
}

// seedCourts seeds court data
func seedCourts(db *gorm.DB) {
	courts := []models.Court{
		{
			Name:        "Studio A",
			Capacity:    10,
			Description: "Reformer Pilates - Premium equipment with personalized instruction",
			IsActive:    true,
		},
		{
			Name:        "Studio B",
			Capacity:    8,
			Description: "Mat Pilates - Classic exercises on comfortable mats",
			IsActive:    true,
		},
		{
			Name:        "Studio C",
			Capacity:    12,
			Description: "Mixed Class - Combination of Reformer and Mat exercises",
			IsActive:    true,
		},
	}

	for _, court := range courts {
		if err := db.Create(&court).Error; err != nil {
			log.Printf("⚠️  Failed to seed court: %s - %v", court.Name, err)
		} else {
			log.Printf("✓ Seeded court: %s", court.Name)
		}
	}
}

// seedTimeslots seeds timeslot data
func seedTimeslots(db *gorm.DB) {
	timeslots := []models.Timeslot{
		{Time: "08:00", Duration: 60, IsActive: true},
		{Time: "10:00", Duration: 60, IsActive: true},
		{Time: "12:00", Duration: 60, IsActive: true},
		{Time: "14:00", Duration: 60, IsActive: true},
		{Time: "16:00", Duration: 60, IsActive: true},
		{Time: "18:00", Duration: 60, IsActive: true},
		{Time: "20:00", Duration: 60, IsActive: true},
	}

	for _, timeslot := range timeslots {
		if err := db.Create(&timeslot).Error; err != nil {
			log.Printf("⚠️  Failed to seed timeslot: %s - %v", timeslot.Time, err)
		} else {
			log.Printf("✓ Seeded timeslot: %s", timeslot.Time)
		}
	}
}

// ClearDatabase clears all data from database (for testing)
func ClearDatabase(db *gorm.DB) error {
	log.Println("🗑️  Clearing database...")

	// Delete in reverse order of foreign keys
	if err := db.Exec("DELETE FROM payments").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM reservations").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM courts").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM timeslots").Error; err != nil {
		return err
	}

	log.Println("✅ Database cleared")
	return nil
}