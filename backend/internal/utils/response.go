package utils

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, err error) {
	c.JSON(400, gin.H{
		"success": false,
		"error":   "Validation failed",
		"details": err.Error(),
	})
}