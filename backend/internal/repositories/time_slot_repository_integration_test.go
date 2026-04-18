package repositories

import (
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeSlotRepository_CreateListAndMarkUnavailableLifecycle(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "slot")
	eventType := createTestEventType(t, ctx, owner.ID, "slot")
	repo := NewTimeSlotRepository()

	slot := &models.TimeSlot{
		OwnerID:     owner.ID,
		EventTypeID: eventType.ID,
		StartTime:   time.Date(2026, 4, 22, 9, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 4, 22, 9, 30, 0, 0, time.UTC),
		IsAvailable: true,
	}

	require.NoError(t, repo.Create(ctx, slot))
	require.NotEmpty(t, slot.ID)

	stored, err := repo.GetByID(ctx, slot.ID)
	require.NoError(t, err)
	assert.Equal(t, slot.ID, stored.ID)
	assert.Equal(t, eventType.ID, stored.EventTypeID)

	available := true
	items, total, err := repo.List(ctx, owner.ID, eventType.ID, 1, 10, &available, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, items, 1)

	require.NoError(t, repo.MarkAsUnavailable(ctx, slot.ID))
	stored, err = repo.GetByID(ctx, slot.ID)
	require.NoError(t, err)
	assert.False(t, stored.IsAvailable)

	require.NoError(t, repo.Delete(ctx, slot.ID))
	_, err = repo.GetByID(ctx, slot.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "time slot not found")
}

func TestTimeSlotRepository_DeleteAvailableInRangePreservesBookedSlots(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "slot-delete-range")
	eventType := createTestEventType(t, ctx, owner.ID, "slot-delete-range")

	startTime := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	bookedSlot := createTestTimeSlot(t, ctx, owner.ID, eventType.ID, startTime)
	availableSlot := createTestTimeSlot(t, ctx, owner.ID, eventType.ID, startTime.Add(time.Hour))

	bookingRepo := NewBookingRepository()
	booking := &models.Booking{
		EventTypeID: eventType.ID,
		SlotID:      &bookedSlot.ID,
		GuestName:   "Booked Guest",
		GuestEmail:  "booked@example.com",
		Timezone:    func() *string { tz := "UTC"; return &tz }(),
		StartTime:   bookedSlot.StartTime,
		EndTime:     bookedSlot.EndTime,
	}
	require.NoError(t, bookingRepo.CreateWithReservedSlot(ctx, booking))

	repo := NewTimeSlotRepository()
	require.NoError(t, repo.DeleteAvailableInRange(ctx, owner.ID, eventType.ID, startTime.Add(-time.Minute), startTime.Add(2*time.Hour)))

	_, err := repo.GetByID(ctx, bookedSlot.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, availableSlot.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "time slot not found")
}

func TestTimeSlotRepository_GetAvailableSlots_ExcludesRegeneratedOverlapWithActiveBooking(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "slot-public-overlap")
	eventType := createTestEventType(t, ctx, owner.ID, "slot-public-overlap")

	startTime := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	bookedSlot := createTestTimeSlot(t, ctx, owner.ID, eventType.ID, startTime)

	bookingRepo := NewBookingRepository()
	booking := &models.Booking{
		EventTypeID: eventType.ID,
		SlotID:      &bookedSlot.ID,
		GuestName:   "Booked Guest",
		GuestEmail:  "booked@example.com",
		Timezone:    func() *string { tz := "UTC"; return &tz }(),
		StartTime:   bookedSlot.StartTime,
		EndTime:     bookedSlot.EndTime,
	}
	require.NoError(t, bookingRepo.CreateWithReservedSlot(ctx, booking))

	repo := NewTimeSlotRepository()
	regeneratedSlot := &models.TimeSlot{
		OwnerID:     owner.ID,
		EventTypeID: eventType.ID,
		StartTime:   bookedSlot.StartTime,
		EndTime:     bookedSlot.EndTime,
		IsAvailable: true,
	}
	require.NoError(t, repo.Create(ctx, regeneratedSlot))
	require.NotEmpty(t, regeneratedSlot.ID)

	available := true
	allAvailable, availableTotal, err := repo.List(ctx, owner.ID, "", 1, 20, &available, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, availableTotal)
	require.Len(t, allAvailable, 1)
	assert.Equal(t, regeneratedSlot.ID, allAvailable[0].ID)

	items, total, err := repo.GetAvailableSlots(ctx, owner.ID, 1, 20, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, items)
}
