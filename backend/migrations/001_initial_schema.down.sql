-- Migration: 001_initial_schema.down.sql

DROP TRIGGER IF EXISTS prevent_booking_overlap ON bookings;
DROP FUNCTION IF EXISTS check_booking_overlap();

DROP TABLE IF EXISTS slot_generation_configs CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS time_slots CASCADE;
DROP TABLE IF EXISTS event_types CASCADE;
DROP TABLE IF EXISTS owners CASCADE;

DROP EXTENSION IF EXISTS "uuid-ossp";
