package models

import "time"

// TimeSlot represents an available time slot for booking
type TimeSlot struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"ownerId"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	IsAvailable bool      `json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CreateTimeSlotRequest represents the request body for creating a time slot
type CreateTimeSlotRequest struct {
	EventTypeID string    `json:"eventTypeId" binding:"required,uuid"`
	StartTime   time.Time `json:"startTime" binding:"required"`
	EndTime     time.Time `json:"endTime" binding:"required"`
}

// SlotGenerationConfig represents configuration for auto-generating slots
type SlotGenerationConfig struct {
	ID                string    `json:"id"`
	OwnerID           string    `json:"ownerId"`
	WorkingHoursStart string    `json:"workingHoursStart"`
	WorkingHoursEnd   string    `json:"workingHoursEnd"`
	IntervalMinutes   int       `json:"intervalMinutes"`
	DaysOfWeek        []int     `json:"daysOfWeek"`
	DateFrom          time.Time `json:"dateFrom"`
	DateTo            time.Time `json:"dateTo"`
	Timezone          string    `json:"timezone"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// SlotGenerationRequest represents the request body for generating slots
type SlotGenerationRequest struct {
	WorkingHoursStart string  `json:"workingHoursStart" binding:"required,datetime=15:04"`
	WorkingHoursEnd   string  `json:"workingHoursEnd" binding:"required,datetime=15:04"`
	IntervalMinutes   int     `json:"intervalMinutes" binding:"required,oneof=15 30"`
	DaysOfWeek        []int   `json:"daysOfWeek" binding:"required"`
	DateFrom          string  `json:"dateFrom" binding:"omitempty,datetime=2006-01-02"`
	DateTo            string  `json:"dateTo" binding:"omitempty,datetime=2006-01-02"`
	Timezone          *string `json:"timezone,omitempty"`
}

// SlotGenerationResult represents the response from slot generation
type SlotGenerationResult struct {
	SlotsCreated   int      `json:"slotsCreated"`
	CreatedSlotIDs []string `json:"createdSlotIds"`
}
