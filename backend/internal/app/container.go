package app

import (
	"context"
	"log"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/handlers"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
	memrepo "github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories/memory"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
)

type ContainerConfig struct {
	DatabaseURL string
}

type Container struct {
	Mode                   StorageMode
	EventTypeHandler       *handlers.EventTypeHandler
	TimeSlotHandler        *handlers.TimeSlotHandler
	BookingHandler         *handlers.BookingHandler
	PublicEventTypeHandler *handlers.PublicEventTypeHandler
	PublicBookingHandler   *handlers.PublicBookingHandler
	OwnerHandler           *handlers.OwnerHandler
	ScheduleHandler        *handlers.ScheduleHandler
	BookingService         *services.BookingService
	TimeSlotService        *services.TimeSlotService
	EventTypeService       *services.EventTypeService
	OwnerService           *services.OwnerService
	ScheduleService        *services.ScheduleService
}

func NewContainer(cfg ContainerConfig) (*Container, error) {
	ctx := context.Background()
	if err := db.InitDB(ctx, cfg.DatabaseURL); err == nil {
		eventTypeService := services.NewEventTypeService(repositories.NewEventTypeRepository())
		timeSlotService := services.NewTimeSlotService(
			repositories.NewTimeSlotRepository(),
			repositories.NewSlotGenerationConfigRepository(),
			repositories.NewOwnerRepository(),
			repositories.NewEventTypeRepository(),
		)
		bookingService := services.NewBookingService(
			repositories.NewBookingRepository(),
			repositories.NewTimeSlotRepository(),
			repositories.NewEventTypeRepository(),
		)
		ownerService := services.NewOwnerService(repositories.NewOwnerRepository())
		scheduleService := services.NewScheduleService(repositories.NewScheduleRepository())

		return &Container{
			Mode:                   StorageModePostgres,
			EventTypeService:       eventTypeService,
			TimeSlotService:        timeSlotService,
			BookingService:         bookingService,
			OwnerService:           ownerService,
			ScheduleService:        scheduleService,
			EventTypeHandler:       handlers.NewEventTypeHandlerWithService(eventTypeService),
			TimeSlotHandler:        handlers.NewTimeSlotHandlerWithService(timeSlotService),
			BookingHandler:         handlers.NewBookingHandlerWithService(bookingService),
			PublicEventTypeHandler: handlers.NewPublicEventTypeHandlerWithServices(eventTypeService, timeSlotService),
			PublicBookingHandler:   handlers.NewPublicBookingHandlerWithService(bookingService),
			OwnerHandler:           handlers.NewOwnerHandlerWithService(ownerService),
			ScheduleHandler:        handlers.NewScheduleHandlerWithService(scheduleService),
		}, nil
	}

	db.CloseDB()
	log.Printf("Database unavailable, starting in in-memory mode")
	store := memrepo.NewStore()
	defaultOwnerID := memrepo.SeedDefaultOwner(store)
	db.DefaultOwnerID = defaultOwnerID

	eventTypeService := services.NewEventTypeService(memrepo.NewEventTypeRepository(store))
	timeSlotService := services.NewTimeSlotService(
		memrepo.NewTimeSlotRepository(store),
		memrepo.NewSlotGenerationConfigRepository(store),
		memrepo.NewOwnerRepository(store),
		memrepo.NewEventTypeRepository(store),
	)
	bookingService := services.NewBookingService(
		memrepo.NewBookingRepository(store),
		memrepo.NewTimeSlotRepository(store),
		memrepo.NewEventTypeRepository(store),
	)
	ownerService := services.NewOwnerService(memrepo.NewOwnerRepository(store))
	scheduleService := services.NewScheduleService(memrepo.NewScheduleRepository(store))

	return &Container{
		Mode:                   StorageModeMemory,
		EventTypeService:       eventTypeService,
		TimeSlotService:        timeSlotService,
		BookingService:         bookingService,
		OwnerService:           ownerService,
		ScheduleService:        scheduleService,
		EventTypeHandler:       handlers.NewEventTypeHandlerWithService(eventTypeService),
		TimeSlotHandler:        handlers.NewTimeSlotHandlerWithService(timeSlotService),
		BookingHandler:         handlers.NewBookingHandlerWithService(bookingService),
		PublicEventTypeHandler: handlers.NewPublicEventTypeHandlerWithServices(eventTypeService, timeSlotService),
		PublicBookingHandler:   handlers.NewPublicBookingHandlerWithService(bookingService),
		OwnerHandler:           handlers.NewOwnerHandlerWithService(ownerService),
		ScheduleHandler:        handlers.NewScheduleHandlerWithService(scheduleService),
	}, nil
}
