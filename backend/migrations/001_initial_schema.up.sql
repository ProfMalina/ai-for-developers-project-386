-- Migration: 001_initial_schema.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Owners table
CREATE TABLE owners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    timezone VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Event types table
CREATE TABLE event_types (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES owners(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes >= 5 AND duration_minutes <= 1440),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Time slots table
CREATE TABLE time_slots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type_id UUID NOT NULL REFERENCES event_types(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_time_range CHECK (end_time > start_time)
);

-- Bookings table
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type_id UUID NOT NULL REFERENCES event_types(id) ON DELETE CASCADE,
    slot_id UUID REFERENCES time_slots(id) ON DELETE SET NULL,
    guest_name VARCHAR(100) NOT NULL,
    guest_email VARCHAR(255) NOT NULL,
    timezone VARCHAR(50),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_booking_time_range CHECK (end_time > start_time)
);

-- Create index for overlapping booking check
CREATE INDEX idx_bookings_time_range ON bookings (start_time, end_time);

-- Create function to prevent overlapping bookings
CREATE OR REPLACE FUNCTION check_booking_overlap() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM bookings
        WHERE id != COALESCE(NEW.id, '00000000-0000-0000-0000-000000000000')
        AND status != 'cancelled'
        AND NEW.start_time < bookings.end_time
        AND NEW.end_time > bookings.start_time
    ) THEN
        RAISE EXCEPTION 'Booking overlaps with existing booking';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for booking overlap check
CREATE TRIGGER prevent_booking_overlap
    BEFORE INSERT OR UPDATE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION check_booking_overlap();

-- Slot generation configs table
CREATE TABLE slot_generation_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES owners(id) ON DELETE CASCADE UNIQUE,
    working_hours_start TIME NOT NULL,
    working_hours_end TIME NOT NULL,
    interval_minutes INTEGER NOT NULL DEFAULT 30 CHECK (interval_minutes IN (15, 30)),
    days_of_week INTEGER[] NOT NULL DEFAULT '{1,2,3,4,5}',
    date_from DATE NOT NULL DEFAULT (CURRENT_DATE + INTERVAL '1 day'),
    date_to DATE NOT NULL DEFAULT (CURRENT_DATE + INTERVAL '31 days'),
    timezone VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_working_hours CHECK (working_hours_end > working_hours_start),
    CONSTRAINT valid_date_range CHECK (date_to > date_from)
);

-- Create indexes for better query performance
CREATE INDEX idx_event_types_owner_id ON event_types(owner_id);
CREATE INDEX idx_time_slots_event_type_id ON time_slots(event_type_id);
CREATE INDEX idx_time_slots_availability ON time_slots(is_available);
CREATE INDEX idx_bookings_event_type_id ON bookings(event_type_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_slot_configs_owner_id ON slot_generation_configs(owner_id);

-- Insert default owner for development
INSERT INTO owners (id, name, email, timezone, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Default Owner',
    'owner@example.com',
    'Europe/Moscow',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;
