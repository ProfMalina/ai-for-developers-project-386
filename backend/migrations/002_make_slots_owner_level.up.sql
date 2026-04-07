-- Migration: 002_make_slots_owner_level.sql
-- Remove event_type_id from time_slots - slots belong to owner, not event type

-- Drop old index
DROP INDEX IF EXISTS idx_time_slots_event_type_id;

-- Add owner_id to time_slots (only if it doesn't exist)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'time_slots' AND column_name = 'owner_id'
    ) THEN
        ALTER TABLE time_slots ADD COLUMN owner_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000001';
        ALTER TABLE time_slots ADD CONSTRAINT fk_time_slots_owner FOREIGN KEY (owner_id) REFERENCES owners(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Remove event_type_id column (only if it exists)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'time_slots' AND column_name = 'event_type_id'
    ) THEN
        ALTER TABLE time_slots DROP COLUMN event_type_id;
    END IF;
END $$;

-- Create index for owner
CREATE INDEX IF NOT EXISTS idx_time_slots_owner_id ON time_slots(owner_id);
