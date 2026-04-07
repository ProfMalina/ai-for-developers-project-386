package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEventTypeService_Create_Success(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	ownerID := "test-owner-id"
	req := models.CreateEventTypeRequest{
		Name:            "Meeting",
		Description:     "30 minute meeting",
		DurationMinutes: 30,
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(et *models.EventType) bool {
		return et.OwnerID == ownerID && et.Name == "Meeting" && et.DurationMinutes == 30
	})).Return(nil)

	result, err := service.Create(context.Background(), ownerID, req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ownerID, result.OwnerID)
	assert.Equal(t, "Meeting", result.Name)
	assert.Equal(t, 30, result.DurationMinutes)
	mockRepo.AssertExpectations(t)
}

func TestEventTypeService_Create_Error(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	ownerID := "test-owner-id"
	req := models.CreateEventTypeRequest{
		Name:            "Meeting",
		Description:     "30 minute meeting",
		DurationMinutes: 30,
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))

	result, err := service.Create(context.Background(), ownerID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEventTypeService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	eventTypeID := "test-event-type-id"
	expected := &models.EventType{
		ID:              eventTypeID,
		Name:            "Meeting",
		DurationMinutes: 30,
	}

	mockRepo.On("GetByID", mock.Anything, eventTypeID).Return(expected, nil)

	result, err := service.GetByID(context.Background(), eventTypeID)

	require.NoError(t, err)
	assert.Equal(t, eventTypeID, result.ID)
	assert.Equal(t, "Meeting", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestEventTypeService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "non-existent").Return(nil, errors.New("not found"))

	result, err := service.GetByID(context.Background(), "non-existent")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestEventTypeService_List_Success(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	ownerID := "test-owner-id"
	eventTypes := []models.EventType{
		{ID: "1", Name: "Meeting"},
		{ID: "2", Name: "Workshop"},
	}

	mockRepo.On("List", mock.Anything, ownerID, 1, 20, "created_at", "desc").
		Return(eventTypes, 2, nil)

	result, err := service.List(context.Background(), ownerID, 1, 20, "created_at", "desc")

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.Pagination.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestEventTypeService_List_DefaultPagination(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	ownerID := "test-owner-id"
	// Service passes empty strings - repository handles defaults
	mockRepo.On("List", mock.Anything, ownerID, 1, 20, "", "").
		Return([]models.EventType{}, 0, nil)

	result, err := service.List(context.Background(), ownerID, 0, 0, "", "")

	require.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestEventTypeService_Update_Success(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	eventTypeID := "test-event-type-id"
	name := "Updated Meeting"
	req := models.UpdateEventTypeRequest{
		Name: &name,
	}

	updated := &models.EventType{
		ID:              eventTypeID,
		Name:            "Updated Meeting",
		DurationMinutes: 30,
	}

	mockRepo.On("Patch", mock.Anything, eventTypeID, req).Return(updated, nil)

	result, err := service.Update(context.Background(), eventTypeID, req)

	require.NoError(t, err)
	assert.Equal(t, "Updated Meeting", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestEventTypeService_Delete_Success(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	eventTypeID := "test-event-type-id"
	mockRepo.On("Delete", mock.Anything, eventTypeID).Return(nil)

	err := service.Delete(context.Background(), eventTypeID)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEventTypeService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockEventTypeRepository)
	service := NewEventTypeService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "non-existent").Return(errors.New("not found"))

	err := service.Delete(context.Background(), "non-existent")

	assert.Error(t, err)
}
