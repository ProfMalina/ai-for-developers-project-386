package services

import (
	"context"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

// EventTypeService handles business logic for event types
type EventTypeService struct {
	repo EventTypeRepository
}

// NewEventTypeService creates a new event type service
func NewEventTypeService(repo EventTypeRepository) *EventTypeService {
	return &EventTypeService{
		repo: repo,
	}
}

// Create creates a new event type
func (s *EventTypeService) Create(ctx context.Context, ownerID string, req models.CreateEventTypeRequest) (*models.EventType, error) {
	et := &models.EventType{
		OwnerID:         ownerID,
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
	}

	if err := s.repo.Create(ctx, et); err != nil {
		return nil, err
	}

	return et, nil
}

// GetByID retrieves an event type by ID
func (s *EventTypeService) GetByID(ctx context.Context, id string) (*models.EventType, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves a paginated list of event types
func (s *EventTypeService) List(ctx context.Context, ownerID string, page, pageSize int, sortBy, sortOrder string) (*models.PaginatedResponse[models.EventType], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, totalItems, err := s.repo.List(ctx, ownerID, page, pageSize, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	pagination := CalculatePagination(page, pageSize, totalItems)

	return &models.PaginatedResponse[models.EventType]{
		Items:      items,
		Pagination: pagination,
	}, nil
}

// Update updates an event type
func (s *EventTypeService) Update(ctx context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error) {
	return s.repo.Patch(ctx, id, req)
}

// Delete deletes an event type
func (s *EventTypeService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
