package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
)

// TimeSlotService handles business logic for time slots
type TimeSlotService struct {
	repo       *repositories.TimeSlotRepository
	configRepo *repositories.SlotGenerationConfigRepository
}

// NewTimeSlotService creates a new time slot service
func NewTimeSlotService() *TimeSlotService {
	return &TimeSlotService{
		repo:       repositories.NewTimeSlotRepository(),
		configRepo: repositories.NewSlotGenerationConfigRepository(),
	}
}

// Create creates a new time slot
func (s *TimeSlotService) Create(ctx context.Context, slot *models.TimeSlot) error {
	return s.repo.Create(ctx, slot)
}

// GetByID retrieves a time slot by ID
func (s *TimeSlotService) GetByID(ctx context.Context, id string) (*models.TimeSlot, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves a paginated list of time slots
func (s *TimeSlotService) List(ctx context.Context, eventTypeID string, page, pageSize int, available *bool, startTime, endTime *time.Time) (*models.PaginatedResponse[models.TimeSlot], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, totalItems, err := s.repo.List(ctx, eventTypeID, page, pageSize, available, startTime, endTime)
	if err != nil {
		return nil, err
	}

	pagination := CalculatePagination(page, pageSize, totalItems)

	return &models.PaginatedResponse[models.TimeSlot]{
		Items:      items,
		Pagination: pagination,
	}, nil
}

// GetAvailableSlots retrieves available slots for an event type
func (s *TimeSlotService) GetAvailableSlots(ctx context.Context, eventTypeID string, page, pageSize int, startTime, endTime *time.Time) (*models.PaginatedResponse[models.TimeSlot], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, totalItems, err := s.repo.GetAvailableSlots(ctx, eventTypeID, page, pageSize, startTime, endTime)
	if err != nil {
		return nil, err
	}

	pagination := CalculatePagination(page, pageSize, totalItems)

	return &models.PaginatedResponse[models.TimeSlot]{
		Items:      items,
		Pagination: pagination,
	}, nil
}

// GenerateSlots auto-generates time slots based on configuration
func (s *TimeSlotService) GenerateSlots(ctx context.Context, ownerID string, req models.SlotGenerationRequest) (*models.SlotGenerationResult, error) {
	// Parse dates
	dateFrom := time.Now().AddDate(0, 0, 1) // tomorrow
	if req.DateFrom != "" {
		var err error
		dateFrom, err = time.Parse("2006-01-02", req.DateFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from format: %w", err)
		}
	}

	dateTo := dateFrom.AddDate(0, 0, 30)
	if req.DateTo != "" {
		var err error
		dateTo, err = time.Parse("2006-01-02", req.DateTo)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to format: %w", err)
		}
	}

	// Parse working hours
	workStart, err := time.Parse("15:04", req.WorkingHoursStart)
	if err != nil {
		return nil, fmt.Errorf("invalid working_hours_start format: %w", err)
	}

	workEnd, err := time.Parse("15:04", req.WorkingHoursEnd)
	if err != nil {
		return nil, fmt.Errorf("invalid working_hours_end format: %w", err)
	}

	// Create days of week map
	daysMap := make(map[int]bool)
	for _, day := range req.DaysOfWeek {
		daysMap[day] = true
	}

	// Save config
	config := &models.SlotGenerationConfig{
		OwnerID:           ownerID,
		WorkingHoursStart: req.WorkingHoursStart,
		WorkingHoursEnd:   req.WorkingHoursEnd,
		IntervalMinutes:   req.IntervalMinutes,
		DaysOfWeek:        req.DaysOfWeek,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
	}
	if req.Timezone != nil {
		config.Timezone = *req.Timezone
	}

	if err := s.configRepo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to save slot generation config: %w", err)
	}

	slotsCreated := 0
	var createdSlotIDs []string

	// Calculate slots for each day
	currentDate := dateFrom
	for !currentDate.After(dateTo) {
		// Check if this day is in the days of week
		// Go's Weekday: Sunday=0, Monday=1, ..., Saturday=6
		// ISO: Monday=1, ..., Sunday=7
		goWeekday := int(currentDate.Weekday())
		isoWeekday := goWeekday
		if goWeekday == 0 {
			isoWeekday = 7 // Sunday
		}

		if daysMap[isoWeekday] {
			// Generate slots for this day
			slotStart := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				workStart.Hour(), workStart.Minute(), 0, 0, time.UTC)

			dayEnd := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				workEnd.Hour(), workEnd.Minute(), 0, 0, time.UTC)

			for slotStart.Before(dayEnd) {
				slotEnd := slotStart.Add(time.Duration(req.IntervalMinutes) * time.Minute)

				if slotEnd.After(dayEnd) {
					break
				}

				// Create the slot in the database
				slot := &models.TimeSlot{
					EventTypeID:   req.EventTypeID,
					StartTime:     slotStart,
					EndTime:       slotEnd,
					IsAvailable:   true,
				}

				if err := s.repo.Create(ctx, slot); err != nil {
					return nil, fmt.Errorf("failed to create slot at %s: %w", slotStart.Format(time.RFC3339), err)
				}

				createdSlotIDs = append(createdSlotIDs, slot.ID)
				slotsCreated++

				slotStart = slotEnd
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &models.SlotGenerationResult{
		SlotsCreated:  slotsCreated,
		CreatedSlotIDs: createdSlotIDs,
	}, nil
}
