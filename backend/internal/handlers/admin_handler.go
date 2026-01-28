package handlers

import (
	"net/http"
	"reservation-api/internal/models"
	"reservation-api/internal/repository"
	"reservation-api/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin requests
type AdminHandler struct {
	courtRepo    *repository.CourtRepository
	timeslotRepo *repository.TimeslotRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	courtRepo *repository.CourtRepository,
	timeslotRepo *repository.TimeslotRepository,
) *AdminHandler {
	return &AdminHandler{
		courtRepo:    courtRepo,
		timeslotRepo: timeslotRepo,
	}
}

// Courts Management

// GetCourts gets all courts
func (h *AdminHandler) GetCourts(c *gin.Context) {
	courts, err := h.courtRepo.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve courts")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Courts retrieved successfully", gin.H{
		"courts": courts,
	})
}

// CreateCourt creates a new court
func (h *AdminHandler) CreateCourt(c *gin.Context) {
	var court models.Court
	if err := c.ShouldBindJSON(&court); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Set default values
	court.IsActive = true

	if err := h.courtRepo.Create(&court); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create court")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Court created successfully", gin.H{
		"court": court,
	})
}

// UpdateCourt updates a court
func (h *AdminHandler) UpdateCourt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid court ID")
		return
	}

	court, err := h.courtRepo.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Court not found")
		return
	}

	var updates models.Court
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update fields
	if updates.Name != "" {
		court.Name = updates.Name
	}
	if updates.Capacity > 0 {
		court.Capacity = updates.Capacity
	}
	if updates.Description != "" {
		court.Description = updates.Description
	}

	if err := h.courtRepo.Update(court); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update court")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Court updated successfully", gin.H{
		"court": court,
	})
}

// DeleteCourt deletes a court
func (h *AdminHandler) DeleteCourt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid court ID")
		return
	}

	if err := h.courtRepo.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete court")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Court deleted successfully", nil)
}

// Timeslots Management

// GetTimeslots gets all timeslots
func (h *AdminHandler) GetTimeslots(c *gin.Context) {
	timeslots, err := h.timeslotRepo.FindAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve timeslots")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Timeslots retrieved successfully", gin.H{
		"timeslots": timeslots,
	})
}

// CreateTimeslot creates a new timeslot
func (h *AdminHandler) CreateTimeslot(c *gin.Context) {
	var timeslot models.Timeslot
	if err := c.ShouldBindJSON(&timeslot); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Set default values
	timeslot.IsActive = true

	if err := h.timeslotRepo.Create(&timeslot); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create timeslot")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Timeslot created successfully", gin.H{
		"timeslot": timeslot,
	})
}

// UpdateTimeslot updates a timeslot
func (h *AdminHandler) UpdateTimeslot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid timeslot ID")
		return
	}

	timeslot, err := h.timeslotRepo.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Timeslot not found")
		return
	}

	var updates models.Timeslot
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Update fields
	if updates.Time != "" {
		timeslot.Time = updates.Time
	}
	if updates.Duration > 0 {
		timeslot.Duration = updates.Duration
	}

	if err := h.timeslotRepo.Update(timeslot); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update timeslot")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Timeslot updated successfully", gin.H{
		"timeslot": timeslot,
	})
}

// DeleteTimeslot deletes a timeslot
func (h *AdminHandler) DeleteTimeslot(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid timeslot ID")
		return
	}

	if err := h.timeslotRepo.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete timeslot")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Timeslot deleted successfully", nil)
}

// GetStatistics gets admin statistics
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	// TODO: Implement statistics gathering
	// This is a placeholder for admin dashboard statistics

	stats := gin.H{
		"total_courts":    3,
		"total_timeslots": 7,
		"message":         "Statistics endpoint - implement as needed",
	}

	utils.SuccessResponse(c, http.StatusOK, "Statistics retrieved successfully", stats)
}