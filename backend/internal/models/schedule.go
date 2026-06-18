package models

import "time"

// ScheduleWindow represents a time window within a day
type ScheduleWindow struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ScheduleBreak represents a break period to exclude from availability
type ScheduleBreak struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// DaySchedule represents schedule configuration for a specific day of week
type DaySchedule struct {
	ID        string           `json:"id,omitempty"`
	OwnerID   string           `json:"ownerId,omitempty"`
	DayOfWeek int              `json:"dayOfWeek"`
	Windows   []ScheduleWindow `json:"windows"`
	Breaks    []ScheduleBreak  `json:"breaks,omitempty"`
}

// ExceptionType defines the type of date exception
type ExceptionType string

const (
	ExceptionTypeCustom  ExceptionType = "custom"
	ExceptionTypeHoliday ExceptionType = "holiday"
)

// DateException represents an exception that overrides default schedule for a date
type DateException struct {
	ID            string           `json:"id,omitempty"`
	OwnerID       string           `json:"ownerId,omitempty"`
	Date          string           `json:"date"`
	ExceptionType ExceptionType    `json:"exceptionType"`
	Windows       []ScheduleWindow `json:"windows,omitempty"`
	Breaks        []ScheduleBreak  `json:"breaks,omitempty"`
	Description   string           `json:"description,omitempty"`
	CreatedAt     time.Time        `json:"createdAt,omitempty"`
	UpdatedAt     time.Time        `json:"updatedAt,omitempty"`
}

// CustomSchedule represents the full schedule for an owner
type CustomSchedule struct {
	OwnerID      string          `json:"ownerId"`
	DaySchedules []DaySchedule   `json:"daySchedules"`
	Exceptions   []DateException `json:"exceptions,omitempty"`
}

// UpsertScheduleRequest represents the request to create/update schedule
type UpsertScheduleRequest struct {
	DaySchedules []DaySchedule   `json:"daySchedules"`
	Exceptions   []DateException `json:"exceptions,omitempty"`
}
