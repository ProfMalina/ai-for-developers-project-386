package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

type ScheduleService struct {
	repo ScheduleRepository
}

func NewScheduleService(repo ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) GetSchedule(ctx context.Context, ownerID string) (*models.CustomSchedule, error) {
	daySchedules, err := s.repo.GetDaySchedules(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get day schedules: %w", err)
	}
	exceptions, err := s.repo.GetDateExceptions(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get date exceptions: %w", err)
	}
	return &models.CustomSchedule{
		OwnerID:      ownerID,
		DaySchedules: daySchedules,
		Exceptions:   exceptions,
	}, nil
}

func (s *ScheduleService) UpsertSchedule(ctx context.Context, ownerID string, req models.UpsertScheduleRequest) (*models.CustomSchedule, error) {
	for i := range req.DaySchedules {
		req.DaySchedules[i].OwnerID = ownerID
		if err := s.repo.UpsertDaySchedule(ctx, &req.DaySchedules[i]); err != nil {
			return nil, fmt.Errorf("failed to upsert day schedule: %w", err)
		}
	}
	for i := range req.Exceptions {
		req.Exceptions[i].OwnerID = ownerID
		if err := s.repo.UpsertDateException(ctx, &req.Exceptions[i]); err != nil {
			return nil, fmt.Errorf("failed to upsert date exception: %w", err)
		}
	}
	return s.GetSchedule(ctx, ownerID)
}

func (s *ScheduleService) UpsertDaySchedules(ctx context.Context, ownerID string, schedules []models.DaySchedule) (*models.CustomSchedule, error) {
	for i := range schedules {
		schedules[i].OwnerID = ownerID
		if err := s.repo.UpsertDaySchedule(ctx, &schedules[i]); err != nil {
			return nil, fmt.Errorf("failed to upsert day schedule: %w", err)
		}
	}
	return s.GetSchedule(ctx, ownerID)
}

func (s *ScheduleService) UpsertDateExceptions(ctx context.Context, ownerID string, exceptions []models.DateException) (*models.CustomSchedule, error) {
	for i := range exceptions {
		exceptions[i].OwnerID = ownerID
		if err := s.repo.UpsertDateException(ctx, &exceptions[i]); err != nil {
			return nil, fmt.Errorf("failed to upsert date exception: %w", err)
		}
	}
	return s.GetSchedule(ctx, ownerID)
}

func (s *ScheduleService) DeleteDateException(ctx context.Context, ownerID, date string) error {
	return s.repo.DeleteDateException(ctx, ownerID, date)
}

func (s *ScheduleService) IsSlotAvailable(ctx context.Context, ownerID string, slotStart, slotEnd time.Time) (bool, error) {
	date := slotStart.Format("2006-01-02")
	exception, err := s.repo.GetByDate(ctx, ownerID, date)
	if err != nil {
		return false, fmt.Errorf("failed to check date exception: %w", err)
	}
	if exception != nil {
		if exception.ExceptionType == models.ExceptionTypeHoliday {
			return false, nil
		}
		if exception.ExceptionType == models.ExceptionTypeCustom {
			return s.checkWindowFit(slotStart, slotEnd, exception.Windows, exception.Breaks), nil
		}
	}
	daySchedules, err := s.repo.GetDaySchedules(ctx, ownerID)
	if err != nil {
		return false, fmt.Errorf("failed to get day schedules: %w", err)
	}
	weekday := int(slotStart.Weekday())
	for _, ds := range daySchedules {
		if ds.DayOfWeek == weekday {
			return s.checkWindowFit(slotStart, slotEnd, ds.Windows, ds.Breaks), nil
		}
	}
	return false, nil
}

func (s *ScheduleService) checkWindowFit(slotStart, slotEnd time.Time, windows []models.ScheduleWindow, breaks []models.ScheduleBreak) bool {
	slotDate := slotStart.Format("2006-01-02")
	for _, w := range windows {
		winStart, err1 := time.Parse("2006-01-02 15:04", slotDate+" "+w.StartTime)
		winEnd, err2 := time.Parse("2006-01-02 15:04", slotDate+" "+w.EndTime)
		if err1 != nil || err2 != nil {
			continue
		}
		if !slotStart.Before(winStart) && !slotEnd.After(winEnd) {
			for _, br := range breaks {
				breakStart, err1 := time.Parse("2006-01-02 15:04", slotDate+" "+br.StartTime)
				breakEnd, err2 := time.Parse("2006-01-02 15:04", slotDate+" "+br.EndTime)
				if err1 != nil || err2 != nil {
					continue
				}
				if slotStart.Before(breakEnd) && slotEnd.After(breakStart) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (s *ScheduleService) FilterSlots(ctx context.Context, ownerID string, slots []models.TimeSlot) ([]models.TimeSlot, error) {
	var filtered []models.TimeSlot
	for _, slot := range slots {
		available, err := s.IsSlotAvailable(ctx, ownerID, slot.StartTime, slot.EndTime)
		if err != nil {
			return nil, err
		}
		if available {
			filtered = append(filtered, slot)
		}
	}
	return filtered, nil
}
