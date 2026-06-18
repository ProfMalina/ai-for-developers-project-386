-- Migration: 004_add_custom_schedule.down.sql
-- Drop custom schedule tables

DROP TRIGGER IF EXISTS update_date_exceptions_updated_at ON date_exceptions;
DROP TRIGGER IF EXISTS update_day_schedules_updated_at ON day_schedules;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS date_exceptions;
DROP TABLE IF EXISTS day_schedules;