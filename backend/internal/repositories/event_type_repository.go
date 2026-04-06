package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// EventTypeRepository handles database operations for event types
type EventTypeRepository struct{}

// NewEventTypeRepository creates a new event type repository
func NewEventTypeRepository() *EventTypeRepository {
	return &EventTypeRepository{}
}

// Create creates a new event type
func (r *EventTypeRepository) Create(ctx context.Context, et *models.EventType) error {
	query := `
		INSERT INTO event_types (id, owner_id, name, description, duration_minutes, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	if et.ID == "" {
		et.ID = uuid.New().String()
	}

	et.IsActive = true

	err := db.Pool.QueryRow(ctx, query, et.ID, et.OwnerID, et.Name, et.Description, et.DurationMinutes, et.IsActive).
		Scan(&et.ID, &et.CreatedAt, &et.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create event type: %w", err)
	}

	return nil
}

// GetByID retrieves an event type by ID
func (r *EventTypeRepository) GetByID(ctx context.Context, id string) (*models.EventType, error) {
	query := `SELECT id, owner_id, name, description, duration_minutes, is_active, created_at, updated_at FROM event_types WHERE id = $1`

	et := &models.EventType{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&et.ID, &et.OwnerID, &et.Name, &et.Description, &et.DurationMinutes, &et.IsActive, &et.CreatedAt, &et.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("event type not found")
		}
		return nil, fmt.Errorf("failed to get event type: %w", err)
	}

	return et, nil
}

// List retrieves a paginated list of event types
func (r *EventTypeRepository) List(ctx context.Context, ownerID string, page, pageSize int, sortBy, sortOrder string) ([]models.EventType, int, error) {
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	allowedSortFields := map[string]bool{
		"name": true, "created_at": true, "updated_at": true, "duration_minutes": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM event_types WHERE owner_id = $1`
	var totalItems int
	err := db.Pool.QueryRow(ctx, countQuery, ownerID).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event types: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, owner_id, name, description, duration_minutes, is_active, created_at, updated_at
		FROM event_types
		WHERE owner_id = $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, sortOrder)

	rows, err := db.Pool.Query(ctx, query, ownerID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list event types: %w", err)
	}
	defer rows.Close()

	var eventTypes []models.EventType
	for rows.Next() {
		var et models.EventType
		err := rows.Scan(&et.ID, &et.OwnerID, &et.Name, &et.Description, &et.DurationMinutes, &et.IsActive, &et.CreatedAt, &et.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event type: %w", err)
		}
		eventTypes = append(eventTypes, et)
	}

	if eventTypes == nil {
		eventTypes = []models.EventType{}
	}

	return eventTypes, totalItems, nil
}

// Update updates an event type
func (r *EventTypeRepository) Update(ctx context.Context, et *models.EventType) error {
	query := `
		UPDATE event_types
		SET name = $2, description = $3, duration_minutes = $4, is_active = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := db.Pool.QueryRow(ctx, query, et.ID, et.Name, et.Description, et.DurationMinutes, et.IsActive).
		Scan(&et.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("event type not found")
		}
		return fmt.Errorf("failed to update event type: %w", err)
	}

	return nil
}

// Patch partially updates an event type
func (r *EventTypeRepository) Patch(ctx context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error) {
	et, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}
	if req.DurationMinutes != nil {
		setClauses = append(setClauses, fmt.Sprintf("duration_minutes = $%d", argIdx))
		args = append(args, *req.DurationMinutes)
		argIdx++
	}
	if req.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *req.IsActive)
		argIdx++
	}

	if len(setClauses) == 0 {
		return et, nil
	}

	// Add id to args BEFORE building query
	args = append(args, id)

	query := fmt.Sprintf("UPDATE event_types SET %s, updated_at = NOW() WHERE id = $%d::uuid RETURNING id, owner_id, name, description, duration_minutes, is_active, created_at, updated_at",
		strings.Join(setClauses, ", "), argIdx+1)

	err = db.Pool.QueryRow(ctx, query, args...).Scan(
		&et.ID, &et.OwnerID, &et.Name, &et.Description, &et.DurationMinutes, &et.IsActive, &et.CreatedAt, &et.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("event type not found")
		}
		return nil, fmt.Errorf("failed to patch event type: %w", err)
	}

	return et, nil
}

// Delete deletes an event type
func (r *EventTypeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM event_types WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("event type not found")
	}

	return nil
}
