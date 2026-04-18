package repositories

import (
	"testing"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventTypeRepository_CRUDAndListLifecycle(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "event-type")
	repo := NewEventTypeRepository()

	eventType := &models.EventType{
		OwnerID:         owner.ID,
		Name:            "Consultation",
		Description:     "Initial description",
		DurationMinutes: 45,
	}

	require.NoError(t, repo.Create(ctx, eventType))
	require.NotEmpty(t, eventType.ID)
	assert.True(t, eventType.IsActive)

	stored, err := repo.GetByID(ctx, eventType.ID)
	require.NoError(t, err)
	assert.Equal(t, "Consultation", stored.Name)
	assert.Equal(t, 45, stored.DurationMinutes)

	items, total, err := repo.List(ctx, owner.ID, 1, 10, "name", "asc")
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	require.Len(t, items, 1)
	assert.Equal(t, eventType.ID, items[0].ID)

	eventType.Name = "Deep Dive"
	eventType.Description = "Updated description"
	eventType.DurationMinutes = 60
	eventType.IsActive = false
	require.NoError(t, repo.Update(ctx, eventType))

	name := "Patched Deep Dive"
	active := true
	patched, err := repo.Patch(ctx, eventType.ID, models.UpdateEventTypeRequest{
		Name:     &name,
		IsActive: &active,
	})
	require.NoError(t, err)
	assert.Equal(t, "Patched Deep Dive", patched.Name)
	assert.True(t, patched.IsActive)
	assert.Equal(t, 60, patched.DurationMinutes)

	require.NoError(t, repo.Delete(ctx, eventType.ID))

	_, err = repo.GetByID(ctx, eventType.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "event type not found")
}

func TestEventTypeRepository_PatchWithoutFieldsReturnsExistingRecord(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	owner := createTestOwner(t, ctx, "event-type-empty-patch")
	repo := NewEventTypeRepository()
	eventType := createTestEventType(t, ctx, owner.ID, "empty-patch")

	patched, err := repo.Patch(ctx, eventType.ID, models.UpdateEventTypeRequest{})
	require.NoError(t, err)
	assert.Equal(t, eventType.ID, patched.ID)
	assert.Equal(t, eventType.Name, patched.Name)
	assert.True(t, patched.IsActive)
}

func TestEventTypeRepository_UpdateAndDeleteNotFound(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	repo := NewEventTypeRepository()

	err := repo.Update(ctx, &models.EventType{
		ID:              "00000000-0000-0000-0000-000000000199",
		Name:            "Missing",
		Description:     "Missing",
		DurationMinutes: 30,
		IsActive:        true,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "event type not found")

	err = repo.Delete(ctx, "00000000-0000-0000-0000-000000000199")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "event type not found")
}
