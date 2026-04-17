package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TimeSlotRepository handles database operations for time slots
type TimeSlotRepository struct{}

// NewTimeSlotRepository creates a new time slot repository
func NewTimeSlotRepository() *TimeSlotRepository {
	return &TimeSlotRepository{}
}

// Create creates a new time slot
func (r *TimeSlotRepository) Create(ctx context.Context, slot *models.TimeSlot) error {
	query := `
		INSERT INTO time_slots (id, owner_id, event_type_id, start_time, end_time, is_available, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	if slot.ID == "" {
		slot.ID = uuid.New().String()
	}

	err := db.Pool.QueryRow(ctx, query, slot.ID, slot.OwnerID, slot.EventTypeID, slot.StartTime, slot.EndTime, slot.IsAvailable).
		Scan(&slot.ID, &slot.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create time slot: %w", err)
	}

	return nil
}

// GetByID retrieves a time slot by ID
func (r *TimeSlotRepository) GetByID(ctx context.Context, id string) (*models.TimeSlot, error) {
	query := `SELECT id, owner_id, COALESCE(event_type_id::text, ''), start_time, end_time, is_available, created_at FROM time_slots WHERE id = $1`

	slot := &models.TimeSlot{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&slot.ID, &slot.OwnerID, &slot.EventTypeID, &slot.StartTime, &slot.EndTime, &slot.IsAvailable, &slot.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("time slot not found")
		}
		return nil, fmt.Errorf("failed to get time slot: %w", err)
	}

	return slot, nil
}

// List retrieves a paginated list of time slots with filters
func (r *TimeSlotRepository) List(ctx context.Context, ownerID, eventTypeID string, page, pageSize int, available *bool, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	offset := (page - 1) * pageSize

	// Build query dynamically based on filters
	query := `SELECT id, owner_id, COALESCE(event_type_id::text, ''), start_time, end_time, is_available, created_at FROM time_slots WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM time_slots WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if ownerID != "" {
		query += fmt.Sprintf(" AND owner_id = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND owner_id = $%d", argIdx)
		args = append(args, ownerID)
		argIdx++
	}

	if eventTypeID != "" {
		query += fmt.Sprintf(" AND event_type_id = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND event_type_id = $%d", argIdx)
		args = append(args, eventTypeID)
		argIdx++
	}

	if available != nil {
		query += fmt.Sprintf(" AND is_available = $%d", argIdx)
		countQuery += fmt.Sprintf(" AND is_available = $%d", argIdx)
		args = append(args, *available)
		argIdx++
	}

	if startTime != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND start_time >= $%d", argIdx)
		args = append(args, *startTime)
		argIdx++
	}

	if endTime != nil {
		query += fmt.Sprintf(" AND end_time <= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND end_time <= $%d", argIdx)
		args = append(args, *endTime)
		argIdx++
	}

	// Count total items
	var totalItems int
	err := db.Pool.QueryRow(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count time slots: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY start_time ASC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list time slots: %w", err)
	}
	defer rows.Close()

	var slots []models.TimeSlot
	for rows.Next() {
		var slot models.TimeSlot
		err := rows.Scan(&slot.ID, &slot.OwnerID, &slot.EventTypeID, &slot.StartTime, &slot.EndTime, &slot.IsAvailable, &slot.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan time slot: %w", err)
		}
		slots = append(slots, slot)
	}

	if slots == nil {
		slots = []models.TimeSlot{}
	}

	return slots, totalItems, nil
}

// GetAvailableSlots retrieves available slots for an event type
func (r *TimeSlotRepository) GetAvailableSlots(ctx context.Context, eventTypeID string, page, pageSize int, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	available := true
	return r.List(ctx, eventTypeID, "", page, pageSize, &available, startTime, endTime)
}

// MarkAsUnavailable marks a slot as unavailable (booked)
func (r *TimeSlotRepository) MarkAsUnavailable(ctx context.Context, id string) error {
	query := `UPDATE time_slots SET is_available = false WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark slot as unavailable: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("time slot not found")
	}

	return nil
}

// DeleteAvailableInRange deletes available slots for an owner and event type in a time range
func (r *TimeSlotRepository) DeleteAvailableInRange(ctx context.Context, ownerID, eventTypeID string, startTime, endTime time.Time) error {
	query, args := buildDeleteAvailableInRangeQuery(ownerID, eventTypeID, startTime, endTime)

	_, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete available time slots in range: %w", err)
	}

	return nil
}

func buildDeleteAvailableInRangeQuery(ownerID, eventTypeID string, startTime, endTime time.Time) (string, []interface{}) {
	query := `
		DELETE FROM time_slots
		WHERE owner_id = $1
		AND (event_type_id = $2 OR event_type_id IS NULL)
		AND is_available = true
		AND start_time < $4
		AND end_time > $3
		AND NOT EXISTS (
			SELECT 1 FROM bookings b
			WHERE b.slot_id = time_slots.id
			AND b.status != 'cancelled'
		)
	`

	args := []interface{}{ownerID, eventTypeID, startTime, endTime}
	return query, args
}

// Delete deletes a time slot
func (r *TimeSlotRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM time_slots WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete time slot: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("time slot not found")
	}

	return nil
}

// DeleteByOwnerID deletes all time slots for an event type
func (r *TimeSlotRepository) DeleteByOwnerID(ctx context.Context, eventTypeID string) error {
	query := `DELETE FROM time_slots WHERE owner_id = $1`

	_, err := db.Pool.Exec(ctx, query, eventTypeID)
	if err != nil {
		return fmt.Errorf("failed to delete time slots for event type: %w", err)
	}

	return nil
}
