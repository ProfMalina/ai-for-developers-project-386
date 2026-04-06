package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// OwnerRepository handles database operations for owners
type OwnerRepository struct{}

// NewOwnerRepository creates a new owner repository
func NewOwnerRepository() *OwnerRepository {
	return &OwnerRepository{}
}

// Create creates a new owner
func (r *OwnerRepository) Create(ctx context.Context, owner *models.Owner) error {
	query := `
		INSERT INTO owners (id, name, email, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	if owner.ID == "" {
		owner.ID = uuid.New().String()
	}

	err := db.Pool.QueryRow(ctx, query, owner.ID, owner.Name, owner.Email, owner.Timezone).
		Scan(&owner.ID, &owner.CreatedAt, &owner.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create owner: %w", err)
	}

	return nil
}

// GetByID retrieves an owner by ID
func (r *OwnerRepository) GetByID(ctx context.Context, id string) (*models.Owner, error) {
	query := `SELECT id, name, email, timezone, created_at, updated_at FROM owners WHERE id = $1`

	owner := &models.Owner{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&owner.ID, &owner.Name, &owner.Email, &owner.Timezone, &owner.CreatedAt, &owner.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("owner not found")
		}
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}

	return owner, nil
}

// Update updates an owner
func (r *OwnerRepository) Update(ctx context.Context, owner *models.Owner) error {
	query := `
		UPDATE owners
		SET name = $2, email = $3, timezone = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := db.Pool.QueryRow(ctx, query, owner.ID, owner.Name, owner.Email, owner.Timezone).
		Scan(&owner.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("owner not found")
		}
		return fmt.Errorf("failed to update owner: %w", err)
	}

	return nil
}

// Delete deletes an owner
func (r *OwnerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM owners WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete owner: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("owner not found")
	}

	return nil
}
