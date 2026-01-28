package handlers

import (
	"net/http"
	"reservation-api/api/dto"
	"reservation-api/internal/middleware"
	"reservation-api/internal/services"
	"reservation-api/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles payment requests
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment creates a payment transaction
// @Summary Create payment
// @Description Create a payment transaction for a reservation
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePaymentRequest true "Payment details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /payments/create [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	payment, paymentURL, snapToken, err := h.paymentService.CreatePayment(userID, req)
	if err != nil {
		// Still return payment info even if Midtrans fails
		if payment != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get client key from config for frontend integration
	clientKey := "SB-Mid-client-xxxxxxxxxx" // This should come from config

	utils.SuccessResponse(c, http.StatusOK, "Payment created successfully", gin.H{
		"payment":     payment,
		"payment_url": paymentURL,
		"snap_token":  snapToken,
		"client_key":  clientKey,
	})
}

// PaymentCallback handles payment callback from Midtrans
// @Summary Payment callback
// @Description Handle payment status callback from Midtrans
// @Tags payments
// @Accept json
// @Produce json
// @Param request body dto.PaymentCallbackRequest true "Callback data"
// @Success 200 {object} map[string]interface{}
// @Router /payments/callback [post]
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	var req dto.PaymentCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	payment, err := h.paymentService.HandleCallback(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment status updated", gin.H{
		"payment": payment,
	})
}

// GetPayment gets payment details
// @Summary Get payment details
// @Description Get details of a specific payment
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	payment, err := h.paymentService.GetPayment(uint(id), userID)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "unauthorized access to payment" {
			status = http.StatusForbidden
		}
		utils.ErrorResponse(c, status, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", gin.H{
		"payment": payment,
	})
}