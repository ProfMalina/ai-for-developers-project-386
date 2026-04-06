package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// PublicEventTypeHandler handles HTTP requests for public event types
type PublicEventTypeHandler struct {
	etService    *services.EventTypeService
	slotService  *services.TimeSlotService
}

// NewPublicEventTypeHandler creates a new public event type handler
func NewPublicEventTypeHandler() *PublicEventTypeHandler {
	return &PublicEventTypeHandler{
		etService:   services.NewEventTypeService(),
		slotService: services.NewTimeSlotService(),
	}
}

// List handles GET /api/public/event-types
func (h *PublicEventTypeHandler) List(c *gin.Context) {
	ownerID := db.DefaultOwnerID

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	result, err := h.etService.List(c.Request.Context(), ownerID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		BadRequest(c, "Failed to list event types: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}

// GetByID handles GET /api/public/event-types/{id}
func (h *PublicEventTypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	et, err := h.etService.GetByID(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "Event type")
		return
	}

	// Only show active event types
	if !et.IsActive {
		NotFound(c, "Event type")
		return
	}

	SuccessResponse(c, http.StatusOK, et)
}

// GetSlots handles GET /api/public/slots
func (h *PublicEventTypeHandler) GetSlots(c *gin.Context) {
	ownerID := db.DefaultOwnerID

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// Parse optional dateFrom and dateTo query parameters
	var startTime, endTime *time.Time
	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		if t, err := time.Parse(time.RFC3339, dateFrom); err == nil {
			startTime = &t
		}
	}
	if dateTo := c.Query("dateTo"); dateTo != "" {
		if t, err := time.Parse(time.RFC3339, dateTo); err == nil {
			endTime = &t
		}
	}

	result, err := h.slotService.GetAvailableSlots(c.Request.Context(), ownerID, page, pageSize, startTime, endTime)
	if err != nil {
		BadRequest(c, "Failed to get available slots: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}
