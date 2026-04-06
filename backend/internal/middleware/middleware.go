package middleware

import (
	"net/http"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// ErrorHandler handles errors and returns proper JSON responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			status := http.StatusInternalServerError
			errorType := "INTERNAL_ERROR"
			message := "Internal server error"

			// Determine error type from status code
			switch c.Writer.Status() {
			case http.StatusBadRequest:
				errorType = "BAD_REQUEST"
				message = "Bad request"
			case http.StatusNotFound:
				errorType = "NOT_FOUND"
				message = "Resource not found"
			case http.StatusConflict:
				errorType = "CONFLICT"
				message = "Resource conflict"
			}

			c.JSON(status, models.ErrorResponse{
				Error:   errorType,
				Message: message,
				Details: err.Error(),
			})
		}
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
