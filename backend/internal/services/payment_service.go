package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reservation-api/api/dto"
	"reservation-api/internal/config"
	"reservation-api/internal/models"
	"reservation-api/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DefaultSessionPrice = 100000.0 // IDR
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo     *repository.PaymentRepository
	reservationRepo *repository.ReservationRepository
	config          *config.Config
}

// NewPaymentService creates a new payment service
func NewPaymentService(
	paymentRepo *repository.PaymentRepository,
	reservationRepo *repository.ReservationRepository,
	cfg *config.Config,
) *PaymentService {
	return &PaymentService{
		paymentRepo:     paymentRepo,
		reservationRepo: reservationRepo,
		config:          cfg,
	}
}

// CreatePayment creates a payment transaction
func (s *PaymentService) CreatePayment(userID uint, req dto.CreatePaymentRequest) (*models.Payment, string, string, error) {
	// Get reservation
	reservation, err := s.reservationRepo.FindByID(req.ReservationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("reservation not found")
		}
		return nil, "", "", err
	}

	// Verify ownership
	if reservation.UserID != userID {
		return nil, "", "", errors.New("unauthorized access to reservation")
	}

	// Check if reservation is cancelled
	if reservation.Status == models.StatusCancelled {
		return nil, "", "", errors.New("cannot pay for cancelled reservation")
	}

	// Check if already paid
	paid, err := s.paymentRepo.CheckPaidByReservationID(req.ReservationID)
	if err != nil {
		return nil, "", "", err
	}
	if paid {
		return nil, "", "", errors.New("reservation already paid")
	}

	// Calculate amount
	amount := DefaultSessionPrice

	// Generate unique transaction ID
	transactionID := fmt.Sprintf("TRX-%s-%d", uuid.New().String()[:8], time.Now().Unix())

	// Create payment record
	payment := &models.Payment{
		ReservationID: req.ReservationID,
		Amount:        amount,
		Status:        models.PaymentPending,
		TransactionID: transactionID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, "", "", errors.New("failed to create payment")
	}

	// Load payment with relations for response
	payment, err = s.paymentRepo.FindByID(payment.ID)
	if err != nil {
		// Continue even if reload fails, payment is already created
		payment.Reservation = *reservation
	}

	// Check if Midtrans is configured
	if s.config.MidtransServerKey == "" || s.config.MidtransClientKey == "" {
		// DUMMY PAYMENT MODE - Skip Midtrans integration
		dummyURL := fmt.Sprintf("http://localhost:3000/payment/dummy?transaction_id=%s", transactionID)
		payment.MidtransURL = dummyURL
		payment.MidtransToken = "DUMMY_TOKEN_" + transactionID

		// Set expiration (24 hours)
		expiredAt := time.Now().Add(24 * time.Hour)
		payment.ExpiredAt = &expiredAt

		if err := s.paymentRepo.Update(payment); err != nil {
			return payment, dummyURL, payment.MidtransToken, err
		}

		return payment, dummyURL, payment.MidtransToken, nil
	}

	// REAL MIDTRANS INTEGRATION (when configured)
	// Create Midtrans transaction
	
	// Handle empty phone (Midtrans requires phone)
	phone := reservation.User.Phone
	if phone == "" {
		phone = "08123456789" // Default dummy phone for users without phone
	}
	
	midtransReq := dto.MidtransRequest{
		TransactionDetails: dto.TransactionDetails{
			OrderID:     transactionID,
			GrossAmount: amount,
		},
		CustomerDetails: dto.CustomerDetails{
			FirstName: reservation.User.Name,
			Email:     reservation.User.Email,
			Phone:     phone, // Use phone or default
		},
		ItemDetails: []dto.ItemDetail{
			{
				ID:       fmt.Sprintf("COURT-%d", reservation.CourtID),
				Price:    amount,
				Quantity: 1,
				Name:     fmt.Sprintf("Pilates Class - %s at %s", reservation.Court.Name, reservation.Timeslot.Time),
			},
		},
	}

	// Call Midtrans API
	midtransResp, err := s.createMidtransTransaction(midtransReq)
	if err != nil {
		// Return payment ID even if Midtrans fails, so user can retry
		return payment, "", "", fmt.Errorf("failed to create payment transaction: %v", err)
	}

	// Update payment with Midtrans info
	payment.MidtransToken = midtransResp.Token
	payment.MidtransURL = midtransResp.RedirectURL

	// Set expiration (24 hours)
	expiredAt := time.Now().Add(24 * time.Hour)
	payment.ExpiredAt = &expiredAt

	if err := s.paymentRepo.Update(payment); err != nil {
		return payment, midtransResp.RedirectURL, midtransResp.Token, err
	}

	return payment, midtransResp.RedirectURL, midtransResp.Token, nil
}

// HandleCallback handles payment callback from Midtrans
func (s *PaymentService) HandleCallback(req dto.PaymentCallbackRequest) (*models.Payment, error) {
	// Find payment by transaction ID
	payment, err := s.paymentRepo.FindByTransactionID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	// Update payment status based on Midtrans callback
	switch req.TransactionStatus {
	case "capture", "settlement":
		payment.Status = models.PaymentPaid
		now := time.Now()
		payment.PaidAt = &now

		// Update reservation status to confirmed
		reservation, err := s.reservationRepo.FindByID(payment.ReservationID)
		if err == nil {
			reservation.Status = models.StatusConfirmed
			s.reservationRepo.Update(reservation)
		}

	case "pending":
		payment.Status = models.PaymentPending

	case "deny", "expire", "cancel":
		payment.Status = models.PaymentFailed

		// Update reservation status to cancelled
		reservation, err := s.reservationRepo.FindByID(payment.ReservationID)
		if err == nil {
			reservation.Status = models.StatusCancelled
			s.reservationRepo.Update(reservation)
		}
	}

	// Save payment
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, errors.New("failed to update payment status")
	}

	// Load full payment details
	payment, err = s.paymentRepo.FindByID(payment.ID)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetPayment gets payment details
func (s *PaymentService) GetPayment(id, userID uint) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	// Verify ownership
	if payment.Reservation.UserID != userID {
		return nil, errors.New("unauthorized access to payment")
	}

	return payment, nil
}

// createMidtransTransaction calls Midtrans Snap API
func (s *PaymentService) createMidtransTransaction(req dto.MidtransRequest) (*dto.MidtransResponse, error) {
	// Check if Midtrans credentials are set
	if s.config.MidtransServerKey == "" {
		return nil, errors.New("midtrans credentials not configured")
	}

	// Marshal request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	url := s.config.MidtransBaseURL + "/transactions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	auth := base64.StdEncoding.EncodeToString([]byte(s.config.MidtransServerKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("midtrans API returned status %d", resp.StatusCode)
	}

	// Parse response
	var midtransResp dto.MidtransResponse
	if err := json.NewDecoder(resp.Body).Decode(&midtransResp); err != nil {
		return nil, err
	}

	return &midtransResp, nil
}