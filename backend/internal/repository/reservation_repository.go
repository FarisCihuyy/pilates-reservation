package repository

import (
	"reservation-api/internal/models"
	"time"

	"gorm.io/gorm"
)

// ReservationRepository handles reservation data operations
type ReservationRepository struct {
	db *gorm.DB
}

// NewReservationRepository creates a new reservation repository
func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// Create creates a new reservation
func (r *ReservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

// FindByID finds a reservation by ID with relations
func (r *ReservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.Preload("User").
		Preload("Court").
		Preload("Timeslot").
		Preload("Payment").
		First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// FindByUserID finds all reservations by user ID
func (r *ReservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Court").
		Preload("Timeslot").
		Preload("Payment").
		Where("user_id = ?", userID).
		Order("date DESC, created_at DESC").
		Find(&reservations).Error
	return reservations, err
}

// Update updates a reservation
func (r *ReservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

// CheckAvailability checks if a court has available capacity for a specific date and timeslot
func (r *ReservationRepository) CheckAvailability(courtID, timeslotID uint, date time.Time) (bool, int, error) {
	// Count ONLY CONFIRMED bookings (exclude pending & cancelled)
	var bookedCount int64
	err := r.db.Model(&models.Reservation{}).
		Where("court_id = ? AND timeslot_id = ? AND date = ? AND status = ?",
			courtID, timeslotID, date, models.StatusConfirmed).
		Count(&bookedCount).Error

	if err != nil {
		return false, 0, err
	}

	return true, int(bookedCount), nil
}

// CountBookedByDateAndTimeslot counts booked reservations for a date and timeslot
func (r *ReservationRepository) CountBookedByDateAndTimeslot(date time.Time, timeslotID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).
		Where("date = ? AND timeslot_id = ? AND status != ?",
			date, timeslotID, models.StatusCancelled).
		Count(&count).Error
	return count, err
}

// GetBookedCourtIDs gets list of booked court IDs for a specific date and timeslot
func (r *ReservationRepository) GetBookedCourtIDs(date time.Time, timeslotID uint) ([]uint, error) {
	var courtIDs []uint
	err := r.db.Model(&models.Reservation{}).
		Where("date = ? AND timeslot_id = ? AND status != ?",
			date, timeslotID, models.StatusCancelled).
		Pluck("court_id", &courtIDs).Error
	return courtIDs, err
}

// GetUpcomingReservations gets upcoming reservations
func (r *ReservationRepository) GetUpcomingReservations(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	today := time.Now().Truncate(24 * time.Hour)

	err := r.db.Preload("Court").
		Preload("Timeslot").
		Preload("Payment").
		Where("user_id = ? AND date >= ? AND status != ?", userID, today, models.StatusCancelled).
		Order("date ASC").
		Find(&reservations).Error

	return reservations, err
}

// GetPastReservations gets past reservations
func (r *ReservationRepository) GetPastReservations(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	today := time.Now().Truncate(24 * time.Hour)

	err := r.db.Preload("Court").
		Preload("Timeslot").
		Preload("Payment").
		Where("user_id = ? AND date < ?", userID, today).
		Order("date DESC").
		Find(&reservations).Error

	return reservations, err
}