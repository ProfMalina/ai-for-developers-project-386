package repositories

import (
	"strings"
	"testing"
	"time"
)

func TestBuildBookingListQueries_ExcludeCancelledByDefault(t *testing.T) {
	query, countQuery, args := buildBookingListQueries(1, 20, "", "", nil, nil)

	if len(args) != 2 {
		t.Fatalf("expected pagination args only, got %d", len(args))
	}

	if !strings.Contains(query, "status != 'cancelled'") { //nolint:misspell // persisted booking status value
		t.Fatalf("expected bookings list query to exclude canceled bookings by default, got %q", query)
	}

	if !strings.Contains(countQuery, "status != 'cancelled'") { //nolint:misspell // persisted booking status value
		t.Fatalf("expected bookings count query to exclude canceled bookings by default, got %q", countQuery)
	}

	if !strings.Contains(query, "ORDER BY start_time asc") {
		t.Fatalf("expected default sort to remain start_time asc, got %q", query)
	}
}

func TestBuildBookingListQueries_PreserveDateFilters(t *testing.T) {
	dateFrom := time.Date(2026, 4, 10, 9, 0, 0, 0, time.UTC)
	dateTo := dateFrom.Add(2 * time.Hour)

	query, countQuery, args := buildBookingListQueries(2, 10, "created_at", "desc", &dateFrom, &dateTo)

	if len(args) != 4 {
		t.Fatalf("expected date filters plus pagination args, got %d", len(args))
	}

	if !strings.Contains(query, "start_time >= $1") || !strings.Contains(query, "end_time <= $2") {
		t.Fatalf("expected query to preserve date filters, got %q", query)
	}

	if !strings.Contains(countQuery, "start_time >= $1") || !strings.Contains(countQuery, "end_time <= $2") {
		t.Fatalf("expected count query to preserve date filters, got %q", countQuery)
	}

	if !strings.Contains(query, "ORDER BY created_at desc LIMIT $3 OFFSET $4") {
		t.Fatalf("expected query to use mapped sort and pagination, got %q", query)
	}
}
