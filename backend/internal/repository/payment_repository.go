package repository

import (
	"reservation-api/internal/models"

	"gorm.io/gorm"
)

// PaymentRepository handles payment data operations
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create creates a new payment
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// FindByID finds a payment by ID
func (r *PaymentRepository) FindByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Reservation").
		Preload("Reservation.User").
		Preload("Reservation.Court").
		Preload("Reservation.Timeslot").
		First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByReservationID finds a payment by reservation ID
func (r *PaymentRepository) FindByReservationID(reservationID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("reservation_id = ?", reservationID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByTransactionID finds a payment by transaction ID
func (r *PaymentRepository) FindByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update updates a payment
func (r *PaymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// CheckPaidByReservationID checks if reservation has been paid
func (r *PaymentRepository) CheckPaidByReservationID(reservationID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Payment{}).
		Where("reservation_id = ? AND status = ?", reservationID, models.PaymentPaid).
		Count(&count).Error
	return count > 0, err
}