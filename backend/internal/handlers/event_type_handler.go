package handlers

import (
	"net/http"
	"strconv"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// EventTypeHandler handles HTTP requests for event types
type EventTypeHandler struct {
	service *services.EventTypeService
}

// NewEventTypeHandler creates a new event type handler
func NewEventTypeHandler() *EventTypeHandler {
	return &EventTypeHandler{
		service: services.NewEventTypeService(repositories.NewEventTypeRepository()),
	}
}

// Create handles POST /api/event-types
func (h *EventTypeHandler) Create(c *gin.Context) {
	ownerID := c.GetString("ownerID")
	if ownerID == "" {
		ownerID = db.DefaultOwnerID
	}

	var req models.CreateEventTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	et, err := h.service.Create(c.Request.Context(), ownerID, req)
	if err != nil {
		BadRequest(c, "Failed to create event type: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, et)
}

// GetByID handles GET /api/event-types/{id}
func (h *EventTypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	et, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "Event type")
		return
	}

	SuccessResponse(c, http.StatusOK, et)
}

// List handles GET /api/event-types
func (h *EventTypeHandler) List(c *gin.Context) {
	ownerID := c.GetString("ownerID")
	if ownerID == "" {
		ownerID = db.DefaultOwnerID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	result, err := h.service.List(c.Request.Context(), ownerID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		BadRequest(c, "Failed to list event types: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}

// Update handles PATCH /api/event-types/{id}
func (h *EventTypeHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateEventTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	et, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if err.Error() == "event type not found" {
			NotFound(c, "Event type")
			return
		}
		BadRequest(c, "Failed to update event type: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, et)
}

// Delete handles DELETE /api/event-types/{id}
func (h *EventTypeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		NotFound(c, "Event type")
		return
	}

	NoContent(c)
}
