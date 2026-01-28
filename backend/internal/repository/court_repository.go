package repository

import (
	"reservation-api/internal/models"

	"gorm.io/gorm"
)

// CourtRepository handles court data operations
type CourtRepository struct {
	db *gorm.DB
}

// NewCourtRepository creates a new court repository
func NewCourtRepository(db *gorm.DB) *CourtRepository {
	return &CourtRepository{db: db}
}

// Create creates a new court
func (r *CourtRepository) Create(court *models.Court) error {
	return r.db.Create(court).Error
}

// FindAll retrieves all courts
func (r *CourtRepository) FindAll() ([]models.Court, error) {
	var courts []models.Court
	err := r.db.Where("is_active = ?", true).Find(&courts).Error
	return courts, err
}

// FindByID finds a court by ID
func (r *CourtRepository) FindByID(id uint) (*models.Court, error) {
	var court models.Court
	err := r.db.First(&court, id).Error
	if err != nil {
		return nil, err
	}
	return &court, nil
}

// Update updates a court
func (r *CourtRepository) Update(court *models.Court) error {
	return r.db.Save(court).Error
}

// Delete soft deletes a court
func (r *CourtRepository) Delete(id uint) error {
	return r.db.Delete(&models.Court{}, id).Error
}

// CountAll counts all active courts
func (r *CourtRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Court{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}