package services

import (
	"context"
	"math"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/repositories"
)

// OwnerService handles business logic for owners
type OwnerService struct {
	repo *repositories.OwnerRepository
}

// NewOwnerService creates a new owner service
func NewOwnerService() *OwnerService {
	return &OwnerService{
		repo: repositories.NewOwnerRepository(),
	}
}

// Create creates a new owner
func (s *OwnerService) Create(ctx context.Context, req models.CreateOwnerRequest) (*models.Owner, error) {
	owner := &models.Owner{
		Name:     req.Name,
		Email:    req.Email,
		Timezone: req.Timezone,
	}

	if err := s.repo.Create(ctx, owner); err != nil {
		return nil, err
	}

	return owner, nil
}

// GetByID retrieves an owner by ID
func (s *OwnerService) GetByID(ctx context.Context, id string) (*models.Owner, error) {
	return s.repo.GetByID(ctx, id)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize, totalItems int) models.Pagination {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return models.Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}
