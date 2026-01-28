package dto

// CreateReservationRequest represents reservation creation request
type CreateReservationRequest struct {
	CourtID    uint   `json:"court_id" binding:"required"`
	TimeslotID uint   `json:"timeslot_id" binding:"required"`
	Date       string `json:"date" binding:"required"` // Format: YYYY-MM-DD
	Notes      string `json:"notes"`
}

// TimeslotAvailability represents timeslot with availability info
type TimeslotAvailability struct {
	ID              uint   `json:"id"`
	Time            string `json:"time"`
	Duration        int    `json:"duration"`
	IsActive        bool   `json:"is_active"`
	Available       bool   `json:"available"`
	BookedCount     int    `json:"booked_count"`
	AvailableCourts int    `json:"available_courts"`
}

// CourtAvailability represents court with availability info
type CourtAvailability struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	Available   bool   `json:"available"`
}