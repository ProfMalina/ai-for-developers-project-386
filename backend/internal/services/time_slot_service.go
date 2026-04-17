package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

// TimeSlotService handles business logic for time slots
type TimeSlotService struct {
	repo       TimeSlotRepository
	configRepo SlotGenerationConfigRepository
	ownerRepo  OwnerRepository
	etRepo     EventTypeRepository
}

// NewTimeSlotService creates a new time slot service
func NewTimeSlotService(repo TimeSlotRepository, configRepo SlotGenerationConfigRepository, ownerRepo OwnerRepository, etRepo EventTypeRepository) *TimeSlotService {
	return &TimeSlotService{
		repo:       repo,
		configRepo: configRepo,
		ownerRepo:  ownerRepo,
		etRepo:     etRepo,
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
func (s *TimeSlotService) List(ctx context.Context, ownerID, eventTypeID string, page, pageSize int, available *bool, startTime, endTime *time.Time) (*models.PaginatedResponse[models.TimeSlot], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	if eventTypeID != "" {
		if _, err := s.etRepo.GetByID(ctx, eventTypeID); err != nil {
			return nil, fmt.Errorf("event type not found: %w", err)
		}
	}

	items, totalItems, err := s.repo.List(ctx, ownerID, eventTypeID, page, pageSize, available, startTime, endTime)
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
func (s *TimeSlotService) GetAvailableSlots(ctx context.Context, ownerID string, page, pageSize int, startTime, endTime *time.Time) (*models.PaginatedResponse[models.TimeSlot], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, _, err := s.repo.GetAvailableSlots(ctx, ownerID, page, pageSize, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Filter out slots that have already started or passed
	now := time.Now()
	var filteredItems []models.TimeSlot
	for _, slot := range items {
		if slot.StartTime.After(now) {
			filteredItems = append(filteredItems, slot)
		}
	}

	// Recalculate pagination based on filtered results
	totalFiltered := len(filteredItems)
	pagination := CalculatePagination(page, pageSize, totalFiltered)

	return &models.PaginatedResponse[models.TimeSlot]{
		Items:      filteredItems,
		Pagination: pagination,
	}, nil
}

// GenerateSlots auto-generates time slots based on configuration
func (s *TimeSlotService) GenerateSlots(ctx context.Context, ownerID, eventTypeID string, req models.SlotGenerationRequest) (*models.SlotGenerationResult, error) {
	eventType, err := s.etRepo.GetByID(ctx, eventTypeID)
	if err != nil {
		return nil, fmt.Errorf("event type not found: %w", err)
	}

	// Always use owner's timezone from database for consistency
	owner, err := s.ownerRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("owner not found: %w", err)
	}

	timezoneName := owner.Timezone
	if req.Timezone != nil && *req.Timezone != "" {
		timezoneName = *req.Timezone
	}

	loc := time.UTC
	if timezoneName != "" {
		if l, err := time.LoadLocation(timezoneName); err == nil {
			loc = l
			fmt.Printf("Using timezone: %s\n", timezoneName)
		} else {
			fmt.Printf("Failed to load timezone %s: %v, falling back to UTC\n", timezoneName, err)
		}
	} else {
		fmt.Println("No timezone configured, using UTC")
	}
	// Parse dates
	nowInLocation := time.Now().In(loc)
	dateFrom := time.Date(nowInLocation.Year(), nowInLocation.Month(), nowInLocation.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, 1)
	if req.DateFrom != "" {
		var err error
		parsed, err := time.ParseInLocation("2006-01-02", req.DateFrom, loc)
		if err != nil {
			return nil, fmt.Errorf("invalid date_from format: %w", err)
		}
		if parsed.Before(dateFrom) {
			dateFrom = dateFrom
		} else {
			dateFrom = parsed
		}
	}

	dateTo := dateFrom.AddDate(0, 0, 30)
	if req.DateTo != "" {
		var err error
		dateTo, err = time.ParseInLocation("2006-01-02", req.DateTo, loc)
		if err != nil {
			return nil, fmt.Errorf("invalid date_to format: %w", err)
		}
	}
	if dateTo.Before(dateFrom) {
		return nil, fmt.Errorf("date_to must not be before date_from")
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

	intervalMinutes := req.IntervalMinutes
	if intervalMinutes == 0 {
		intervalMinutes = 30
	}

	daysOfWeek := req.DaysOfWeek
	if len(daysOfWeek) == 0 {
		daysOfWeek = []int{1, 2, 3, 4, 5}
	}

	// Create days of week map
	daysMap := make(map[int]bool)
	for _, day := range daysOfWeek {
		daysMap[day] = true
	}

	// Save config
	config := &models.SlotGenerationConfig{
		OwnerID:           ownerID,
		WorkingHoursStart: req.WorkingHoursStart,
		WorkingHoursEnd:   req.WorkingHoursEnd,
		IntervalMinutes:   intervalMinutes,
		DaysOfWeek:        daysOfWeek,
		DateFrom:          dateFrom,
		DateTo:            dateTo,
	}
	if req.Timezone != nil {
		config.Timezone = *req.Timezone
	} else {
		config.Timezone = timezoneName
	}

	if err := s.configRepo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to save slot generation config: %w", err)
	}

	slotsCreated := 0
	slotsSkipped := 0
	existingSlots, _, err := s.repo.List(ctx, ownerID, eventTypeID, 1, 100000, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing slots: %w", err)
	}
	existing := make(map[string]struct{}, len(existingSlots))
	for _, slot := range existingSlots {
		existing[slot.StartTime.UTC().Format(time.RFC3339)+"|"+slot.EndTime.UTC().Format(time.RFC3339)] = struct{}{}
	}

	// Calculate slots for each day
	currentDate := dateFrom
	for !currentDate.After(dateTo) {
		weekday := int(currentDate.Weekday())

		if daysMap[weekday] {
			// Generate slots for this day in the OWNER'S timezone
			slotStart := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				workStart.Hour(), workStart.Minute(), 0, 0, loc)

			dayEnd := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(),
				workEnd.Hour(), workEnd.Minute(), 0, 0, loc)

			for slotStart.Before(dayEnd) {
				slotEnd := slotStart.Add(time.Duration(eventType.DurationMinutes) * time.Minute)

				if slotEnd.After(dayEnd) {
					break
				}

				key := slotStart.UTC().Format(time.RFC3339) + "|" + slotEnd.UTC().Format(time.RFC3339)
				if _, exists := existing[key]; exists {
					slotsSkipped++
					slotStart = slotStart.Add(time.Duration(intervalMinutes) * time.Minute)
					continue
				}

				// Create the slot in the database (times are automatically converted to UTC)
				slot := &models.TimeSlot{
					OwnerID:     ownerID,
					EventTypeID: eventTypeID,
					StartTime:   slotStart.UTC(),
					EndTime:     slotEnd.UTC(),
					IsAvailable: true,
				}

				if err := s.repo.Create(ctx, slot); err != nil {
					return nil, fmt.Errorf("failed to create slot at %s: %w", slotStart.Format(time.RFC3339), err)
				}

				existing[key] = struct{}{}
				slotsCreated++

				slotStart = slotStart.Add(time.Duration(intervalMinutes) * time.Minute)
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return &models.SlotGenerationResult{
		SlotsCreated: slotsCreated,
		SlotsSkipped: slotsSkipped,
		DateFrom:     dateFrom.Format("2006-01-02"),
		DateTo:       dateTo.Format("2006-01-02"),
	}, nil
}
