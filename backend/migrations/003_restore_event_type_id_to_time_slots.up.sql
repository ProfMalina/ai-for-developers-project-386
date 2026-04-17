-- Migration: 003_restore_event_type_id_to_time_slots.up.sql
-- Restore event_type_id on time_slots so owner slot filtering can follow the TypeSpec contract.

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'time_slots' AND column_name = 'event_type_id'
    ) THEN
        ALTER TABLE time_slots ADD COLUMN event_type_id UUID REFERENCES event_types(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_time_slots_event_type_id ON time_slots(event_type_id);
