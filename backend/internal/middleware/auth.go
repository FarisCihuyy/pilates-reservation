package middleware

import (
	"net/http"
	"reservation-api/internal/config"
	"reservation-api/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and sets user info in context
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// AdminMiddleware checks if user has admin privileges
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement admin check
		// For now, we'll allow all authenticated users to access admin routes
		// In production, you should check user role from database

		// Example implementation:
		// userID := c.GetUint("user_id")
		// user, err := userRepo.FindByID(userID)
		// if err != nil || !user.IsAdmin {
		//     c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		//     c.Abort()
		//     return
		// }

		c.Next()
	}
}

// GetUserID gets user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetEmail gets email from context
func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	return email.(string), true
}