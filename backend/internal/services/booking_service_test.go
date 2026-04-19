package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBookingService_Create_Success(t *testing.T) {
	// Setup mocks
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	// Test data
	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	futureTime := time.Now().Add(24 * time.Hour)
	slot := &models.TimeSlot{
		ID:          slotID,
		StartTime:   futureTime,
		EndTime:     futureTime.Add(30 * time.Minute),
		IsAvailable: true,
	}

	eventType := &models.EventType{
		ID:              eventTypeID,
		Name:            "Meeting",
		DurationMinutes: 30,
	}

	// Setup expectations
	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(slot, nil)
	mockBookingRepo.On("CheckOverlap", mock.Anything, slot.StartTime, slot.EndTime).Return(false, nil)
	mockBookingRepo.On("CreateWithReservedSlot", mock.Anything, mock.MatchedBy(func(b *models.Booking) bool {
		return b.EventTypeID == eventTypeID && b.GuestName == "John Doe"
	})).Return(nil)

	// Execute
	result, err := service.Create(context.Background(), req)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, eventTypeID, result.EventTypeID)
	assert.Equal(t, "John Doe", result.GuestName)
	assert.Equal(t, "confirmed", result.Status)
	mockSlotRepo.AssertNotCalled(t, "MarkAsUnavailable", mock.Anything, slotID)

	mockBookingRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockEtRepo.AssertExpectations(t)
}

func TestBookingService_Create_EventTypeNotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	req := models.CreateBookingRequest{
		EventTypeID: "non-existent-id",
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	mockEtRepo.On("GetByID", mock.Anything, "non-existent-id").Return(nil, errors.New("event type not found"))

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "event type not found")

	mockEtRepo.AssertExpectations(t)
}

func TestBookingService_Create_SlotNotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	eventType := &models.EventType{ID: eventTypeID, Name: "Meeting"}

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(nil, errors.New("slot not found"))

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "time slot not found")
}

func TestBookingService_Create_SlotNotAvailable(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	eventType := &models.EventType{ID: eventTypeID}
	slot := &models.TimeSlot{
		ID:          slotID,
		IsAvailable: false,
	}

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(slot, nil)

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "time slot is already booked")
}

func TestBookingService_Create_SlotAlreadyStarted(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	eventType := &models.EventType{ID: eventTypeID}
	pastTime := time.Now().Add(-1 * time.Hour)
	slot := &models.TimeSlot{
		ID:          slotID,
		StartTime:   pastTime,
		IsAvailable: true,
	}

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(slot, nil)

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "cannot book a time slot that has already started")
}

func TestBookingService_Create_OverlapDetected(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	eventType := &models.EventType{ID: eventTypeID}
	futureTime := time.Now().Add(24 * time.Hour)
	slot := &models.TimeSlot{
		ID:          slotID,
		StartTime:   futureTime,
		EndTime:     futureTime.Add(30 * time.Minute),
		IsAvailable: true,
	}

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(slot, nil)
	mockBookingRepo.On("CheckOverlap", mock.Anything, slot.StartTime, slot.EndTime).Return(true, nil)

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "selected time slot is already booked")
}

func TestBookingService_GetByID_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	bookingID := "test-booking-id"
	expectedBooking := &models.Booking{
		ID:        bookingID,
		GuestName: "John Doe",
	}

	mockBookingRepo.On("GetByID", mock.Anything, bookingID).Return(expectedBooking, nil)

	result, err := service.GetByID(context.Background(), bookingID)

	require.NoError(t, err)
	assert.Equal(t, bookingID, result.ID)
	assert.Equal(t, "John Doe", result.GuestName)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_GetByID_NotFound(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	bookingID := "non-existent-id"
	mockBookingRepo.On("GetByID", mock.Anything, bookingID).Return(nil, errors.New("booking not found"))

	result, err := service.GetByID(context.Background(), bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBookingService_List_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	bookings := []models.Booking{
		{ID: "1", GuestName: "John"},
		{ID: "2", GuestName: "Jane"},
	}

	mockBookingRepo.On("List", mock.Anything, 1, 20, "created_at", "desc", (*time.Time)(nil), (*time.Time)(nil)).
		Return(bookings, 2, nil)

	result, err := service.List(context.Background(), 1, 20, "created_at", "desc", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.Pagination.TotalItems)
	assert.Equal(t, 1, result.Pagination.TotalPages)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_List_DefaultPagination(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	// Service passes empty strings for sortBy/sortOrder when not provided - repository sets defaults
	mockBookingRepo.On("List", mock.Anything, 1, 20, "", "", (*time.Time)(nil), (*time.Time)(nil)).
		Return([]models.Booking{}, 0, nil)

	// Test with invalid page and pageSize
	result, err := service.List(context.Background(), 0, 0, "", "", nil, nil)

	require.NoError(t, err)
	assert.NotNil(t, result)
	// Should use defaults: page=1, pageSize=20
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_Cancel_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	bookingID := "test-booking-id"
	booking := &models.Booking{
		ID:     bookingID,
		Status: "confirmed",
	}

	mockBookingRepo.On("GetByID", mock.Anything, bookingID).Return(booking, nil)
	mockBookingRepo.On("Cancel", mock.Anything, bookingID).Return(nil)

	err := service.Cancel(context.Background(), bookingID)

	require.NoError(t, err)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_Cancel_AlreadyCancelled(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	bookingID := "test-booking-id"
	booking := &models.Booking{
		ID:     bookingID,
		Status: "cancelled", //nolint:misspell // persisted booking status value
	}

	mockBookingRepo.On("GetByID", mock.Anything, bookingID).Return(booking, nil)

	err := service.Cancel(context.Background(), bookingID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "booking is already canceled")
}

func TestBookingService_Delete_Success(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	bookingID := "test-booking-id"
	mockBookingRepo.On("Delete", mock.Anything, bookingID).Return(nil)

	err := service.Delete(context.Background(), bookingID)

	require.NoError(t, err)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_Create_NoSlotID(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	eventType := &models.EventType{ID: eventTypeID}
	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "slot ID is required")
}

func TestBookingService_List_WithDateFilters(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	dateFrom := time.Now().Add(-time.Hour)
	dateTo := time.Now().Add(time.Hour)
	mockBookingRepo.On("List", mock.Anything, 1, 20, "created_at", "desc", &dateFrom, &dateTo).
		Return([]models.Booking{}, 0, nil)

	result, err := service.List(context.Background(), 1, 20, "created_at", "desc", &dateFrom, &dateTo)

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockBookingRepo.AssertExpectations(t)
}

func TestBookingService_Create_DBError(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	futureTime := time.Now().Add(24 * time.Hour)
	slot := &models.TimeSlot{
		ID:          slotID,
		StartTime:   futureTime,
		EndTime:     futureTime.Add(30 * time.Minute),
		IsAvailable: true,
	}

	eventType := &models.EventType{ID: eventTypeID}

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(slot, nil)
	mockBookingRepo.On("CheckOverlap", mock.Anything, slot.StartTime, slot.EndTime).Return(false, nil)
	mockBookingRepo.On("CreateWithReservedSlot", mock.Anything, mock.Anything).Return(fmt.Errorf("database error"))

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
	mockSlotRepo.AssertNotCalled(t, "MarkAsUnavailable", mock.Anything, slotID)
}

func TestBookingService_Create_AtomicReservationFailure(t *testing.T) {
	mockBookingRepo := new(MockBookingRepository)
	mockSlotRepo := new(MockTimeSlotRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewBookingService(mockBookingRepo, mockSlotRepo, mockEtRepo)

	slotID := "test-slot-id"
	eventTypeID := "test-event-type-id"
	req := models.CreateBookingRequest{
		EventTypeID: eventTypeID,
		SlotID:      &slotID,
		GuestName:   "John Doe",
		GuestEmail:  "john@example.com",
	}

	futureTime := time.Now().Add(24 * time.Hour)
	slot := &models.TimeSlot{
		ID:          slotID,
		StartTime:   futureTime,
		EndTime:     futureTime.Add(30 * time.Minute),
		IsAvailable: true,
	}

	eventType := &models.EventType{ID: eventTypeID}

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockSlotRepo.On("GetByID", mock.Anything, slotID).Return(slot, nil)
	mockBookingRepo.On("CheckOverlap", mock.Anything, slot.StartTime, slot.EndTime).Return(false, nil)
	mockBookingRepo.On("CreateWithReservedSlot", mock.Anything, mock.Anything).Return(fmt.Errorf("slot reservation failed"))

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "slot reservation failed")
	mockSlotRepo.AssertNotCalled(t, "MarkAsUnavailable", mock.Anything, slotID)
	mockBookingRepo.AssertExpectations(t)
	mockSlotRepo.AssertExpectations(t)
	mockEtRepo.AssertExpectations(t)
}
