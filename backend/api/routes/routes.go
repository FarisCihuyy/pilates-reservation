package routes

import (
	"reservation-api/internal/config"
	"reservation-api/internal/handlers"
	"reservation-api/internal/middleware"
	"reservation-api/internal/repository"
	"reservation-api/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	courtRepo := repository.NewCourtRepository(db)
	timeslotRepo := repository.NewTimeslotRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	reservationService := services.NewReservationService(reservationRepo, courtRepo, timeslotRepo)
	paymentService := services.NewPaymentService(paymentRepo, reservationRepo, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	reservationHandler := handlers.NewReservationHandler(reservationService, courtRepo, timeslotRepo)
	paymentHandler := handlers.NewPaymentHandler(paymentService)
	adminHandler := handlers.NewAdminHandler(courtRepo, timeslotRepo)

	// Setup middleware
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Pilates Reservation API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes - Authentication
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public routes - Browse available slots
		v1.GET("/dates", reservationHandler.GetAvailableDates)
		v1.GET("/timeslots", reservationHandler.GetTimeslots)
		v1.GET("/courts", reservationHandler.GetAvailableCourts)

		// Protected routes - Require authentication
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Reservations
			reservations := protected.Group("/reservations")
			{
				reservations.POST("", reservationHandler.CreateReservation)
				reservations.GET("", reservationHandler.GetUserReservations)
				reservations.GET("/:id", reservationHandler.GetReservation)
				reservations.PUT("/:id/cancel", reservationHandler.CancelReservation)
			}

			// Payments
			payments := protected.Group("/payments")
			{
				payments.POST("/create", paymentHandler.CreatePayment)
				payments.POST("/callback", paymentHandler.PaymentCallback)
				payments.GET("/:id", paymentHandler.GetPayment)
			}

			// User profile
			profile := protected.Group("/profile")
			{
				profile.GET("", authHandler.GetProfile)
				profile.PUT("", authHandler.UpdateProfile)
			}
		}

		// Admin routes - For managing courts and timeslots
		admin := v1.Group("/admin")
		// admin.Use(middleware.AdminMiddleware()) // TODO: Implement admin middleware
		{
			// Courts management
			courts := admin.Group("/courts")
			{
				courts.GET("", adminHandler.GetCourts)
				courts.POST("", adminHandler.CreateCourt)
				courts.PUT("/:id", adminHandler.UpdateCourt)
				courts.DELETE("/:id", adminHandler.DeleteCourt)
			}

			// Timeslots management
			timeslots := admin.Group("/timeslots")
			{
				timeslots.GET("", adminHandler.GetTimeslots)
				timeslots.POST("", adminHandler.CreateTimeslot)
				timeslots.PUT("/:id", adminHandler.UpdateTimeslot)
				timeslots.DELETE("/:id", adminHandler.DeleteTimeslot)
			}

			// Dashboard statistics
			admin.GET("/stats", adminHandler.GetStatistics)
		}
	}

	// Legacy routes for backward compatibility
	setupLegacyRoutes(router, authHandler, reservationHandler, paymentHandler, adminHandler, cfg)
}

// setupLegacyRoutes sets up backward compatible routes
func setupLegacyRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	reservationHandler *handlers.ReservationHandler,
	paymentHandler *handlers.PaymentHandler,
	adminHandler *handlers.AdminHandler,
	cfg *config.Config,
) {
	api := router.Group("/api")
	{
		// Auth
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Public
		api.GET("/dates", reservationHandler.GetAvailableDates)
		api.GET("/timeslots", reservationHandler.GetTimeslots)
		api.GET("/courts", reservationHandler.GetAvailableCourts)

		// Protected
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.POST("/reservations", reservationHandler.CreateReservation)
			protected.GET("/reservations", reservationHandler.GetUserReservations)
			protected.GET("/reservations/:id", reservationHandler.GetReservation)

			protected.POST("/payment/create", paymentHandler.CreatePayment)
			protected.POST("/payment/callback", paymentHandler.PaymentCallback)
		}

		// Admin
		admin := api.Group("/admin")
		{
			admin.POST("/courts", adminHandler.CreateCourt)
			admin.POST("/timeslots", adminHandler.CreateTimeslot)
		}
	}
}