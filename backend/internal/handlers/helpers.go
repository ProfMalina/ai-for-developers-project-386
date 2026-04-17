package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		Message:     "Request validation failed",
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

// InvalidTime sends a 400 invalid time response
func InvalidTime(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, "INVALID_TIME", message)
}

// Conflict sends a 409 error response
func Conflict(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message)
}

// NoContent sends a 204 response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func errorAsValidation(err error, target *validator.ValidationErrors) bool {
	return errors.As(err, target)
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func validationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "uuid":
		return "must be a valid UUID"
	case "min":
		return "is below the minimum value"
	case "max":
		return "exceeds the maximum value"
	case "datetime":
		return "must match the required datetime format"
	case "oneof":
		return "must be one of the allowed values"
	default:
		return "is invalid"
	}
}
