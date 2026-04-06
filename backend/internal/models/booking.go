package models

import "time"

// Booking represents a confirmed booking
type Booking struct {
	ID          string    `json:"id"`
	EventTypeID string    `json:"eventTypeId"`
	SlotID      *string   `json:"slotId,omitempty"`
	GuestName   string    `json:"guestName"`
	GuestEmail  string    `json:"guestEmail"`
	Timezone    *string   `json:"timezone,omitempty"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CreateBookingRequest represents the request body for creating a booking
type CreateBookingRequest struct {
	EventTypeID string  `json:"eventTypeId" binding:"required,uuid"`
	SlotID      *string `json:"slotId,omitempty" binding:"omitempty,uuid"`
	GuestName   string  `json:"guestName" binding:"required,min=1,max=100"`
	GuestEmail  string  `json:"guestEmail" binding:"required,email,max=255"`
	Timezone    *string `json:"timezone,omitempty" binding:"omitempty,max=50"`
}

// UpdateBookingRequest represents the request body for updating a booking
type UpdateBookingRequest struct {
	GuestName *string `json:"guestName,omitempty" binding:"omitempty,min=1,max=100"`
	GuestEmail *string `json:"guestEmail,omitempty" binding:"omitempty,email,max=255"`
	Timezone  *string `json:"timezone,omitempty" binding:"omitempty,max=50"`
}
