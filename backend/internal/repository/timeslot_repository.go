package repository

import (
	"reservation-api/internal/models"

	"gorm.io/gorm"
)

// TimeslotRepository handles timeslot data operations
type TimeslotRepository struct {
	db *gorm.DB
}

// NewTimeslotRepository creates a new timeslot repository
func NewTimeslotRepository(db *gorm.DB) *TimeslotRepository {
	return &TimeslotRepository{db: db}
}

// Create creates a new timeslot
func (r *TimeslotRepository) Create(timeslot *models.Timeslot) error {
	return r.db.Create(timeslot).Error
}

// FindAll retrieves all timeslots
func (r *TimeslotRepository) FindAll() ([]models.Timeslot, error) {
	var timeslots []models.Timeslot
	err := r.db.Where("is_active = ?", true).Order("time ASC").Find(&timeslots).Error
	return timeslots, err
}

// FindByID finds a timeslot by ID
func (r *TimeslotRepository) FindByID(id uint) (*models.Timeslot, error) {
	var timeslot models.Timeslot
	err := r.db.First(&timeslot, id).Error
	if err != nil {
		return nil, err
	}
	return &timeslot, nil
}

// Update updates a timeslot
func (r *TimeslotRepository) Update(timeslot *models.Timeslot) error {
	return r.db.Save(timeslot).Error
}

// Delete soft deletes a timeslot
func (r *TimeslotRepository) Delete(id uint) error {
	return r.db.Delete(&models.Timeslot{}, id).Error
}