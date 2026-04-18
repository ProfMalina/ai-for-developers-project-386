package repositories

import (
	"testing"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOwnerRepository_CRUDLifecycle(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	repo := NewOwnerRepository()

	owner := &models.Owner{
		Name:     "Alice Example",
		Email:    "alice@example.com",
		Timezone: "Europe/Moscow",
	}

	require.NoError(t, repo.Create(ctx, owner))
	require.NotEmpty(t, owner.ID)
	assert.False(t, owner.CreatedAt.IsZero())
	assert.False(t, owner.UpdatedAt.IsZero())

	stored, err := repo.GetByID(ctx, owner.ID)
	require.NoError(t, err)
	assert.Equal(t, owner.Name, stored.Name)
	assert.Equal(t, owner.Email, stored.Email)

	owner.Name = "Alice Updated"
	owner.Timezone = "UTC"
	require.NoError(t, repo.Update(ctx, owner))
	assert.False(t, owner.UpdatedAt.IsZero())

	updated, err := repo.GetByID(ctx, owner.ID)
	require.NoError(t, err)
	assert.Equal(t, "Alice Updated", updated.Name)
	assert.Equal(t, "UTC", updated.Timezone)

	require.NoError(t, repo.Delete(ctx, owner.ID))

	_, err = repo.GetByID(ctx, owner.ID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "owner not found")
}

func TestOwnerRepository_UpdateAndDeleteNotFound(t *testing.T) {
	ctx := setupRepositoryTestDB(t)
	repo := NewOwnerRepository()

	err := repo.Update(ctx, &models.Owner{
		ID:       "00000000-0000-0000-0000-000000000099",
		Name:     "Missing",
		Email:    "missing@example.com",
		Timezone: "UTC",
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "owner not found")

	err = repo.Delete(ctx, "00000000-0000-0000-0000-000000000099")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "owner not found")
}
