package app

import (
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(container *Container, env string) *gin.Engine {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",
			"time":         time.Now().Format(time.RFC3339),
			"storageMode":  container.Mode,
			"degradedMode": container.Mode == StorageModeMemory,
		})
	})

	api := r.Group("/api")
	{
		api.POST("/event-types", container.EventTypeHandler.Create)
		api.GET("/event-types", container.EventTypeHandler.List)
		api.GET("/event-types/:id", container.EventTypeHandler.GetByID)
		api.PATCH("/event-types/:id", container.EventTypeHandler.Update)
		api.DELETE("/event-types/:id", container.EventTypeHandler.Delete)
		api.POST("/event-types/:eventTypeId/slots/generate", container.TimeSlotHandler.GenerateSlots)
		api.GET("/slots", container.TimeSlotHandler.List)
		api.GET("/bookings", container.BookingHandler.List)
		api.GET("/bookings/:id", container.BookingHandler.GetByID)
		api.DELETE("/bookings/:id", container.BookingHandler.Cancel)
	}

	public := r.Group("/api/public")
	{
		public.GET("/event-types", container.PublicEventTypeHandler.List)
		public.GET("/event-types/:id", container.PublicEventTypeHandler.GetByID)
		public.GET("/slots", container.PublicEventTypeHandler.GetSlots)
		public.POST("/bookings", container.PublicBookingHandler.Create)
	}

	return r
}
