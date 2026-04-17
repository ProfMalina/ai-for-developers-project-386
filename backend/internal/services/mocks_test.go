package services

import (
	"context"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockBookingRepository is a mock implementation of booking repository
type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockBookingRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockBookingRepository) List(ctx context.Context, page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) ([]models.Booking, int, error) {
	args := m.Called(ctx, page, pageSize, sortBy, sortOrder, dateFrom, dateTo)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Booking), args.Int(1), args.Error(2)
}

func (m *MockBookingRepository) CheckOverlap(ctx context.Context, startTime, endTime time.Time) (bool, error) {
	args := m.Called(ctx, startTime, endTime)
	return args.Bool(0), args.Error(1)
}

func (m *MockBookingRepository) Cancel(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBookingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTimeSlotRepository is a mock implementation of time slot repository
type MockTimeSlotRepository struct {
	mock.Mock
}

func (m *MockTimeSlotRepository) Create(ctx context.Context, slot *models.TimeSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockTimeSlotRepository) GetByID(ctx context.Context, id string) (*models.TimeSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TimeSlot), args.Error(1)
}

func (m *MockTimeSlotRepository) List(ctx context.Context, ownerID, eventTypeID string, page, pageSize int, available *bool, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	args := m.Called(ctx, ownerID, eventTypeID, page, pageSize, available, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.TimeSlot), args.Int(1), args.Error(2)
}

func (m *MockTimeSlotRepository) GetAvailableSlots(ctx context.Context, ownerID string, page, pageSize int, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	args := m.Called(ctx, ownerID, page, pageSize, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.TimeSlot), args.Int(1), args.Error(2)
}

func (m *MockTimeSlotRepository) MarkAsUnavailable(ctx context.Context, slotID string) error {
	args := m.Called(ctx, slotID)
	return args.Error(0)
}

// MockEventTypeRepository is a mock implementation of event type repository
type MockEventTypeRepository struct {
	mock.Mock
}

func (m *MockEventTypeRepository) Create(ctx context.Context, eventType *models.EventType) error {
	args := m.Called(ctx, eventType)
	return args.Error(0)
}

func (m *MockEventTypeRepository) GetByID(ctx context.Context, id string) (*models.EventType, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventType), args.Error(1)
}

func (m *MockEventTypeRepository) List(ctx context.Context, ownerID string, page, pageSize int, sortBy, sortOrder string) ([]models.EventType, int, error) {
	args := m.Called(ctx, ownerID, page, pageSize, sortBy, sortOrder)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.EventType), args.Int(1), args.Error(2)
}

func (m *MockEventTypeRepository) Patch(ctx context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventType), args.Error(1)
}

func (m *MockEventTypeRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockOwnerRepository is a mock implementation of owner repository
type MockOwnerRepository struct {
	mock.Mock
}

func (m *MockOwnerRepository) Create(ctx context.Context, owner *models.Owner) error {
	args := m.Called(ctx, owner)
	return args.Error(0)
}

func (m *MockOwnerRepository) GetByID(ctx context.Context, id string) (*models.Owner, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Owner), args.Error(1)
}

// MockSlotGenerationConfigRepository is a mock implementation of slot generation config repository
type MockSlotGenerationConfigRepository struct {
	mock.Mock
}

func (m *MockSlotGenerationConfigRepository) Create(ctx context.Context, config *models.SlotGenerationConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}
