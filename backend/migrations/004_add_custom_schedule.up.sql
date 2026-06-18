-- Migration: 004_add_custom_schedule.up.sql
-- Add custom schedule tables for flexible availability windows, breaks, and exceptions

-- Day schedules (schedule for each day of week)
CREATE TABLE IF NOT EXISTS day_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    day_of_week INTEGER NOT NULL CHECK (day_of_week >= 0 AND day_of_week <= 6),
    windows JSONB NOT NULL DEFAULT '[]',
    breaks JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(owner_id, day_of_week)
);

CREATE INDEX IF NOT EXISTS idx_day_schedules_owner_id ON day_schedules(owner_id);

-- Date exceptions (custom schedule or holidays for specific dates)
CREATE TABLE IF NOT EXISTS date_exceptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    exception_type VARCHAR(20) NOT NULL CHECK (exception_type IN ('custom', 'holiday')),
    windows JSONB,
    breaks JSONB,
    description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(owner_id, date)
);

CREATE INDEX IF NOT EXISTS idx_date_exceptions_owner_id ON date_exceptions(owner_id);
CREATE INDEX IF NOT EXISTS idx_date_exceptions_date ON date_exceptions(date);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE TRIGGER update_day_schedules_updated_at
    BEFORE UPDATE ON day_schedules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE OR REPLACE TRIGGER update_date_exceptions_updated_at
    BEFORE UPDATE ON date_exceptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();