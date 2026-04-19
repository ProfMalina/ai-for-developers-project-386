package memory

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
)

type BookingRepository struct{ store *Store }

func NewBookingRepository(store *Store) *BookingRepository { return &BookingRepository{store: store} }

func (r *BookingRepository) Create(_ context.Context, booking *models.Booking) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if booking.ID == "" {
		booking.ID = uuid.NewString()
	}
	if booking.Status == "" {
		booking.Status = "confirmed"
	}
	booking.CreatedAt = time.Now().UTC()
	r.store.bookings[booking.ID] = *booking
	return nil
}

func (r *BookingRepository) CreateWithReservedSlot(_ context.Context, booking *models.Booking) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if booking.SlotID != nil {
		slot, ok := r.store.timeSlots[*booking.SlotID]
		if !ok {
			return fmt.Errorf("time slot not found")
		}
		if !slot.IsAvailable {
			return fmt.Errorf("time slot is already booked")
		}
		for _, existing := range r.store.bookings {
			if existing.Status == "cancelled" { //nolint:misspell // persisted booking status value
				continue
			}
			if existing.StartTime.Before(booking.EndTime) && existing.EndTime.After(booking.StartTime) {
				return fmt.Errorf("selected time slot is already booked")
			}
		}
		slot.IsAvailable = false
		r.store.timeSlots[*booking.SlotID] = slot
	}
	if booking.ID == "" {
		booking.ID = uuid.NewString()
	}
	if booking.Status == "" {
		booking.Status = "confirmed"
	}
	booking.CreatedAt = time.Now().UTC()
	r.store.bookings[booking.ID] = *booking
	return nil
}

func (r *BookingRepository) GetByID(_ context.Context, id string) (*models.Booking, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	b, ok := r.store.bookings[id]
	if !ok {
		return nil, fmt.Errorf("booking not found")
	}
	copy := b
	return &copy, nil
}

func (r *BookingRepository) List(_ context.Context, page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) ([]models.Booking, int, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	items := make([]models.Booking, 0)
	for _, booking := range r.store.bookings {
		if booking.Status == "cancelled" { //nolint:misspell // persisted booking status value
			continue
		}
		if dateFrom != nil && booking.StartTime.Before(*dateFrom) {
			continue
		}
		if dateTo != nil && booking.EndTime.After(*dateTo) {
			continue
		}
		items = append(items, booking)
	}
	sort.Slice(items, func(i, j int) bool {
		desc := sortOrder == "desc"
		less := items[i].StartTime.Before(items[j].StartTime)
		if sortBy == "created_at" || sortBy == "createdAt" {
			less = items[i].CreatedAt.Before(items[j].CreatedAt)
		}
		if sortBy == "guest_name" || sortBy == "guestName" {
			less = items[i].GuestName < items[j].GuestName
		}
		if sortBy == "status" {
			less = items[i].Status < items[j].Status
		}
		if desc {
			return !less
		}
		return less
	})
	total := len(items)
	start := (page - 1) * pageSize
	if start >= total {
		return []models.Booking{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return append([]models.Booking(nil), items[start:end]...), total, nil
}

func (r *BookingRepository) CheckOverlap(_ context.Context, startTime, endTime time.Time) (bool, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, booking := range r.store.bookings {
		if booking.Status == "cancelled" { //nolint:misspell // persisted booking status value
			continue
		}
		if booking.StartTime.Before(endTime) && booking.EndTime.After(startTime) {
			return true, nil
		}
	}
	return false, nil
}

func (r *BookingRepository) Cancel(_ context.Context, id string) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	b, ok := r.store.bookings[id]
	if !ok {
		return fmt.Errorf("booking not found")
	}
	b.Status = "cancelled" //nolint:misspell // persisted booking status value
	r.store.bookings[id] = b
	return nil
}

func (r *BookingRepository) Delete(_ context.Context, id string) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if _, ok := r.store.bookings[id]; !ok {
		return fmt.Errorf("booking not found")
	}
	delete(r.store.bookings, id)
	return nil
}
