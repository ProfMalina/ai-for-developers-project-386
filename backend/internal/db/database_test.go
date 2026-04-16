package db

import (
	"context"
	"testing"
)

func TestGetEnvIntReturnsParsedValue(t *testing.T) {
	t.Setenv("TEST_INT_ENV", "42")

	if got := getEnvInt("TEST_INT_ENV", 10); got != 42 {
		t.Fatalf("expected parsed value 42, got %d", got)
	}
}

func TestGetEnvIntFallsBackForInvalidValue(t *testing.T) {
	t.Setenv("TEST_INT_ENV", "invalid")

	if got := getEnvInt("TEST_INT_ENV", 10); got != 10 {
		t.Fatalf("expected default value 10, got %d", got)
	}
}

func TestRunMigrationsRequiresInitializedDatabase(t *testing.T) {
	Pool = nil

	err := RunMigrations(context.Background(), "./migrations")
	if err == nil {
		t.Fatal("expected error when database pool is nil")
	}
}

func TestSeedDefaultOwnerRequiresInitializedDatabase(t *testing.T) {
	Pool = nil

	err := SeedDefaultOwner(context.Background())
	if err == nil {
		t.Fatal("expected error when database pool is nil")
	}
}

func TestCloseDBWithNilPoolDoesNotPanic(t *testing.T) {
	Pool = nil
	CloseDB()
}
