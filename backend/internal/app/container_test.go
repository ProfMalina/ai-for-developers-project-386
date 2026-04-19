package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContainer_FallsBackToMemoryWhenDatabaseInitFails(t *testing.T) {
	container, err := NewContainer(ContainerConfig{
		DatabaseURL: "postgres://invalid:invalid@127.0.0.1:1/booking_db?sslmode=disable",
	})
	require.NoError(t, err)
	assert.Equal(t, StorageModeMemory, container.Mode)
	require.NotNil(t, container.BookingService)
	require.NotNil(t, container.TimeSlotService)
	require.NotNil(t, container.PublicBookingHandler)
}

func TestNewContainer_UsesPostgresWhenDatabaseInitSucceeds(t *testing.T) {
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	container, err := NewContainer(ContainerConfig{DatabaseURL: testDBURL})
	require.NoError(t, err)
	assert.Equal(t, StorageModePostgres, container.Mode)
	require.NotNil(t, container.EventTypeHandler)
	require.NotNil(t, container.PublicEventTypeHandler)
}
