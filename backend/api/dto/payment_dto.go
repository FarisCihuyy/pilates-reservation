package dto

// CreatePaymentRequest represents payment creation request
type CreatePaymentRequest struct {
	ReservationID uint `json:"reservation_id" binding:"required"`
}

// PaymentCallbackRequest represents Midtrans callback
type PaymentCallbackRequest struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
}

// MidtransRequest represents request to Midtrans API
type MidtransRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails    CustomerDetails    `json:"customer_details"`
	ItemDetails        []ItemDetail       `json:"item_details"`
}

// TransactionDetails for Midtrans
type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount float64 `json:"gross_amount"`
}

// CustomerDetails for Midtrans
type CustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// ItemDetail for Midtrans
type ItemDetail struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Name     string  `json:"name"`
}

// MidtransResponse represents response from Midtrans
type MidtransResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}