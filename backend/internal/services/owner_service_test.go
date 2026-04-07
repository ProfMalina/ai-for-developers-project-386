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

func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		totalItems int
		expected   struct {
			totalPages int
			hasNext    bool
			hasPrev    bool
		}
	}{
		{
			name:       "first page with results",
			page:       1,
			pageSize:   10,
			totalItems: 25,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 3, hasNext: true, hasPrev: false},
		},
		{
			name:       "middle page",
			page:       2,
			pageSize:   10,
			totalItems: 25,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 3, hasNext: true, hasPrev: true},
		},
		{
			name:       "last page",
			page:       3,
			pageSize:   10,
			totalItems: 25,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 3, hasNext: false, hasPrev: true},
		},
		{
			name:       "empty results",
			page:       1,
			pageSize:   10,
			totalItems: 0,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 0, hasNext: false, hasPrev: false},
		},
		{
			name:       "single item",
			page:       1,
			pageSize:   10,
			totalItems: 1,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 1, hasNext: false, hasPrev: false},
		},
		{
			name:       "exact page size",
			page:       1,
			pageSize:   10,
			totalItems: 10,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 1, hasNext: false, hasPrev: false},
		},
		{
			name:       "large page size",
			page:       1,
			pageSize:   100,
			totalItems: 50,
			expected: struct {
				totalPages int
				hasNext    bool
				hasPrev    bool
			}{totalPages: 1, hasNext: false, hasPrev: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePagination(tt.page, tt.pageSize, tt.totalItems)

			assert.Equal(t, tt.page, result.Page)
			assert.Equal(t, tt.pageSize, result.PageSize)
			assert.Equal(t, tt.totalItems, result.TotalItems)
			assert.Equal(t, tt.expected.totalPages, result.TotalPages)
			assert.Equal(t, tt.expected.hasNext, result.HasNext)
			assert.Equal(t, tt.expected.hasPrev, result.HasPrev)
		})
	}
}

func TestOwnerService_Create_Success(t *testing.T) {
	mockOwnerRepo := new(MockOwnerRepository)
	service := NewOwnerService(mockOwnerRepo)

	req := models.CreateOwnerRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Timezone: "UTC",
	}

	mockOwnerRepo.On("Create", mock.Anything, mock.MatchedBy(func(owner *models.Owner) bool {
		return owner.Name == "John Doe" && owner.Email == "john@example.com" && owner.Timezone == "UTC"
	})).Return(nil)

	result, err := service.Create(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, "UTC", result.Timezone)
	mockOwnerRepo.AssertExpectations(t)
}

func TestOwnerService_Create_Error(t *testing.T) {
	mockOwnerRepo := new(MockOwnerRepository)
	service := NewOwnerService(mockOwnerRepo)

	req := models.CreateOwnerRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Timezone: "UTC",
	}

	mockOwnerRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("database error"))

	result, err := service.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database error")
}

func TestOwnerService_GetByID_Success(t *testing.T) {
	mockOwnerRepo := new(MockOwnerRepository)
	service := NewOwnerService(mockOwnerRepo)

	ownerID := "test-owner-id"
	expectedOwner := &models.Owner{
		ID:       ownerID,
		Name:     "John Doe",
		Email:    "john@example.com",
		Timezone: "UTC",
	}

	mockOwnerRepo.On("GetByID", mock.Anything, ownerID).Return(expectedOwner, nil)

	result, err := service.GetByID(context.Background(), ownerID)

	require.NoError(t, err)
	assert.Equal(t, ownerID, result.ID)
	assert.Equal(t, "John Doe", result.Name)
	mockOwnerRepo.AssertExpectations(t)
}

func TestOwnerService_GetByID_NotFound(t *testing.T) {
	mockOwnerRepo := new(MockOwnerRepository)
	service := NewOwnerService(mockOwnerRepo)

	ownerID := "non-existent-id"
	mockOwnerRepo.On("GetByID", mock.Anything, ownerID).Return(nil, errors.New("owner not found"))

	result, err := service.GetByID(context.Background(), ownerID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "owner not found")
}
