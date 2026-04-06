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

// SlotGenerationConfigRepository handles database operations for slot generation configs
type SlotGenerationConfigRepository struct{}

// NewSlotGenerationConfigRepository creates a new slot generation config repository
func NewSlotGenerationConfigRepository() *SlotGenerationConfigRepository {
	return &SlotGenerationConfigRepository{}
}

// Create creates or updates a slot generation config
func (r *SlotGenerationConfigRepository) Create(ctx context.Context, config *models.SlotGenerationConfig) error {
	query := `
		INSERT INTO slot_generation_configs (id, owner_id, working_hours_start, working_hours_end, interval_minutes, days_of_week, date_from, date_to, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (owner_id) DO UPDATE SET
			working_hours_start = EXCLUDED.working_hours_start,
			working_hours_end = EXCLUDED.working_hours_end,
			interval_minutes = EXCLUDED.interval_minutes,
			days_of_week = EXCLUDED.days_of_week,
			date_from = EXCLUDED.date_from,
			date_to = EXCLUDED.date_to,
			timezone = EXCLUDED.timezone,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	if config.ID == "" {
		config.ID = uuid.New().String()
	}

	err := db.Pool.QueryRow(ctx, query, config.ID, config.OwnerID, config.WorkingHoursStart,
		config.WorkingHoursEnd, config.IntervalMinutes, config.DaysOfWeek,
		config.DateFrom, config.DateTo, config.Timezone).
		Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create slot generation config: %w", err)
	}

	return nil
}

// GetByOwnerID retrieves a slot generation config by owner ID
func (r *SlotGenerationConfigRepository) GetByOwnerID(ctx context.Context, ownerID string) (*models.SlotGenerationConfig, error) {
	query := `
		SELECT id, owner_id, working_hours_start, working_hours_end, interval_minutes, days_of_week, date_from, date_to, timezone, created_at, updated_at
		FROM slot_generation_configs WHERE owner_id = $1
	`

	config := &models.SlotGenerationConfig{}
	err := db.Pool.QueryRow(ctx, query, ownerID).Scan(
		&config.ID, &config.OwnerID, &config.WorkingHoursStart, &config.WorkingHoursEnd,
		&config.IntervalMinutes, &config.DaysOfWeek, &config.DateFrom, &config.DateTo,
		&config.Timezone, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("slot generation config not found")
		}
		return nil, fmt.Errorf("failed to get slot generation config: %w", err)
	}

	return config, nil
}
