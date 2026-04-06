package handlers

import (
	"net/http"
	"strconv"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// BookingHandler handles HTTP requests for bookings (Owner API)
type BookingHandler struct {
	service *services.BookingService
}

// NewBookingHandler creates a new booking handler
func NewBookingHandler() *BookingHandler {
	return &BookingHandler{
		service: services.NewBookingService(),
	}
}

// List handles GET /api/bookings
func (h *BookingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	result, err := h.service.List(c.Request.Context(), page, pageSize, sortBy, sortOrder, status)
	if err != nil {
		BadRequest(c, "Failed to list bookings: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}

// GetByID handles GET /api/bookings/{id}
func (h *BookingHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	booking, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "Booking")
		return
	}

	SuccessResponse(c, http.StatusOK, booking)
}

// Cancel handles DELETE /api/bookings/{id}
func (h *BookingHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Cancel(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "Booking")
		return
	}

	NoContent(c)
}
