package memory

import (
	"context"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryBookingRepository_RejectsOverlappingActiveBookings(t *testing.T) {
	store := NewStore()
	repo := NewBookingRepository(store)

	start := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	end := start.Add(30 * time.Minute)

	require.NoError(t, repo.Create(context.Background(), &models.Booking{
		EventTypeID: "et-1",
		GuestName:   "First",
		GuestEmail:  "first@example.com",
		StartTime:   start,
		EndTime:     end,
		Status:      "confirmed",
	}))

	overlaps, err := repo.CheckOverlap(context.Background(), start, end)
	require.NoError(t, err)
	assert.True(t, overlaps)
}

func TestMemoryTimeSlotRepository_GetAvailableSlots_ExcludesBookedOverlap(t *testing.T) {
	store := NewStore()
	ownerID := SeedDefaultOwner(store)
	bookingRepo := NewBookingRepository(store)
	timeSlotRepo := NewTimeSlotRepository(store)

	start := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	slot := &models.TimeSlot{
		OwnerID:     ownerID,
		EventTypeID: "et-1",
		StartTime:   start,
		EndTime:     start.Add(30 * time.Minute),
		IsAvailable: true,
	}
	require.NoError(t, timeSlotRepo.Create(context.Background(), slot))

	require.NoError(t, bookingRepo.Create(context.Background(), &models.Booking{
		EventTypeID: "et-1",
		GuestName:   "Guest",
		GuestEmail:  "guest@example.com",
		StartTime:   slot.StartTime,
		EndTime:     slot.EndTime,
		Status:      "confirmed",
	}))

	items, total, err := timeSlotRepo.GetAvailableSlots(context.Background(), ownerID, 1, 20, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, items)
}

func TestMemorySlotGenerationConfigRepository_Create_UpsertsByOwner(t *testing.T) {
	store := NewStore()
	repo := NewSlotGenerationConfigRepository(store)
	ctx := context.Background()

	first := &models.SlotGenerationConfig{OwnerID: "owner-1", WorkingHoursStart: "09:00", WorkingHoursEnd: "18:00", IntervalMinutes: 30}
	require.NoError(t, repo.Create(ctx, first))
	require.NotEmpty(t, first.ID)

	second := &models.SlotGenerationConfig{OwnerID: "owner-1", WorkingHoursStart: "10:00", WorkingHoursEnd: "17:00", IntervalMinutes: 15}
	require.NoError(t, repo.Create(ctx, second))

	stored, err := repo.GetByOwnerID(ctx, "owner-1")
	require.NoError(t, err)
	assert.Equal(t, first.ID, stored.ID)
	assert.Equal(t, "10:00", stored.WorkingHoursStart)
	assert.Equal(t, 15, stored.IntervalMinutes)
	assert.Len(t, store.slotConfigs, 1)
}
