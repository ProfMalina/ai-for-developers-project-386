package services

import (
	"context"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
)

type mockScheduleRepo struct {
	daySchedules []models.DaySchedule
	exceptions   []models.DateException
}

func (m *mockScheduleRepo) GetDaySchedules(ctx context.Context, ownerID string) ([]models.DaySchedule, error) {
	return m.daySchedules, nil
}

func (m *mockScheduleRepo) UpsertDaySchedule(ctx context.Context, schedule *models.DaySchedule) error {
	found := false
	for i, s := range m.daySchedules {
		if s.DayOfWeek == schedule.DayOfWeek {
			m.daySchedules[i] = *schedule
			found = true
			break
		}
	}
	if !found {
		m.daySchedules = append(m.daySchedules, *schedule)
	}
	return nil
}

func (m *mockScheduleRepo) GetDateExceptions(ctx context.Context, ownerID string) ([]models.DateException, error) {
	return m.exceptions, nil
}

func (m *mockScheduleRepo) UpsertDateException(ctx context.Context, exception *models.DateException) error {
	found := false
	for i, e := range m.exceptions {
		if e.Date == exception.Date {
			m.exceptions[i] = *exception
			found = true
			break
		}
	}
	if !found {
		m.exceptions = append(m.exceptions, *exception)
	}
	return nil
}

func (m *mockScheduleRepo) DeleteDateException(ctx context.Context, ownerID, date string) error {
	for i, e := range m.exceptions {
		if e.Date == date {
			m.exceptions = append(m.exceptions[:i], m.exceptions[i+1:]...)
			break
		}
	}
	return nil
}

func (m *mockScheduleRepo) GetByDate(ctx context.Context, ownerID string, date string) (*models.DateException, error) {
	for _, e := range m.exceptions {
		if e.Date == date {
			return &e, nil
		}
	}
	return nil, nil
}

func TestIsSlotAvailable_NormalWorkday(t *testing.T) {
	repo := &mockScheduleRepo{
		daySchedules: []models.DaySchedule{
			{
				DayOfWeek: 1,
				Windows: []models.ScheduleWindow{
					{StartTime: "09:00", EndTime: "17:00"},
				},
			},
		},
	}
	service := NewScheduleService(repo)

	slotStart := time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC)
	slotEnd := slotStart.Add(time.Hour)

	available, err := service.IsSlotAvailable(context.Background(), "owner1", slotStart, slotEnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available {
		t.Error("expected slot to be available during normal working hours")
	}
}

func TestIsSlotAvailable_DayWithBreak(t *testing.T) {
	repo := &mockScheduleRepo{
		daySchedules: []models.DaySchedule{
			{
				DayOfWeek: 1,
				Windows: []models.ScheduleWindow{
					{StartTime: "09:00", EndTime: "17:00"},
				},
				Breaks: []models.ScheduleBreak{
					{StartTime: "12:00", EndTime: "13:00"},
				},
			},
		},
	}
	service := NewScheduleService(repo)

	slotStart := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)
	slotEnd := slotStart.Add(time.Hour)

	available, err := service.IsSlotAvailable(context.Background(), "owner1", slotStart, slotEnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if available {
		t.Error("expected slot to be unavailable during break time")
	}
}

func TestIsSlotAvailable_CustomException(t *testing.T) {
	repo := &mockScheduleRepo{
		daySchedules: []models.DaySchedule{
			{
				DayOfWeek: 1,
				Windows: []models.ScheduleWindow{
					{StartTime: "09:00", EndTime: "17:00"},
				},
			},
		},
		exceptions: []models.DateException{
			{
				Date:          "2026-04-20",
				ExceptionType: models.ExceptionTypeCustom,
				Windows: []models.ScheduleWindow{
					{StartTime: "10:00", EndTime: "14:00"},
				},
			},
		},
	}
	service := NewScheduleService(repo)

	slotStart := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	slotEnd := slotStart.Add(time.Hour)

	available, err := service.IsSlotAvailable(context.Background(), "owner1", slotStart, slotEnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available {
		t.Error("expected slot to be available during custom exception window")
	}
}

func TestIsSlotAvailable_Holiday(t *testing.T) {
	repo := &mockScheduleRepo{
		daySchedules: []models.DaySchedule{
			{
				DayOfWeek: 1,
				Windows: []models.ScheduleWindow{
					{StartTime: "09:00", EndTime: "17:00"},
				},
			},
		},
		exceptions: []models.DateException{
			{
				Date:          "2026-04-20",
				ExceptionType: models.ExceptionTypeHoliday,
				Description:   "New Year",
			},
		},
	}
	service := NewScheduleService(repo)

	slotStart := time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC)
	slotEnd := slotStart.Add(time.Hour)

	available, err := service.IsSlotAvailable(context.Background(), "owner1", slotStart, slotEnd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if available {
		t.Error("expected slot to be unavailable on holiday")
	}
}

func TestFilterSlots(t *testing.T) {
	repo := &mockScheduleRepo{
		daySchedules: []models.DaySchedule{
			{
				DayOfWeek: 1,
				Windows: []models.ScheduleWindow{
					{StartTime: "09:00", EndTime: "17:00"},
				},
				Breaks: []models.ScheduleBreak{
					{StartTime: "12:00", EndTime: "13:00"},
				},
			},
		},
	}
	service := NewScheduleService(repo)

	slots := []models.TimeSlot{
		{
			StartTime:   time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC),
			EndTime:     time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC),
			IsAvailable: true,
		},
		{
			StartTime:   time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC),
			EndTime:     time.Date(2026, 4, 20, 13, 0, 0, 0, time.UTC),
			IsAvailable: true,
		},
		{
			StartTime:   time.Date(2026, 4, 20, 14, 0, 0, 0, time.UTC),
			EndTime:     time.Date(2026, 4, 20, 15, 0, 0, 0, time.UTC),
			IsAvailable: true,
		},
	}

	filtered, err := service.FilterSlots(context.Background(), "owner1", slots)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered) != 2 {
		t.Errorf("expected 2 available slots, got %d", len(filtered))
	}
}
