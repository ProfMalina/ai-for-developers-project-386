package handlers

import (
	"net/http"
	"strings"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PublicBookingHandler handles HTTP requests for public bookings
type PublicBookingHandler struct {
	service *services.BookingService
}

// NewPublicBookingHandler creates a new public booking handler
func NewPublicBookingHandler() *PublicBookingHandler {
	return &PublicBookingHandler{
		service: services.NewBookingService(
			repositories.NewBookingRepository(),
			repositories.NewTimeSlotRepository(),
			repositories.NewEventTypeRepository(),
		),
	}
}

// Create handles POST /api/public/bookings
func (h *PublicBookingHandler) Create(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrors validator.ValidationErrors
		if errorsAs := validator.ValidationErrors(nil); strings.Contains(err.Error(), "validation") || strings.Contains(err.Error(), "Field validation") {
			_ = errorsAs
		}
		if ok := errorAsValidation(err, &validationErrors); ok {
			fieldErrors := make([]models.FieldError, 0, len(validationErrors))
			for _, fieldErr := range validationErrors {
				fieldErrors = append(fieldErrors, models.FieldError{
					Field:   lowerFirst(fieldErr.Field()),
					Message: validationMessage(fieldErr),
				})
			}
			ValidationError(c, fieldErrors)
			return
		}
		BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate required fields
	if req.EventTypeID == "" {
		BadRequest(c, "eventTypeId is required")
		return
	}

	if req.SlotID == nil || *req.SlotID == "" {
		BadRequest(c, "slotId is required")
		return
	}

	if req.GuestName == "" {
		BadRequest(c, "guestName is required")
		return
	}

	if req.GuestEmail == "" {
		BadRequest(c, "guestEmail is required")
		return
	}

	booking, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		errText := err.Error()
		if strings.Contains(errText, "selected time slot is already booked") || strings.Contains(errText, "already booked") || strings.Contains(errText, "overlap") {
			Conflict(c, "Selected time slot is already booked")
			return
		}
		if strings.Contains(errText, "event type not found") {
			NotFound(c, "Event type")
			return
		}
		if strings.Contains(errText, "time slot not found") {
			NotFound(c, "Time slot")
			return
		}
		if strings.Contains(errText, "already started or passed") {
			InvalidTime(c, "Cannot book a slot that has already started or ended")
			return
		}
		BadRequest(c, "Failed to create booking: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, booking)
}
