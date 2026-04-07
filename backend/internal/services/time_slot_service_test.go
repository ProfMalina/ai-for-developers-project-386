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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

	mockRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))

	result, err := service.GetByID(context.Background(), "non-existent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTimeSlotService_List_Success(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

	ownerID := "test-owner-id"
	slots := []models.TimeSlot{
		{ID: "1", OwnerID: ownerID, IsAvailable: true},
		{ID: "2", OwnerID: ownerID, IsAvailable: false},
	}

	mockRepo.On("List", mock.Anything, ownerID, 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).
		Return(slots, 2, nil)

	result, err := service.List(context.Background(), ownerID, 1, 20, nil, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.Pagination.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_List_DefaultPagination(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

	ownerID := "test-owner-id"
	mockRepo.On("List", mock.Anything, ownerID, 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).
		Return([]models.TimeSlot{}, 0, nil)

	result, err := service.List(context.Background(), ownerID, 0, 0, nil, nil, nil)

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertCalled(t, "List", mock.Anything, ownerID, 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil))
}

func TestTimeSlotService_List_WithFilters(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

	ownerID := "test-owner-id"
	available := true
	startTime := time.Now()
	endTime := time.Now().Add(24 * time.Hour)

	mockRepo.On("List", mock.Anything, ownerID, 1, 20, &available, &startTime, &endTime).
		Return([]models.TimeSlot{}, 0, nil)

	result, err := service.List(context.Background(), ownerID, 1, 20, &available, &startTime, &endTime)

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestTimeSlotService_GetAvailableSlots_Success(t *testing.T) {
	mockRepo := new(MockTimeSlotRepository)
	mockConfigRepo := new(MockSlotGenerationConfigRepository)
	mockOwnerRepo := new(MockOwnerRepository)

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

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

	service := NewTimeSlotService(mockRepo, mockConfigRepo, mockOwnerRepo)

	ownerID := "test-owner-id"
	mockRepo.On("List", mock.Anything, ownerID, 1, 20, (*bool)(nil), (*time.Time)(nil), (*time.Time)(nil)).
		Return(nil, 0, errors.New("database error"))

	result, err := service.List(context.Background(), ownerID, 1, 20, nil, nil, nil)

	assert.Error(t, err)
	assert.Nil(t, result)
}
