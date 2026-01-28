package services

import (
	"errors"
	"reservation-api/api/dto"
	"reservation-api/internal/config"
	"reservation-api/internal/models"
	"reservation-api/internal/repository"
	"reservation-api/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register registers a new user
func (s *AuthService) Register(req dto.RegisterRequest) (*models.User, string, error) {
	// Check if email already exists
	exists, err := s.userRepo.Exists(req.Email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", errors.New("failed to create user")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	// Don't return password
	user.Password = ""

	return user, token, nil
}

// Login authenticates a user
func (s *AuthService) Login(req dto.LoginRequest) (*models.User, string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid email or password")
		}
		return nil, "", err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, "", errors.New("user account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	// Don't return password
	user.Password = ""

	return user, token, nil
}

// GetProfile gets user profile
func (s *AuthService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Don't return password
	user.Password = ""

	return user, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*models.User, error) {
	// Find user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	// Don't return password
	user.Password = ""

	return user, nil
}