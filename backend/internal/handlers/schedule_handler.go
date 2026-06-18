package handlers

import (
	"net/http"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service *services.ScheduleService
}

func NewScheduleHandlerWithService(service *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	ownerID := c.GetString("ownerId")
	if ownerID == "" {
		ownerID = "default-owner"
	}
	schedule, err := h.service.GetSchedule(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) UpsertSchedule(c *gin.Context) {
	ownerID := c.GetString("ownerId")
	if ownerID == "" {
		ownerID = "default-owner"
	}
	var req models.UpsertScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": err.Error()})
		return
	}
	schedule, err := h.service.UpsertSchedule(c.Request.Context(), ownerID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) UpsertDaySchedules(c *gin.Context) {
	ownerID := c.GetString("ownerId")
	if ownerID == "" {
		ownerID = "default-owner"
	}
	var schedules []models.DaySchedule
	if err := c.ShouldBindJSON(&schedules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": err.Error()})
		return
	}
	schedule, err := h.service.UpsertDaySchedules(c.Request.Context(), ownerID, schedules)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) UpsertDateExceptions(c *gin.Context) {
	ownerID := c.GetString("ownerId")
	if ownerID == "" {
		ownerID = "default-owner"
	}
	var exceptions []models.DateException
	if err := c.ShouldBindJSON(&exceptions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": err.Error()})
		return
	}
	schedule, err := h.service.UpsertDateExceptions(c.Request.Context(), ownerID, exceptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

func (h *ScheduleHandler) DeleteDateException(c *gin.Context) {
	ownerID := c.GetString("ownerId")
	if ownerID == "" {
		ownerID = "default-owner"
	}
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "BAD_REQUEST", "message": "date is required"})
		return
	}
	if err := h.service.DeleteDateException(c.Request.Context(), ownerID, date); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
