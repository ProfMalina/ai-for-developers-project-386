package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/config"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/handlers"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	ctx := context.Background()
	if err := db.InitDB(ctx, cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.CloseDB()

	// Run migrations
	migrationsDir := "./migrations"
	if err := db.RunMigrations(ctx, migrationsDir); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
		log.Println("Continuing with existing schema...")
	}

	// Seed default owner for development
	if err := db.SeedDefaultOwner(ctx); err != nil {
		log.Printf("Warning: Failed to seed default owner: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	r := gin.Default()

	// Apply middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	// Initialize handlers
	// ownerHandler := handlers.NewOwnerHandler()
	eventTypeHandler := handlers.NewEventTypeHandler()
	timeSlotHandler := handlers.NewTimeSlotHandler()
	bookingHandler := handlers.NewBookingHandler()
	publicEventTypeHandler := handlers.NewPublicEventTypeHandler()
	publicBookingHandler := handlers.NewPublicBookingHandler()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Owner API routes
	api := r.Group("/api")
	{
		// Event types
		api.POST("/event-types", eventTypeHandler.Create)
		api.GET("/event-types", eventTypeHandler.List)
		api.GET("/event-types/:id", eventTypeHandler.GetByID)
		api.PATCH("/event-types/:id", eventTypeHandler.Update)
		api.DELETE("/event-types/:id", eventTypeHandler.Delete)

		// Slot generation
		api.POST("/slots/generate", timeSlotHandler.GenerateSlots)

		// Slots
		api.GET("/slots", timeSlotHandler.List)

		// Bookings
		api.GET("/bookings", bookingHandler.List)
		api.GET("/bookings/:id", bookingHandler.GetByID)
		api.DELETE("/bookings/:id", bookingHandler.Cancel)
	}

	// Guest (Public) API routes
	public := r.Group("/api/public")
	{
		// Event types
		public.GET("/event-types", publicEventTypeHandler.List)
		public.GET("/event-types/:id", publicEventTypeHandler.GetByID)
		public.GET("/slots", publicEventTypeHandler.GetSlots)

		// Bookings
		public.POST("/bookings", publicBookingHandler.Create)
	}

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
