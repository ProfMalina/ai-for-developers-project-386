package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ScheduleRepository struct{}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{}
}

func (r *ScheduleRepository) GetDaySchedules(ctx context.Context, ownerID string) ([]models.DaySchedule, error) {
	query := `
		SELECT id, owner_id, day_of_week, windows, breaks, created_at, updated_at
		FROM day_schedules WHERE owner_id = $1 ORDER BY day_of_week
	`
	rows, err := db.Pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get day schedules: %w", err)
	}
	defer rows.Close()

	var schedules []models.DaySchedule
	for rows.Next() {
		var s models.DaySchedule
		var windowsJSON, breaksJSON []byte
		err := rows.Scan(&s.ID, &s.OwnerID, &s.DayOfWeek, &windowsJSON, &breaksJSON, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to scan day schedule: %w", err)
		}
		if err := json.Unmarshal(windowsJSON, &s.Windows); err != nil {
			return nil, fmt.Errorf("failed to unmarshal windows: %w", err)
		}
		if len(breaksJSON) > 0 {
			if err := json.Unmarshal(breaksJSON, &s.Breaks); err != nil {
				return nil, fmt.Errorf("failed to unmarshal breaks: %w", err)
			}
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *ScheduleRepository) UpsertDaySchedule(ctx context.Context, schedule *models.DaySchedule) error {
	windowsJSON, err := json.Marshal(schedule.Windows)
	if err != nil {
		return fmt.Errorf("failed to marshal windows: %w", err)
	}
	breaksJSON, err := json.Marshal(schedule.Breaks)
	if err != nil {
		return fmt.Errorf("failed to marshal breaks: %w", err)
	}
	if schedule.ID == "" {
		schedule.ID = uuid.New().String()
	}
	query := `
		INSERT INTO day_schedules (id, owner_id, day_of_week, windows, breaks, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (owner_id, day_of_week) DO UPDATE SET
			windows = EXCLUDED.windows,
			breaks = EXCLUDED.breaks,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	err = db.Pool.QueryRow(ctx, query, schedule.ID, schedule.OwnerID, schedule.DayOfWeek, windowsJSON, breaksJSON).
		Scan(&schedule.ID, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to upsert day schedule: %w", err)
	}
	return nil
}

func (r *ScheduleRepository) GetDateExceptions(ctx context.Context, ownerID string) ([]models.DateException, error) {
	query := `
		SELECT id, owner_id, date, exception_type, windows, breaks, description, created_at, updated_at
		FROM date_exceptions WHERE owner_id = $1 ORDER BY date
	`
	rows, err := db.Pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get date exceptions: %w", err)
	}
	defer rows.Close()

	var exceptions []models.DateException
	for rows.Next() {
		var e models.DateException
		var windowsJSON, breaksJSON []byte
		var description *string
		err := rows.Scan(&e.ID, &e.OwnerID, &e.Date, &e.ExceptionType, &windowsJSON, &breaksJSON, &description, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to scan date exception: %w", err)
		}
		if description != nil {
			e.Description = *description
		}
		if len(windowsJSON) > 0 && windowsJSON[0] == '[' {
			if err := json.Unmarshal(windowsJSON, &e.Windows); err != nil {
				return nil, fmt.Errorf("failed to unmarshal windows: %w", err)
			}
		}
		if len(breaksJSON) > 0 && breaksJSON[0] == '[' {
			if err := json.Unmarshal(breaksJSON, &e.Breaks); err != nil {
				return nil, fmt.Errorf("failed to unmarshal breaks: %w", err)
			}
		}
		exceptions = append(exceptions, e)
	}
	return exceptions, nil
}

func (r *ScheduleRepository) UpsertDateException(ctx context.Context, exception *models.DateException) error {
	var windowsJSON, breaksJSON []byte
	var err error
	if len(exception.Windows) > 0 {
		windowsJSON, err = json.Marshal(exception.Windows)
		if err != nil {
			return fmt.Errorf("failed to marshal windows: %w", err)
		}
	}
	if len(exception.Breaks) > 0 {
		breaksJSON, err = json.Marshal(exception.Breaks)
		if err != nil {
			return fmt.Errorf("failed to marshal breaks: %w", err)
		}
	}
	if exception.ID == "" {
		exception.ID = uuid.New().String()
	}
	query := `
		INSERT INTO date_exceptions (id, owner_id, date, exception_type, windows, breaks, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (owner_id, date) DO UPDATE SET
			exception_type = EXCLUDED.exception_type,
			windows = EXCLUDED.windows,
			breaks = EXCLUDED.breaks,
			description = EXCLUDED.description,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	err = db.Pool.QueryRow(ctx, query, exception.ID, exception.OwnerID, exception.Date, exception.ExceptionType,
		windowsJSON, breaksJSON, toNullString(exception.Description)).
		Scan(&exception.ID, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to upsert date exception: %w", err)
	}
	return nil
}

func (r *ScheduleRepository) DeleteDateException(ctx context.Context, ownerID, date string) error {
	query := `DELETE FROM date_exceptions WHERE owner_id = $1 AND date = $2`
	result, err := db.Pool.Exec(ctx, query, ownerID, date)
	if err != nil {
		return fmt.Errorf("failed to delete date exception: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("date exception not found")
	}
	return nil
}

func (r *ScheduleRepository) GetByDate(ctx context.Context, ownerID string, date string) (*models.DateException, error) {
	query := `
		SELECT id, owner_id, date, exception_type, windows, breaks, description, created_at, updated_at
		FROM date_exceptions WHERE owner_id = $1 AND date = $2
	`
	var e models.DateException
	var windowsJSON, breaksJSON []byte
	var description *string
	err := db.Pool.QueryRow(ctx, query, ownerID, date).Scan(
		&e.ID, &e.OwnerID, &e.Date, &e.ExceptionType, &windowsJSON, &breaksJSON, &description, nil, nil)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get date exception: %w", err)
	}
	if description != nil {
		e.Description = *description
	}
	if len(windowsJSON) > 0 && windowsJSON[0] == '[' {
		json.Unmarshal(windowsJSON, &e.Windows)
	}
	if len(breaksJSON) > 0 && breaksJSON[0] == '[' {
		json.Unmarshal(breaksJSON, &e.Breaks)
	}
	return &e, nil
}

func (r *ScheduleRepository) DeleteDaySchedules(ctx context.Context, ownerID string, dayOfWeek int) error {
	query := `DELETE FROM day_schedules WHERE owner_id = $1 AND day_of_week = $2`
	_, err := db.Pool.Exec(ctx, query, ownerID, dayOfWeek)
	return err
}

func toNullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
