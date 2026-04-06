package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// TimeSlotHandler handles HTTP requests for time slots
type TimeSlotHandler struct {
	service *services.TimeSlotService
}

// NewTimeSlotHandler creates a new time slot handler
func NewTimeSlotHandler() *TimeSlotHandler {
	return &TimeSlotHandler{
		service: services.NewTimeSlotService(),
	}
}

// List handles GET /api/slots
func (h *TimeSlotHandler) List(c *gin.Context) {
	eventTypeID := c.Query("eventTypeId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var available *bool
	if availStr := c.Query("available"); availStr != "" {
		avail := availStr == "true"
		available = &avail
	}

	var startTime *time.Time
	if startStr := c.Query("startTime"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			startTime = &t
		}
	}

	var endTime *time.Time
	if endStr := c.Query("endTime"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			endTime = &t
		}
	}

	result, err := h.service.List(c.Request.Context(), eventTypeID, page, pageSize, available, startTime, endTime)
	if err != nil {
		BadRequest(c, "Failed to list time slots: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}

// GenerateSlots handles POST /api/event-types/{id}/slots/generate
func (h *TimeSlotHandler) GenerateSlots(c *gin.Context) {
	ownerID := c.GetString("ownerID")
	if ownerID == "" {
		ownerID = db.DefaultOwnerID
	}

	eventTypeID := c.Param("id")

	var req models.SlotGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Set eventTypeId from URL path if not provided in request body
	if req.EventTypeID == "" {
		req.EventTypeID = eventTypeID
	}

	// Validate days of week
	if len(req.DaysOfWeek) < 1 || len(req.DaysOfWeek) > 7 {
		BadRequest(c, "days_of_week must contain 1-7 items")
		return
	}

	result, err := h.service.GenerateSlots(c.Request.Context(), ownerID, req)
	if err != nil {
		BadRequest(c, "Failed to generate slots: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, result)
}
