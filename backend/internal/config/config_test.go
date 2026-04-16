package config

import "testing"

func TestLoadConfigUsesPortEnvVar(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("SERVER_PORT", "8080")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.ServerPort != "9090" {
		t.Fatalf("expected ServerPort to use PORT env var, got %q", cfg.ServerPort)
	}
}

func TestLoadConfigFallsBackToServerPort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("SERVER_PORT", "7070")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if cfg.ServerPort != "7070" {
		t.Fatalf("expected ServerPort to fall back to SERVER_PORT, got %q", cfg.ServerPort)
	}
}
