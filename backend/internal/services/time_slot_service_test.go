package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTimeSlotService_Create_Success(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	slot := &models.TimeSlot{
		ID:          "test-slot-id",
		OwnerID:     "test-owner-id",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(30 * time.Minute),
		IsAvailable: true,
	}

	mockRepo.On("Create", mock.Anything, slot).Return(nil)

	err := service.Create(context.Background(), slot)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_Create_Error(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	slot := &models.TimeSlot{
		ID: "test-slot-id",
	}

	mockRepo.On("Create", mock.Anything, slot).Return(errors.New("database error"))

	err := service.Create(context.Background(), slot)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestTimeSlotService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	slotID := "test-slot-id"
	expected := &models.TimeSlot{
		ID:          slotID,
		OwnerID:     "test-owner-id",
		IsAvailable: true,
	}

	mockRepo.On("GetByID", mock.Anything, slotID).Return(expected, nil)

	result, err := service.GetByID(context.Background(), slotID)

	require.NoError(t, err)
	assert.Equal(t, slotID, result.ID)
	assert.Equal(t, "test-owner-id", result.OwnerID)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	mockRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))

	result, err := service.GetByID(context.Background(), "non-existent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTimeSlotService_List_Success(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	slots := []models.TimeSlot{
		{ID: "1", OwnerID: ownerID, IsAvailable: true},
		{ID: "2", OwnerID: ownerID, IsAvailable: false},
	}

	mockRepo.On("List", mock.Anything, ownerID, "", 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).
		Return(slots, 2, nil)

	result, err := service.List(context.Background(), ownerID, "", 1, 20, nil, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.Pagination.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_List_DefaultPagination(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	mockRepo.On("List", mock.Anything, ownerID, "", 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).
		Return([]models.TimeSlot{}, 0, nil)

	result, err := service.List(context.Background(), ownerID, "", 0, 0, nil, nil, nil)

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertCalled(t, "List", mock.Anything, ownerID, "", 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil))
}

func TestTimeSlotService_List_WithFilters(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	available := true
	startTime := time.Now()
	endTime := time.Now().Add(24 * time.Hour)

	mockRepo.On("List", mock.Anything, ownerID, "", 1, 20, &available, &startTime, &endTime).
		Return([]models.TimeSlot{}, 0, nil)

	result, err := service.List(context.Background(), ownerID, "", 1, 20, &available, &startTime, &endTime)

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_GetAvailableSlots_Success(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	futureTime := time.Now().Add(24 * time.Hour)
	slots := []models.TimeSlot{
		{ID: "1", OwnerID: ownerID, StartTime: futureTime, IsAvailable: true},
		{ID: "2", OwnerID: ownerID, StartTime: futureTime.Add(30 * time.Minute), IsAvailable: true},
	}

	mockRepo.On("GetAvailableSlots", mock.Anything, ownerID, 1, 20, (*time.Time)(nil), (*time.Time)(nil)).
		Return(slots, 2, nil)

	result, err := service.GetAvailableSlots(context.Background(), ownerID, 1, 20, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_GetAvailableSlots_FiltersPastSlots(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	pastTime := time.Now().Add(-1 * time.Hour)
	futureTime := time.Now().Add(1 * time.Hour)

	// One slot in the past, one in the future
	slots := []models.TimeSlot{
		{ID: "1", OwnerID: ownerID, StartTime: pastTime, IsAvailable: true},
		{ID: "2", OwnerID: ownerID, StartTime: futureTime, IsAvailable: true},
	}

	mockRepo.On("GetAvailableSlots", mock.Anything, ownerID, 1, 20, (*time.Time)(nil), (*time.Time)(nil)).
		Return(slots, 2, nil)

	result, err := service.GetAvailableSlots(context.Background(), ownerID, 1, 20, nil, nil)

	require.NoError(t, err)
	// Should filter out the past slot, so only 1 slot remains
	assert.Len(t, result.Items, 1)
	assert.Equal(t, 1, result.Pagination.TotalItems)
}

func TestTimeSlotService_GetAvailableSlots_EmptyList(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	mockRepo.On("GetAvailableSlots", mock.Anything, ownerID, 1, 20, (*time.Time)(nil), (*time.Time)(nil)).
		Return([]models.TimeSlot{}, 0, nil)

	result, err := service.GetAvailableSlots(context.Background(), ownerID, 1, 20, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result.Items, 0)
	assert.Equal(t, 0, result.Pagination.TotalItems)
}

func TestTimeSlotService_List_Error(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	mockRepo.On("List", mock.Anything, ownerID, "", 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).
		Return(nil, 0, errors.New("database error"))

	result, err := service.List(context.Background(), ownerID, "", 1, 20, nil, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTimeSlotService_GenerateSlots_ReplacesExistingSlotsInWindow(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	eventTypeID := "event-type-id"
	req := models.SlotGenerationRequest{
		WorkingHoursStart: "09:00",
		WorkingHoursEnd:   "11:00",
		IntervalMinutes:   30,
		DaysOfWeek:        []int{1},
		DateFrom:          "2026-04-20",
		DateTo:            "2026-04-20",
	}

	owner := &models.Owner{ID: ownerID, Timezone: "Europe/Moscow"}
	eventType := &models.EventType{ID: eventTypeID, DurationMinutes: 60}
	windowStart := time.Date(2026, 4, 20, 0, 0, 0, 0, time.FixedZone("MSK", 3*60*60)).UTC()
	windowEnd := time.Date(2026, 4, 21, 0, 0, 0, 0, time.FixedZone("MSK", 3*60*60)).UTC()

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockOwnerRepo.On("GetByID", mock.Anything, ownerID).Return(owner, nil)
	mockConfigRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.SlotGenerationConfig")).Return(nil)
	mockRepo.On("DeleteAvailableInRange", mock.Anything, ownerID, eventTypeID, windowStart, windowEnd).Return(nil)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(slot *models.TimeSlot) bool {
		duration := slot.EndTime.Sub(slot.StartTime)
		return slot.OwnerID == ownerID && slot.EventTypeID == eventTypeID && duration == time.Hour
	})).Return(nil).Times(3)

	result, err := service.GenerateSlots(context.Background(), ownerID, eventTypeID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.SlotsCreated)
	assert.Equal(t, 0, result.SlotsSkipped)
	assert.Equal(t, "2026-04-20", result.DateFrom)
	assert.Equal(t, "2026-04-20", result.DateTo)
	mockRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockOwnerRepo.AssertExpectations(t)
	mockEtRepo.AssertExpectations(t)
}

func TestTimeSlotService_GenerateSlots_AppliesDefaultsAndSundayNumbering(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	eventTypeID := "event-type-id"
	nowUTC := time.Now().UTC()
	targetDate := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	for targetDate.Weekday() != time.Sunday {
		targetDate = targetDate.AddDate(0, 0, 1)
	}
	targetDateStr := targetDate.Format("2006-01-02")
	req := models.SlotGenerationRequest{
		WorkingHoursStart: "09:00",
		WorkingHoursEnd:   "10:00",
		DateFrom:          targetDateStr,
		DateTo:            targetDateStr,
		DaysOfWeek:        []int{0},
	}

	owner := &models.Owner{ID: ownerID, Timezone: "UTC"}
	eventType := &models.EventType{ID: eventTypeID, DurationMinutes: 30}
	windowStart := targetDate
	windowEnd := targetDate.AddDate(0, 0, 1)

	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(eventType, nil)
	mockOwnerRepo.On("GetByID", mock.Anything, ownerID).Return(owner, nil)
	mockConfigRepo.On("Create", mock.Anything, mock.MatchedBy(func(config *models.SlotGenerationConfig) bool {
		return config.IntervalMinutes == 30 && len(config.DaysOfWeek) == 1 && config.DaysOfWeek[0] == 0 && config.Timezone == "UTC"
	})).Return(nil)
	mockRepo.On("DeleteAvailableInRange", mock.Anything, ownerID, eventTypeID, windowStart, windowEnd).Return(nil)
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(slot *models.TimeSlot) bool {
		return slot.EventTypeID == eventTypeID
	})).Return(nil).Twice()

	result, err := service.GenerateSlots(context.Background(), ownerID, eventTypeID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.SlotsCreated)
	assert.Equal(t, 0, result.SlotsSkipped)
	mockRepo.AssertExpectations(t)
	mockConfigRepo.AssertExpectations(t)
	mockOwnerRepo.AssertExpectations(t)
	mockEtRepo.AssertExpectations(t)
}

func TestTimeSlotService_List_WithEventTypeFilter(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)
	mockEtRepo := new(MockEventTypeRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo, mockEtRepo)

	ownerID := "test-owner-id"
	eventTypeID := "event-type-id"
	mockEtRepo.On("GetByID", mock.Anything, eventTypeID).Return(&models.EventType{ID: eventTypeID}, nil)
	mockRepo.On("List", mock.Anything, ownerID, eventTypeID, 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).Return([]models.TimeSlot{}, 0, nil)

	result, err := service.List(context.Background(), ownerID, eventTypeID, 1, 20, nil, nil, nil)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Items, 0)
	mockRepo.AssertExpectations(t)
	mockEtRepo.AssertExpectations(t)
}
