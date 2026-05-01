## Summary

Analysis of the codebase reveals several modules with insufficient test coverage. Prioritization is based on:
1. **Number of dependents** - modules used by many other components
2. **Frequency of recent changes** - modules with active development
3. **Complexity** - line count and number of functions

---

## HIGH PRIORITY

### 1. `backend/internal/handlers/helpers.go`
- **Lines:** 90
- **Functions:** 10 helper functions
- **Dependents:** ALL handlers import this
- **Status:** No dedicated tests
- **Rationale:** Critical utility functions used throughout the entire HTTP layer. Any bug here affects every API endpoint.

### 2. `backend/internal/handlers/public_booking_handler.go`
- **Lines:** 104
- **Functions:** 1 major endpoint (Create booking)
- **Dependents:** Public guest booking flow
- **Status:** No dedicated tests
- **Rationale:** Core public-facing booking endpoint with complex validation logic. High traffic, critical business path.

### 3. `backend/internal/handlers/public_event_type_handler.go`
- **Lines:** 104
- **Functions:** 3 methods (List, GetByID, GetSlots)
- **Dependents:** Public event type listing, slot availability
- **Status:** No dedicated tests
- **Rationale:** Public API for guests to browse available services and slots.

### 4. `backend/internal/handlers/time_slot_handler.go`
- **Lines:** 144
- **Functions:** 2 methods (List, GenerateSlots)
- **Dependents:** Owner slot management, slot generation
- **Status:** No dedicated tests
- **Rationale:** Core owner functionality for managing time slots.

---

## MEDIUM PRIORITY

### 5. `backend/internal/handlers/event_type_handler.go`
- **Lines:** 120
- **Functions:** 5 methods (CRUD + List)
- **Status:** No dedicated tests
- **Rationale:** Full CRUD for event types, important for owner workflow.

### 6. `backend/internal/handlers/owner_handler.go`
- **Lines:** 56
- **Functions:** 2 methods (Create, GetByID)
- **Status:** No dedicated tests
- **Rationale:** Basic owner management.

### 7. `backend/internal/repositories/memory/store.go`
- **Lines:** ~25
- **Functions:** In-memory storage core
- **Dependents:** All memory repositories
- **Status:** No dedicated tests
- **Rationale:** Core data store for in-memory implementations used in tests and dev.

### 8-12. Memory Repositories (partial coverage only)
- `memory/booking_repository.go` - only partial tests
- `memory/event_type_repository.go` - only partial tests
- `memory/owner_repository.go` - only partial tests
- `memory/slot_config_repository.go` - only partial tests
- `memory/time_slot_repository.go` - only partial tests

---

## Recommendations

1. Start with `helpers.go` - lowest level, highest impact
2. Then `public_booking_handler.go` - critical business path
3. Add tests for `public_event_type_handler.go` for public API coverage
4. Add handler tests for owner CRUD operations

Current test coverage is focused on services and some handlers, but the handler layer (especially public endpoints and helpers) needs attention.

---
*Analysis date: 2026-05-01*