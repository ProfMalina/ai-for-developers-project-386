package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/app"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/config"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	container, err := app.NewContainer(app.ContainerConfig{DatabaseURL: cfg.DatabaseURL})
	if err != nil {
		log.Fatalf("Failed to initialize application container: %v", err)
	}
	if container.Mode == app.StorageModePostgres {
		defer db.CloseDB()
		migrationsDir := "./migrations"
		if err := db.RunMigrations(ctx, migrationsDir); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
			log.Println("Continuing with existing schema...")
		}
		if err := db.SeedDefaultOwner(ctx); err != nil {
			log.Printf("Warning: Failed to seed default owner: %v", err)
		}
	}

	r := app.NewRouter(container, cfg.Env)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)

	// Create a channel to listen for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Give the server 5 seconds to gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()
	log.Println("Server stopped")
}
