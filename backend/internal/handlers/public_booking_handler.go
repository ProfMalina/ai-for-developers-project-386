package handlers

import (
	"net/http"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
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
		// Check for conflict error
		if err.Error() == "selected time slot is already booked" {
			Conflict(c, "Selected time slot is already booked")
			return
		}
		BadRequest(c, "Failed to create booking: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, booking)
}
