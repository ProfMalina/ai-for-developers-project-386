package memory

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
)

type TimeSlotRepository struct{ store *Store }

func NewTimeSlotRepository(store *Store) *TimeSlotRepository {
	return &TimeSlotRepository{store: store}
}

func (r *TimeSlotRepository) Create(_ context.Context, slot *models.TimeSlot) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if slot.ID == "" {
		slot.ID = uuid.NewString()
	}
	slot.CreatedAt = time.Now().UTC()
	r.store.timeSlots[slot.ID] = *slot
	return nil
}

func (r *TimeSlotRepository) GetByID(_ context.Context, id string) (*models.TimeSlot, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	slot, ok := r.store.timeSlots[id]
	if !ok {
		return nil, fmt.Errorf("time slot not found")
	}
	copy := slot
	return &copy, nil
}

func (r *TimeSlotRepository) List(_ context.Context, ownerID, eventTypeID string, page, pageSize int, available *bool, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	items := make([]models.TimeSlot, 0)
	for _, slot := range r.store.timeSlots {
		if ownerID != "" && slot.OwnerID != ownerID {
			continue
		}
		if eventTypeID != "" && slot.EventTypeID != eventTypeID {
			continue
		}
		if available != nil && slot.IsAvailable != *available {
			continue
		}
		if startTime != nil && slot.StartTime.Before(*startTime) {
			continue
		}
		if endTime != nil && slot.EndTime.After(*endTime) {
			continue
		}
		items = append(items, slot)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].StartTime.Before(items[j].StartTime) })
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []models.TimeSlot{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return append([]models.TimeSlot(nil), items[start:end]...), total, nil
}

func (r *TimeSlotRepository) GetAvailableSlots(_ context.Context, ownerID string, page, pageSize int, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	items := make([]models.TimeSlot, 0)
	for _, slot := range r.store.timeSlots {
		if slot.OwnerID != ownerID || !slot.IsAvailable {
			continue
		}
		if startTime != nil && slot.StartTime.Before(*startTime) {
			continue
		}
		if endTime != nil && slot.EndTime.After(*endTime) {
			continue
		}
		overlaps := false
		for _, booking := range r.store.bookings {
			if booking.Status == "cancelled" { //nolint:misspell // persisted booking status value
				continue
			}
			if booking.StartTime.Before(slot.EndTime) && booking.EndTime.After(slot.StartTime) {
				overlaps = true
				break
			}
		}
		if overlaps {
			continue
		}
		items = append(items, slot)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].StartTime.Before(items[j].StartTime) })
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []models.TimeSlot{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return append([]models.TimeSlot(nil), items[start:end]...), total, nil
}

func (r *TimeSlotRepository) DeleteAvailableInRange(_ context.Context, ownerID, eventTypeID string, startTime, endTime time.Time) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	for id, slot := range r.store.timeSlots {
		if slot.OwnerID != ownerID || !slot.IsAvailable {
			continue
		}
		if slot.EventTypeID != eventTypeID && slot.EventTypeID != "" {
			continue
		}
		if !slot.StartTime.Before(endTime) || !slot.EndTime.After(startTime) {
			continue
		}
		preserve := false
		for _, booking := range r.store.bookings {
			if booking.SlotID != nil && *booking.SlotID == id && booking.Status != "cancelled" { //nolint:misspell // persisted booking status value
				preserve = true
				break
			}
		}
		if !preserve {
			delete(r.store.timeSlots, id)
		}
	}
	return nil
}

func (r *TimeSlotRepository) MarkAsUnavailable(_ context.Context, slotID string) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	slot, ok := r.store.timeSlots[slotID]
	if !ok {
		return fmt.Errorf("time slot not found")
	}
	slot.IsAvailable = false
	r.store.timeSlots[slotID] = slot
	return nil
}
