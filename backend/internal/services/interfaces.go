package services

import (
	"context"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

// BookingRepository defines the interface for booking data operations
type BookingRepository interface {
	Create(ctx context.Context, booking *models.Booking) error
	GetByID(ctx context.Context, id string) (*models.Booking, error)
	List(ctx context.Context, page, pageSize int, sortBy, sortOrder string, status *string) ([]models.Booking, int, error)
	CheckOverlap(ctx context.Context, startTime, endTime time.Time) (bool, error)
	Cancel(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// TimeSlotRepository defines the interface for time slot data operations
type TimeSlotRepository interface {
	Create(ctx context.Context, slot *models.TimeSlot) error
	GetByID(ctx context.Context, id string) (*models.TimeSlot, error)
	List(ctx context.Context, ownerID string, page, pageSize int, available *bool, startTime, endTime *time.Time) ([]models.TimeSlot, int, error)
	GetAvailableSlots(ctx context.Context, ownerID string, page, pageSize int, startTime, endTime *time.Time) ([]models.TimeSlot, int, error)
	MarkAsUnavailable(ctx context.Context, slotID string) error
}

// EventTypeRepository defines the interface for event type data operations
type EventTypeRepository interface {
	Create(ctx context.Context, eventType *models.EventType) error
	GetByID(ctx context.Context, id string) (*models.EventType, error)
	List(ctx context.Context, ownerID string, page, pageSize int, sortBy, sortOrder string) ([]models.EventType, int, error)
	Patch(ctx context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error)
	Delete(ctx context.Context, id string) error
}

// OwnerRepository defines the interface for owner data operations
type OwnerRepository interface {
	Create(ctx context.Context, owner *models.Owner) error
	GetByID(ctx context.Context, id string) (*models.Owner, error)
}

// SlotGenerationConfigRepository defines the interface for slot generation config data operations
type SlotGenerationConfigRepository interface {
	Create(ctx context.Context, config *models.SlotGenerationConfig) error
}
