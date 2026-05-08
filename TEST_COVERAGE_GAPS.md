# Test Coverage Gaps - Prioritized List

## Summary

This document catalogs modules lacking test coverage, prioritized by: (1) number of dependents (handlers are called by router/container), (2) complexity (lines of code), (3) business criticality.

---

## Backend - Priority 1 (High Traffic / Business Critical)

| Module | Lines | Dependents | Notes |
|--------|-------|------------|-------|
| `handlers/public_booking_handler.go` | 104 | Container, Router | Public booking creation - highest traffic endpoint |
| `handlers/booking_handler.go` | 97 | Container, Router | Owner booking list/get/cancel |
| `handlers/event_type_handler.go` | 120 | Container, Router | Full CRUD for event types |
| `handlers/time_slot_handler.go` | 144 | Container, Router | Slot generation + listing |

---

## Backend - Priority 2 (Lower Traffic / Supporting)

| Module | Lines | Dependents | Notes |
|--------|-------|------------|-------|
| `handlers/public_event_type_handler.go` | 104 | Container, Router | Public event type listing |
| `handlers/owner_handler.go` | 56 | Container, Router | Owner CRUD operations |
| `handlers/helpers.go` | - | All handlers | Utility functions - lower priority |

---

## Backend - Priority 3 (Repositories)

| Module | Lines | Test Status |
|--------|-------|-------------|
| `repositories/slot_config_repository.go` | 75 | No tests |
| `repositories/memory/store.go` | 26 | No tests |
| `repositories/memory/*.go` (all) | No tests |

Note: Main DB repositories have integration tests; memory implementations are untested.

---

## Frontend - Low Priority Gaps

| Module | Status |
|--------|--------|
| `main.tsx` | Entry point - likely not worth testing |
| `types/api.ts` | Type definitions only |

Most frontend components already have test coverage.

---

## Recommended Approach

1. **Start with Priority 1 handlers** - they handle the core API surface
2. Use `handler_test.go` and `handler_success_test.go` as reference patterns
3. Focus on: request validation, error responses, service integration
4. Memory repository tests can be added alongside handler tests using mock stores