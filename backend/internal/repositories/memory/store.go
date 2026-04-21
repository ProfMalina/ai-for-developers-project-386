package memory

import (
	"sync"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

type Store struct {
	mu           sync.RWMutex
	owners       map[string]models.Owner
	eventTypes   map[string]models.EventType
	timeSlots    map[string]models.TimeSlot
	bookings     map[string]models.Booking
	slotConfigs  map[string]models.SlotGenerationConfig
	daySchedules map[string][]models.DaySchedule
	exceptions   map[string][]models.DateException
}

func NewStore() *Store {
	return &Store{
		owners:       make(map[string]models.Owner),
		eventTypes:   make(map[string]models.EventType),
		timeSlots:    make(map[string]models.TimeSlot),
		bookings:     make(map[string]models.Booking),
		slotConfigs:  make(map[string]models.SlotGenerationConfig),
		daySchedules: make(map[string][]models.DaySchedule),
		exceptions:   make(map[string][]models.DateException),
	}
}
