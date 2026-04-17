package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

// BookingService handles business logic for bookings
type BookingService struct {
	repo     BookingRepository
	slotRepo TimeSlotRepository
	etRepo   EventTypeRepository
}

// NewBookingService creates a new booking service
func NewBookingService(repo BookingRepository, slotRepo TimeSlotRepository, etRepo EventTypeRepository) *BookingService {
	return &BookingService{
		repo:     repo,
		slotRepo: slotRepo,
		etRepo:   etRepo,
	}
}

// Create creates a new booking with conflict checking
func (s *BookingService) Create(ctx context.Context, req models.CreateBookingRequest) (*models.Booking, error) {
	// Verify event type exists
	_, err := s.etRepo.GetByID(ctx, req.EventTypeID)
	if err != nil {
		return nil, fmt.Errorf("event type not found: %w", err)
	}

	// Calculate start and end times
	// For now, we'll use current time + duration as a placeholder
	// In reality, this should be based on the selected slot
	booking := &models.Booking{
		EventTypeID: req.EventTypeID,
		SlotID:      req.SlotID,
		GuestName:   req.GuestName,
		GuestEmail:  req.GuestEmail,
		Timezone:    req.Timezone,
		Status:      "confirmed",
	}

	// If slot is provided, use its times
	if req.SlotID != nil {
		slot, err := s.slotRepo.GetByID(ctx, *req.SlotID)
		if err != nil {
			return nil, fmt.Errorf("time slot not found: %w", err)
		}

		if !slot.IsAvailable {
			return nil, fmt.Errorf("time slot is already booked")
		}

		// Prevent booking slots that have already started or passed
		if slot.StartTime.Before(time.Now()) {
			return nil, fmt.Errorf("cannot book a time slot that has already started or passed")
		}

		booking.StartTime = slot.StartTime
		booking.EndTime = slot.EndTime
	} else {
		// This should not happen - slot should always be selected
		return nil, fmt.Errorf("slot ID is required")
	}

	// Check for overlapping bookings
	overlaps, err := s.repo.CheckOverlap(ctx, booking.StartTime, booking.EndTime)
	if err != nil {
		return nil, err
	}

	if overlaps {
		return nil, fmt.Errorf("selected time slot is already booked")
	}

	// Create the booking
	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, err
	}

	// Mark the slot as unavailable
	if booking.SlotID != nil {
		if err := s.slotRepo.MarkAsUnavailable(ctx, *booking.SlotID); err != nil {
			// Log error but don't fail the booking
			fmt.Printf("Warning: failed to mark slot as unavailable: %v\n", err)
		}
	}

	return booking, nil
}

// GetByID retrieves a booking by ID
func (s *BookingService) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves a paginated list of bookings
func (s *BookingService) List(ctx context.Context, page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) (*models.PaginatedResponse[models.Booking], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, totalItems, err := s.repo.List(ctx, page, pageSize, sortBy, sortOrder, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	pagination := CalculatePagination(page, pageSize, totalItems)

	return &models.PaginatedResponse[models.Booking]{
		Items:      items,
		Pagination: pagination,
	}, nil
}

// Cancel cancels a booking
func (s *BookingService) Cancel(ctx context.Context, id string) error {
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if booking.Status == "cancelled" {
		return fmt.Errorf("booking is already cancelled")
	}

	if err := s.repo.Cancel(ctx, id); err != nil {
		return err
	}

	// Free up the slot if it was assigned
	if booking.SlotID != nil {
		// In a real implementation, you'd mark the slot as available again
		// For now, we'll leave it as is
	}

	return nil
}

// Delete deletes a booking
func (s *BookingService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
