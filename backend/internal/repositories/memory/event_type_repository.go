package memory

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
)

type EventTypeRepository struct{ store *Store }

func NewEventTypeRepository(store *Store) *EventTypeRepository {
	return &EventTypeRepository{store: store}
}

func (r *EventTypeRepository) Create(_ context.Context, et *models.EventType) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if et.ID == "" {
		et.ID = uuid.NewString()
	}
	if !et.IsActive {
		et.IsActive = true
	}
	now := time.Now().UTC()
	et.CreatedAt = now
	et.UpdatedAt = now
	r.store.eventTypes[et.ID] = *et
	return nil
}

func (r *EventTypeRepository) GetByID(_ context.Context, id string) (*models.EventType, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	et, ok := r.store.eventTypes[id]
	if !ok {
		return nil, fmt.Errorf("event type not found")
	}
	copy := et
	return &copy, nil
}

func (r *EventTypeRepository) List(_ context.Context, ownerID string, page, pageSize int, sortBy, sortOrder string) ([]models.EventType, int, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	items := make([]models.EventType, 0)
	for _, et := range r.store.eventTypes {
		if ownerID == "" || et.OwnerID == ownerID {
			items = append(items, et)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		desc := strings.EqualFold(sortOrder, "desc")
		less := items[i].CreatedAt.Before(items[j].CreatedAt)
		if sortBy == "name" {
			less = items[i].Name < items[j].Name
		}
		if desc {
			return !less
		}
		return less
	})
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []models.EventType{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return append([]models.EventType(nil), items[start:end]...), total, nil
}

func (r *EventTypeRepository) Patch(_ context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	et, ok := r.store.eventTypes[id]
	if !ok {
		return nil, fmt.Errorf("event type not found")
	}
	if req.Name != nil {
		et.Name = *req.Name
	}
	if req.Description != nil {
		et.Description = *req.Description
	}
	if req.DurationMinutes != nil {
		et.DurationMinutes = *req.DurationMinutes
	}
	if req.IsActive != nil {
		et.IsActive = *req.IsActive
	}
	et.UpdatedAt = time.Now().UTC()
	r.store.eventTypes[id] = et
	copy := et
	return &copy, nil
}

func (r *EventTypeRepository) Delete(_ context.Context, id string) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if _, ok := r.store.eventTypes[id]; !ok {
		return fmt.Errorf("event type not found")
	}
	delete(r.store.eventTypes, id)
	return nil
}
