package repositories

import (
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlotGenerationConfigRepository_CreateAndUpsertByOwner(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "slot-config")
	repo := NewSlotGenerationConfigRepository()

	config := &models.SlotGenerationConfig{
		OwnerID:           owner.ID,
		WorkingHoursStart: "09:00",
		WorkingHoursEnd:   "17:00",
		IntervalMinutes:   30,
		DaysOfWeek:        []int{1, 2, 3, 4, 5},
		DateFrom:          time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
		DateTo:            time.Date(2026, 4, 27, 0, 0, 0, 0, time.UTC),
		Timezone:          "Europe/Moscow",
	}

	require.NoError(t, repo.Create(ctx, config))
	require.NotEmpty(t, config.ID)

	stored, err := repo.GetByOwnerID(ctx, owner.ID)
	require.NoError(t, err)
	assert.Equal(t, 30, stored.IntervalMinutes)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, stored.DaysOfWeek)

	updated := &models.SlotGenerationConfig{
		ID:                "00000000-0000-0000-0000-000000000555",
		OwnerID:           owner.ID,
		WorkingHoursStart: "10:00",
		WorkingHoursEnd:   "18:00",
		IntervalMinutes:   15,
		DaysOfWeek:        []int{1, 3, 5},
		DateFrom:          time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		DateTo:            time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
		Timezone:          "UTC",
	}

	require.NoError(t, repo.Create(ctx, updated))

	stored, err = repo.GetByOwnerID(ctx, owner.ID)
	require.NoError(t, err)
	assert.Equal(t, config.ID, stored.ID)
	assert.Equal(t, "10:00:00", stored.WorkingHoursStart)
	assert.Equal(t, "18:00:00", stored.WorkingHoursEnd)
	assert.Equal(t, 15, stored.IntervalMinutes)
	assert.Equal(t, []int{1, 3, 5}, stored.DaysOfWeek)
	assert.Equal(t, "UTC", stored.Timezone)
	assert.True(t, stored.UpdatedAt.After(stored.CreatedAt) || stored.UpdatedAt.Equal(stored.CreatedAt))
}

func TestSlotGenerationConfigRepository_GetByOwnerIDNotFound(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	repo := NewSlotGenerationConfigRepository()

	_, err := repo.GetByOwnerID(ctx, "00000000-0000-0000-0000-000000000299")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "slot generation config not found")
}
