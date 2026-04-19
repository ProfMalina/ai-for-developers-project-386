package repositories

import (
	"strings"
	"testing"
	"time"
)

func TestBuildDeleteAvailableInRangeQuery_UsesOverlapSemantics(t *testing.T) {
	startTime := time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 4, 20, 17, 0, 0, 0, time.UTC)

	query, args := buildDeleteAvailableInRangeQuery("owner-id", "event-type-id", startTime, endTime)

	if len(args) != 4 {
		t.Fatalf("expected 4 query args, got %d", len(args))
	}

	if !strings.Contains(query, "start_time < $4") || !strings.Contains(query, "end_time > $3") {
		t.Fatalf("expected overlap semantics in query, got %q", query)
	}

	if strings.Contains(query, "start_time >= $3") {
		t.Fatalf("expected query to stop using start_time-only window semantics, got %q", query)
	}
}

func TestBuildDeleteAvailableInRangeQuery_PreservesBookedSlotsAndAllowsLegacyNullEventType(t *testing.T) {
	query, _ := buildDeleteAvailableInRangeQuery("owner-id", "event-type-id", time.Now().UTC(), time.Now().Add(time.Hour).UTC())

	if !strings.Contains(query, "NOT EXISTS") || !strings.Contains(query, "FROM bookings b") {
		t.Fatalf("expected query to preserve slots referenced by active bookings, got %q", query)
	}

	if !strings.Contains(query, "b.status != 'cancelled'") { //nolint:misspell // persisted booking status value
		t.Fatalf("expected query to preserve non-cancelled bookings, got %q", query)
	}

	if !strings.Contains(query, "(event_type_id = $2 OR event_type_id IS NULL)") {
		t.Fatalf("expected query to allow legacy NULL event_type_id rows, got %q", query)
	}
}

func TestBuildGetAvailableSlotsQuery_ExcludesIntervalsWithActiveBookings(t *testing.T) {
	startTime := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 4, 21, 0, 0, 0, 0, time.UTC)

	query, args := buildGetAvailableSlotsQuery("owner-id", 1, 20, &startTime, &endTime)

	if len(args) != 5 {
		t.Fatalf("expected 5 query args, got %d", len(args))
	}

	if !strings.Contains(query, "owner_id = $1") {
		t.Fatalf("expected available-slots query to filter by owner, got %q", query)
	}

	if !strings.Contains(query, "is_available = true") {
		t.Fatalf("expected available-slots query to require available rows, got %q", query)
	}

	if !strings.Contains(query, "NOT EXISTS") || !strings.Contains(query, "FROM bookings b") {
		t.Fatalf("expected available-slots query to exclude overlapping active bookings, got %q", query)
	}

	if !strings.Contains(query, "b.status != 'cancelled'") {
		t.Fatalf("expected available-slots query to ignore cancelled bookings only, got %q", query)
	}

	if !strings.Contains(query, "b.start_time < time_slots.end_time") || !strings.Contains(query, "b.end_time > time_slots.start_time") {
		t.Fatalf("expected overlap-based booking exclusion in available-slots query, got %q", query)
	}
}
