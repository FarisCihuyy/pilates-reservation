package handlers

import (
	"net/http"
	"reservation-api/api/dto"
	"reservation-api/internal/middleware"
	"reservation-api/internal/repository"
	"reservation-api/internal/services"
	"reservation-api/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ReservationHandler handles reservation requests
type ReservationHandler struct {
	reservationService *services.ReservationService
	courtRepo          *repository.CourtRepository
	timeslotRepo       *repository.TimeslotRepository
}

// NewReservationHandler creates a new reservation handler
func NewReservationHandler(
	reservationService *services.ReservationService,
	courtRepo *repository.CourtRepository,
	timeslotRepo *repository.TimeslotRepository,
) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		courtRepo:          courtRepo,
		timeslotRepo:       timeslotRepo,
	}
}

// GetAvailableDates gets available dates for booking
// @Summary Get available dates
// @Description Get list of available dates for next 30 days
// @Tags public
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /dates [get]
func (h *ReservationHandler) GetAvailableDates(c *gin.Context) {
	dates := []string{}
	today := time.Now()

	for i := 0; i < 30; i++ {
		date := today.AddDate(0, 0, i)
		dates = append(dates, date.Format("2006-01-02"))
	}

	utils.SuccessResponse(c, http.StatusOK, "Dates retrieved successfully", gin.H{
		"dates": dates,
	})
}

// GetTimeslots gets timeslots with optional availability check
// @Summary Get timeslots
// @Description Get all timeslots with availability info if date is provided
// @Tags public
// @Produce json
// @Param date query string false "Date in YYYY-MM-DD format"
// @Success 200 {object} map[string]interface{}
// @Router /timeslots [get]
func (h *ReservationHandler) GetTimeslots(c *gin.Context) {
	dateStr := c.Query("date")

	// If date provided, return with availability
	if dateStr != "" {
		timeslots, err := h.reservationService.GetTimeslotsAvailability(dateStr)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		utils.SuccessResponse(c, http.StatusOK, "Timeslots with availability retrieved successfully", gin.H{
			"timeslots": timeslots,
		})
		return
	}

	// Return all timeslots without availability
	timeslots, err := h.timeslotRepo.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve timeslots")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Timeslots retrieved successfully", gin.H{
		"timeslots": timeslots,
	})
}

// GetAvailableCourts gets available courts for a date and timeslot
// @Summary Get available courts
// @Description Get courts availability for specific date and timeslot
// @Tags public
// @Produce json
// @Param date query string true "Date in YYYY-MM-DD format"
// @Param timeslot_id query int true "Timeslot ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /courts [get]
func (h *ReservationHandler) GetAvailableCourts(c *gin.Context) {
	dateStr := c.Query("date")
	timeslotIDStr := c.Query("timeslot_id")

	if dateStr == "" || timeslotIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "date and timeslot_id are required")
		return
	}

	timeslotID, err := strconv.ParseUint(timeslotIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid timeslot_id")
		return
	}

	courts, err := h.reservationService.GetCourtsAvailability(dateStr, uint(timeslotID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Courts availability retrieved successfully", gin.H{
		"courts": courts,
	})
}

// CreateReservation creates a new reservation
// @Summary Create reservation
// @Description Create a new reservation for a court
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateReservationRequest true "Reservation details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	reservation, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "court is already booked for this timeslot" {
			status = http.StatusConflict
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Reservation created successfully. Please proceed to payment.", gin.H{
		"reservation": reservation,
	})
}

// GetUserReservations gets all reservations for the logged-in user
// @Summary Get user reservations
// @Description Get all reservations for the authenticated user
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /reservations [get]
func (h *ReservationHandler) GetUserReservations(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	reservations, err := h.reservationService.GetUserReservations(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve reservations")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reservations retrieved successfully", gin.H{
		"reservations": reservations,
	})
}

// GetReservation gets a single reservation by ID
// @Summary Get reservation details
// @Description Get details of a specific reservation
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /reservations/{id} [get]
func (h *ReservationHandler) GetReservation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reservation ID")
		return
	}

	reservation, err := h.reservationService.GetReservation(uint(id), userID)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "unauthorized access to reservation" {
			status = http.StatusForbidden
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reservation retrieved successfully", gin.H{
		"reservation": reservation,
	})
}

// CancelReservation cancels a reservation
// @Summary Cancel reservation
// @Description Cancel an existing reservation
// @Tags reservations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Reservation ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /reservations/{id}/cancel [put]
func (h *ReservationHandler) CancelReservation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid reservation ID")
		return
	}

	reservation, err := h.reservationService.CancelReservation(uint(id), userID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unauthorized access to reservation" {
			status = http.StatusForbidden
		} else if err.Error() == "reservation not found" {
			status = http.StatusNotFound
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Reservation cancelled successfully", gin.H{
		"reservation": reservation,
	})
}