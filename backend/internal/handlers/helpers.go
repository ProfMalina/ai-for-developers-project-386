package handlers

import (
	"net/http"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c *gin.Context, status int, errorType, message string) {
	c.JSON(status, models.ErrorResponse{
		Error:   errorType,
		Message: message,
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, fieldErrors []models.FieldError) {
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Error:       "VALIDATION_ERROR",
		Message:     "Validation failed",
		FieldErrors: fieldErrors,
	})
}

// NotFound sends a 404 error response
func NotFound(c *gin.Context, resource string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", resource+" not found")
}

// BadRequest sends a 400 error response
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Conflict sends a 409 error response
func Conflict(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message)
}

// NoContent sends a 204 response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
