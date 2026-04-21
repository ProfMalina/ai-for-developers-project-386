package memory

import (
	"context"
	"slices"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

type ScheduleRepository struct {
	store *Store
}

func NewScheduleRepository(store *Store) *ScheduleRepository {
	return &ScheduleRepository{store: store}
}

func (r *ScheduleRepository) GetDaySchedules(ctx context.Context, ownerID string) ([]models.DaySchedule, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	return r.store.daySchedules[ownerID], nil
}

func (r *ScheduleRepository) UpsertDaySchedule(ctx context.Context, schedule *models.DaySchedule) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	schedules := r.store.daySchedules[schedule.OwnerID]
	idx := slices.IndexFunc(schedules, func(s models.DaySchedule) bool {
		return s.DayOfWeek == schedule.DayOfWeek
	})
	if idx >= 0 {
		schedules[idx] = *schedule
	} else {
		schedules = append(schedules, *schedule)
	}
	r.store.daySchedules[schedule.OwnerID] = schedules
	return nil
}

func (r *ScheduleRepository) GetDateExceptions(ctx context.Context, ownerID string) ([]models.DateException, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	return r.store.exceptions[ownerID], nil
}

func (r *ScheduleRepository) UpsertDateException(ctx context.Context, exception *models.DateException) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	exceptions := r.store.exceptions[exception.OwnerID]
	idx := slices.IndexFunc(exceptions, func(e models.DateException) bool {
		return e.Date == exception.Date
	})
	if idx >= 0 {
		exceptions[idx] = *exception
	} else {
		exceptions = append(exceptions, *exception)
	}
	r.store.exceptions[exception.OwnerID] = exceptions
	return nil
}

func (r *ScheduleRepository) DeleteDateException(ctx context.Context, ownerID, date string) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	exceptions := r.store.exceptions[ownerID]
	idx := slices.IndexFunc(exceptions, func(e models.DateException) bool {
		return e.Date == date
	})
	if idx < 0 {
		return nil
	}
	r.store.exceptions[ownerID] = append(exceptions[:idx], exceptions[idx+1:]...)
	return nil
}

func (r *ScheduleRepository) GetByDate(ctx context.Context, ownerID string, date string) (*models.DateException, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, e := range r.store.exceptions[ownerID] {
		if e.Date == date {
			return &e, nil
		}
	}
	return nil, nil
}

func (r *ScheduleRepository) DeleteDaySchedules(ctx context.Context, ownerID string, dayOfWeek int) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	schedules := r.store.daySchedules[ownerID]
	schedules = slices.DeleteFunc(schedules, func(s models.DaySchedule) bool {
		return s.DayOfWeek == dayOfWeek
	})
	r.store.daySchedules[ownerID] = schedules
	return nil
}
