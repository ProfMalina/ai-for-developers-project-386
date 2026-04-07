package models

import "time"

// EventType represents a type of event that can be booked
type EventType struct {
	ID              string    `json:"id"`
	OwnerID         string    `json:"ownerId"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"durationMinutes"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CreateEventTypeRequest represents the request body for creating an event type
type CreateEventTypeRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	Description     string `json:"description" binding:"required,min=1,max=500"`
	DurationMinutes int    `json:"durationMinutes" binding:"required,min=5,max=1440"`
}

// UpdateEventTypeRequest represents the request body for updating an event type
// All fields are optional; only provided fields will be updated
type UpdateEventTypeRequest struct {
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	DurationMinutes *int    `json:"durationMinutes,omitempty"`
	IsActive        *bool   `json:"isActive,omitempty"`
}
