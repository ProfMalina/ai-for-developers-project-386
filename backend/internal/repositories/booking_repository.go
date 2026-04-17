package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// BookingRepository handles database operations for bookings
type BookingRepository struct{}

// NewBookingRepository creates a new booking repository
func NewBookingRepository() *BookingRepository {
	return &BookingRepository{}
}

// Create creates a new booking
func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	query := `
		INSERT INTO bookings (id, event_type_id, slot_id, guest_name, guest_email, timezone, start_time, end_time, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`

	if booking.ID == "" {
		booking.ID = uuid.New().String()
	}

	if booking.Status == "" {
		booking.Status = "confirmed"
	}

	err := db.Pool.QueryRow(ctx, query, booking.ID, booking.EventTypeID, booking.SlotID,
		booking.GuestName, booking.GuestEmail, booking.Timezone, booking.StartTime, booking.EndTime, booking.Status).
		Scan(&booking.ID, &booking.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

// CreateWithReservedSlot atomically reserves a slot and creates a booking.
func (r *BookingRepository) CreateWithReservedSlot(ctx context.Context, booking *models.Booking) error {
	if booking.SlotID == nil {
		return r.Create(ctx, booking)
	}

	if booking.ID == "" {
		booking.ID = uuid.New().String()
	}

	if booking.Status == "" {
		booking.Status = "confirmed"
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin booking transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	reserveQuery := `UPDATE time_slots SET is_available = false WHERE id = $1 AND is_available = true`
	result, err := tx.Exec(ctx, reserveQuery, *booking.SlotID)
	if err != nil {
		return fmt.Errorf("failed to reserve time slot: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("time slot is already booked")
	}

	createQuery := `
		INSERT INTO bookings (id, event_type_id, slot_id, guest_name, guest_email, timezone, start_time, end_time, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		RETURNING id, created_at
	`

	err = tx.QueryRow(ctx, createQuery, booking.ID, booking.EventTypeID, booking.SlotID,
		booking.GuestName, booking.GuestEmail, booking.Timezone, booking.StartTime, booking.EndTime, booking.Status).
		Scan(&booking.ID, &booking.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "overlaps") {
			return fmt.Errorf("selected time slot is already booked")
		}
		return fmt.Errorf("failed to create booking: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit booking transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a booking by ID
func (r *BookingRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	query := `
		SELECT id, event_type_id, slot_id, guest_name, guest_email, timezone, start_time, end_time, status, created_at
		FROM bookings WHERE id = $1
	`

	booking := &models.Booking{}
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&booking.ID, &booking.EventTypeID, &booking.SlotID, &booking.GuestName,
		&booking.GuestEmail, &booking.Timezone, &booking.StartTime, &booking.EndTime,
		&booking.Status, &booking.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return booking, nil
}

// List retrieves a paginated list of bookings with filters and sorting
func (r *BookingRepository) List(ctx context.Context, page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) ([]models.Booking, int, error) {
	query, countQuery, args := buildBookingListQueries(page, pageSize, sortBy, sortOrder, dateFrom, dateTo)

	// Count total items
	var totalItems int
	err := db.Pool.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&totalItems)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(&booking.ID, &booking.EventTypeID, &booking.SlotID, &booking.GuestName,
			&booking.GuestEmail, &booking.Timezone, &booking.StartTime, &booking.EndTime,
			&booking.Status, &booking.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	if bookings == nil {
		bookings = []models.Booking{}
	}

	return bookings, totalItems, nil
}

func buildBookingListQueries(page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) (string, string, []interface{}) {
	if sortBy == "" {
		sortBy = "start_time"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}

	allowedSortFields := map[string]bool{
		"created_at": true,
		"createdAt":  true,
		"start_time": true,
		"startTime":  true,
		"guest_name": true,
		"guestName":  true,
		"status":     true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "start_time"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	sortFieldMap := map[string]string{
		"createdAt": "created_at",
		"startTime": "start_time",
		"guestName": "guest_name",
	}
	if mapped, ok := sortFieldMap[sortBy]; ok {
		sortBy = mapped
	}

	offset := (page - 1) * pageSize

	query := `SELECT id, event_type_id, slot_id, guest_name, guest_email, timezone, start_time, end_time, status, created_at FROM bookings WHERE status != 'cancelled'`
	countQuery := `SELECT COUNT(*) FROM bookings WHERE status != 'cancelled'`
	args := []interface{}{}
	argIdx := 1

	if dateFrom != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND start_time >= $%d", argIdx)
		args = append(args, *dateFrom)
		argIdx++
	}

	if dateTo != nil {
		query += fmt.Sprintf(" AND end_time <= $%d", argIdx)
		countQuery += fmt.Sprintf(" AND end_time <= $%d", argIdx)
		args = append(args, *dateTo)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, sortOrder, argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	return query, countQuery, args
}

// CheckOverlap checks if a new booking would overlap with existing bookings
func (r *BookingRepository) CheckOverlap(ctx context.Context, startTime, endTime time.Time) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE status != 'cancelled'
			AND $1 < end_time
			AND $2 > start_time
		)
	`

	var overlaps bool
	err := db.Pool.QueryRow(ctx, query, startTime, endTime).Scan(&overlaps)
	if err != nil {
		return false, fmt.Errorf("failed to check booking overlap: %w", err)
	}

	return overlaps, nil
}

// Update updates a booking
func (r *BookingRepository) Update(ctx context.Context, booking *models.Booking) error {
	query := `
		UPDATE bookings
		SET guest_name = $2, guest_email = $3, timezone = $4, start_time = $5, end_time = $6, status = $7
		WHERE id = $1
	`

	result, err := db.Pool.Exec(ctx, query, booking.ID, booking.GuestName, booking.GuestEmail,
		booking.Timezone, booking.StartTime, booking.EndTime, booking.Status)
	if err != nil {
		if strings.Contains(err.Error(), "overlaps") {
			return fmt.Errorf("booking overlaps with existing booking")
		}
		return fmt.Errorf("failed to update booking: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// Patch partially updates a booking
func (r *BookingRepository) Patch(ctx context.Context, id string, req models.UpdateBookingRequest) (*models.Booking, error) {
	booking, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.GuestName != nil {
		setClauses = append(setClauses, fmt.Sprintf("guest_name = $%d", argIdx))
		args = append(args, *req.GuestName)
		argIdx++
	}
	if req.GuestEmail != nil {
		setClauses = append(setClauses, fmt.Sprintf("guest_email = $%d", argIdx))
		args = append(args, *req.GuestEmail)
		argIdx++
	}
	if req.Timezone != nil {
		setClauses = append(setClauses, fmt.Sprintf("timezone = $%d", argIdx))
		args = append(args, *req.Timezone)
		argIdx++
	}

	if len(setClauses) == 0 {
		return booking, nil
	}

	args = append(args, id)

	query := fmt.Sprintf("UPDATE bookings SET %s WHERE id = $%d RETURNING id, event_type_id, slot_id, guest_name, guest_email, timezone, start_time, end_time, status, created_at",
		strings.Join(setClauses, ", "), argIdx+1)

	err = db.Pool.QueryRow(ctx, query, args...).Scan(
		&booking.ID, &booking.EventTypeID, &booking.SlotID, &booking.GuestName,
		&booking.GuestEmail, &booking.Timezone, &booking.StartTime, &booking.EndTime,
		&booking.Status, &booking.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "overlaps") {
			return nil, fmt.Errorf("booking overlaps with existing booking")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to patch booking: %w", err)
	}

	return booking, nil
}

// Cancel cancels a booking
func (r *BookingRepository) Cancel(ctx context.Context, id string) error {
	query := `UPDATE bookings SET status = 'cancelled' WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}

// Delete deletes a booking
func (r *BookingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM bookings WHERE id = $1`

	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found")
	}

	return nil
}
