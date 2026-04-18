package repositories

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	testDBOnce    sync.Once
	testDBInitErr error
)

func setupRepositoryTestDB(t *testing.T) context.Context {
	t.Helper()

	ctx := context.Background()
	testDBOnce.Do(func() {
		dsn := os.Getenv("TEST_DATABASE_URL")
		if dsn == "" {
			testDBInitErr = fmt.Errorf("TEST_DATABASE_URL is not set")
			return
		}

		testDBInitErr = db.InitDB(ctx, dsn)
		if testDBInitErr != nil {
			return
		}

		testDBInitErr = db.RunMigrations(ctx, repositoryMigrationsDir())
	})

	if testDBInitErr != nil {
		t.Skipf("repository integration database is unavailable: %v", testDBInitErr)
	}

	resetRepositoryTestDatabase(t, ctx)
	return ctx
}

func repositoryMigrationsDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return filepath.Join("..", "..", "migrations")
	}

	return filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations")
}

func resetRepositoryTestDatabase(t *testing.T, ctx context.Context) {
	t.Helper()

	_, err := db.Pool.Exec(ctx, `TRUNCATE TABLE bookings, time_slots, event_types, slot_generation_configs, owners RESTART IDENTITY CASCADE`)
	require.NoError(t, err)
	db.DefaultOwnerID = ""
}

func createTestOwner(t *testing.T, ctx context.Context, suffix string) *models.Owner {
	t.Helper()

	owner := &models.Owner{
		Name:     fmt.Sprintf("Owner %s", suffix),
		Email:    fmt.Sprintf("owner-%s@example.com", suffix),
		Timezone: "Europe/Moscow",
	}

	require.NoError(t, NewOwnerRepository().Create(ctx, owner))
	require.NotEmpty(t, owner.ID)
	return owner
}

func createTestEventType(t *testing.T, ctx context.Context, ownerID string, suffix string) *models.EventType {
	t.Helper()

	et := &models.EventType{
		OwnerID:         ownerID,
		Name:            fmt.Sprintf("Event %s", suffix),
		Description:     fmt.Sprintf("Description %s", suffix),
		DurationMinutes: 30,
	}

	require.NoError(t, NewEventTypeRepository().Create(ctx, et))
	require.NotEmpty(t, et.ID)
	return et
}

func createTestTimeSlot(t *testing.T, ctx context.Context, ownerID, eventTypeID string, startTime time.Time) *models.TimeSlot {
	t.Helper()

	slot := &models.TimeSlot{
		OwnerID:     ownerID,
		EventTypeID: eventTypeID,
		StartTime:   startTime.UTC(),
		EndTime:     startTime.Add(30 * time.Minute).UTC(),
		IsAvailable: true,
	}

	require.NoError(t, NewTimeSlotRepository().Create(ctx, slot))
	require.NotEmpty(t, slot.ID)
	return slot
}
