package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
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
		service: services.NewTimeSlotService(
			repositories.NewTimeSlotRepository(),
			repositories.NewSlotGenerationConfigRepository(),
			repositories.NewOwnerRepository(),
			repositories.NewEventTypeRepository(),
		),
	}
}

// List handles GET /api/slots
func (h *TimeSlotHandler) List(c *gin.Context) {
	ownerID := c.GetString("ownerID")
	if ownerID == "" {
		ownerID = db.DefaultOwnerID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var available *bool
	if availStr := c.Query("isAvailable"); availStr != "" {
		avail := availStr == "true"
		available = &avail
	} else if availStr := c.Query("available"); availStr != "" {
		avail := availStr == "true"
		available = &avail
	}

	var startTime *time.Time
	if startStr := c.Query("dateFrom"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			BadRequest(c, "dateFrom must be RFC3339")
			return
		}
		startTime = &t
	} else if startStr := c.Query("startTime"); startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err == nil {
			startTime = &t
		}
	}

	var endTime *time.Time
	if endStr := c.Query("dateTo"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			BadRequest(c, "dateTo must be RFC3339")
			return
		}
		endTime = &t
	} else if endStr := c.Query("endTime"); endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err == nil {
			endTime = &t
		}
	}

	eventTypeID := strings.TrimSpace(c.Query("eventTypeId"))

	result, err := h.service.List(c.Request.Context(), ownerID, eventTypeID, page, pageSize, available, startTime, endTime)
	if err != nil {
		if strings.Contains(err.Error(), "event type not found") {
			NotFound(c, "Event type")
			return
		}
		BadRequest(c, "Failed to list time slots: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}

// GenerateSlots handles POST /api/event-types/{eventTypeId}/slots/generate
func (h *TimeSlotHandler) GenerateSlots(c *gin.Context) {
	ownerID := c.GetString("ownerID")
	if ownerID == "" {
		ownerID = db.DefaultOwnerID
	}
	eventTypeID := strings.TrimSpace(c.Param("eventTypeId"))
	if eventTypeID == "" {
		BadRequest(c, "eventTypeId is required")
		return
	}

	var req models.SlotGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Validate days of week when provided
	if len(req.DaysOfWeek) > 7 {
		BadRequest(c, "days_of_week must contain 1-7 items")
		return
	}
	for _, day := range req.DaysOfWeek {
		if day < 0 || day > 6 {
			BadRequest(c, "daysOfWeek values must be between 0 and 6")
			return
		}
	}

	result, err := h.service.GenerateSlots(c.Request.Context(), ownerID, eventTypeID, req)
	if err != nil {
		if strings.Contains(err.Error(), "event type not found") {
			NotFound(c, "Event type")
			return
		}
		BadRequest(c, "Failed to generate slots: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, result)
}
