package repositories

import (
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookingRepository_CreateWithReservedSlotAndCancelLifecycle(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "booking")
	eventType := createTestEventType(t, ctx, owner.ID, "booking")
	startTime := time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC)
	slot := createTestTimeSlot(t, ctx, owner.ID, eventType.ID, startTime)
	timezone := "UTC"

	repo := NewBookingRepository()
	booking := &models.Booking{
		EventTypeID: eventType.ID,
		SlotID:      &slot.ID,
		GuestName:   "Jane Doe",
		GuestEmail:  "jane@example.com",
		Timezone:    &timezone,
		StartTime:   slot.StartTime,
		EndTime:     slot.EndTime,
	}

	require.NoError(t, repo.CreateWithReservedSlot(ctx, booking))
	require.NotEmpty(t, booking.ID)
	assert.Equal(t, "confirmed", booking.Status)

	stored, err := repo.GetByID(ctx, booking.ID)
	require.NoError(t, err)
	assert.Equal(t, "Jane Doe", stored.GuestName)

	reservedSlot, err := NewTimeSlotRepository().GetByID(ctx, slot.ID)
	require.NoError(t, err)
	assert.False(t, reservedSlot.IsAvailable)

	overlaps, err := repo.CheckOverlap(ctx, slot.StartTime, slot.EndTime)
	require.NoError(t, err)
	assert.True(t, overlaps)

	items, total, err := repo.List(ctx, 1, 10, "startTime", "asc", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, items, 1)

	require.NoError(t, repo.Cancel(ctx, booking.ID))
	items, total, err = repo.List(ctx, 1, 10, "startTime", "asc", nil, nil)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, items, 0)

	require.NoError(t, repo.Delete(ctx, booking.ID))
	_, err = repo.GetByID(ctx, booking.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "booking not found")
}

func TestBookingRepository_UpdateAndPatchLifecycle(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "booking-patch")
	eventType := createTestEventType(t, ctx, owner.ID, "booking-patch")
	repo := NewBookingRepository()
	initialTimezone := "UTC"

	booking := &models.Booking{
		EventTypeID: eventType.ID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
		Timezone:    &initialTimezone,
		StartTime:   time.Date(2026, 4, 21, 10, 0, 0, 0, time.UTC),
		EndTime:     time.Date(2026, 4, 21, 10, 30, 0, 0, time.UTC),
	}

	require.NoError(t, repo.Create(ctx, booking))

	booking.GuestName = "John Updated"
	updatedTimezone := "Europe/Moscow"
	booking.Timezone = &updatedTimezone
	require.NoError(t, repo.Update(ctx, booking))

	patchedName := "John Patched"
	patchedEmail := "patched@example.com"
	patched, err := repo.Patch(ctx, booking.ID, models.UpdateBookingRequest{
		GuestName:  &patchedName,
		GuestEmail: &patchedEmail,
	})
	require.NoError(t, err)
	assert.Equal(t, "John Patched", patched.GuestName)
	assert.Equal(t, "patched@example.com", patched.GuestEmail)
}
