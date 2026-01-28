package dto

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}