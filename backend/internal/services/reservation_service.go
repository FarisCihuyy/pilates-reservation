package services

import (
	"errors"
	"reservation-api/api/dto"
	"reservation-api/internal/models"
	"reservation-api/internal/repository"
	"time"

	"gorm.io/gorm"
)

// ReservationService handles reservation business logic
type ReservationService struct {
	reservationRepo *repository.ReservationRepository
	courtRepo       *repository.CourtRepository
	timeslotRepo    *repository.TimeslotRepository
}

// NewReservationService creates a new reservation service
func NewReservationService(
	reservationRepo *repository.ReservationRepository,
	courtRepo *repository.CourtRepository,
	timeslotRepo *repository.TimeslotRepository,
) *ReservationService {
	return &ReservationService{
		reservationRepo: reservationRepo,
		courtRepo:       courtRepo,
		timeslotRepo:    timeslotRepo,
	}
}

// CreateReservation creates a new reservation
func (s *ReservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*models.Reservation, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format. Use YYYY-MM-DD")
	}

	// Check if date is in the past
	today := time.Now().Truncate(24 * time.Hour)
	if date.Before(today) {
		return nil, errors.New("cannot book past dates")
	}

	// Verify court exists
	court, err := s.courtRepo.FindByID(req.CourtID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("court not found")
		}
		return nil, err
	}

	if !court.IsActive {
		return nil, errors.New("court is not active")
	}

	// Verify timeslot exists
	timeslot, err := s.timeslotRepo.FindByID(req.TimeslotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("timeslot not found")
		}
		return nil, err
	}

	if !timeslot.IsActive {
		return nil, errors.New("timeslot is not active")
	}

	// Check if court has available capacity (Group Class)
	available, bookedCount, err := s.reservationRepo.CheckAvailability(req.CourtID, req.TimeslotID, date)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.New("failed to check availability")
	}

	// Check if court is full
	if bookedCount >= court.Capacity {
		return nil, errors.New("this class is already full. Please select another court or timeslot.")
	}

	// Create reservation
	reservation := &models.Reservation{
		UserID:     userID,
		CourtID:    req.CourtID,
		TimeslotID: req.TimeslotID,
		Date:       date,
		Status:     models.StatusPending,
		Notes:      req.Notes,
	}

	if err := s.reservationRepo.Create(reservation); err != nil {
		return nil, errors.New("failed to create reservation")
	}

	// Load relations
	reservation, err = s.reservationRepo.FindByID(reservation.ID)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

// GetUserReservations gets all reservations for a user
func (s *ReservationService) GetUserReservations(userID uint) ([]models.Reservation, error) {
	return s.reservationRepo.FindByUserID(userID)
}

// GetReservation gets a single reservation
func (s *ReservationService) GetReservation(id, userID uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	// Verify ownership
	if reservation.UserID != userID {
		return nil, errors.New("unauthorized access to reservation")
	}

	return reservation, nil
}

// CancelReservation cancels a reservation
func (s *ReservationService) CancelReservation(id, userID uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	// Verify ownership
	if reservation.UserID != userID {
		return nil, errors.New("unauthorized access to reservation")
	}

	// Check if can be cancelled
	if !reservation.CanBeCancelled() {
		return nil, errors.New("reservation cannot be cancelled")
	}

	// Check if reservation date is in the past
	if reservation.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("cannot cancel past reservations")
	}

	// Update status
	reservation.Status = models.StatusCancelled
	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, errors.New("failed to cancel reservation")
	}

	return reservation, nil
}

// GetTimeslotsAvailability gets timeslots with availability info for a date
func (s *ReservationService) GetTimeslotsAvailability(dateStr string) ([]dto.TimeslotAvailability, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format. Use YYYY-MM-DD")
	}

	// Get all timeslots
	timeslots, err := s.timeslotRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Get total courts
	totalCourts, err := s.courtRepo.CountAll()
	if err != nil {
		return nil, err
	}

	// Build availability info
	var result []dto.TimeslotAvailability
	for _, ts := range timeslots {
		bookedCount, err := s.reservationRepo.CountBookedByDateAndTimeslot(date, ts.ID)
		if err != nil {
			return nil, err
		}

		availableCourts := int(totalCourts) - int(bookedCount)
		available := availableCourts > 0

		result = append(result, dto.TimeslotAvailability{
			ID:              ts.ID,
			Time:            ts.Time,
			Duration:        ts.Duration,
			IsActive:        ts.IsActive,
			Available:       available,
			BookedCount:     int(bookedCount),
			AvailableCourts: availableCourts,
		})
	}

	return result, nil
}

// GetCourtsAvailability gets courts availability for a date and timeslot
func (s *ReservationService) GetCourtsAvailability(dateStr string, timeslotID uint) ([]dto.CourtAvailability, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New("invalid date format. Use YYYY-MM-DD")
	}

	// Get all courts
	courts, err := s.courtRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Get booked court IDs
	bookedCourtIDs, err := s.reservationRepo.GetBookedCourtIDs(date, timeslotID)
	if err != nil {
		return nil, err
	}

	// Build availability info
	var result []dto.CourtAvailability
	for _, court := range courts {
		available := true
		for _, bookedID := range bookedCourtIDs {
			if court.ID == bookedID {
				available = false
				break
			}
		}

		result = append(result, dto.CourtAvailability{
			ID:          court.ID,
			Name:        court.Name,
			Capacity:    court.Capacity,
			Description: court.Description,
			IsActive:    court.IsActive,
			Available:   available,
		})
	}

	return result, nil
}