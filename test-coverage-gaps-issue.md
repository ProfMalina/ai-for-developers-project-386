# Test Coverage Gaps - Prioritized List

## Summary

This issue identifies modules and functions that lack test coverage, prioritized by:
1. **Number of dependents** - how many other modules depend on this
2. **Recent change frequency** - how often the module has been modified
3. **Complexity** - business logic complexity and size

---

## Priority 1: High Impact - Handlers (HIGH DEPENDENTS)

All 7 handler files are completely untested. They are the HTTP layer that all API clients depend on.

| File | Dependents | Complexity | Notes |
|------|------------|------------|-------|
| `handlers/public_booking_handler.go` | Public API clients, highest traffic | High | 104 lines, complex validation logic |
| `handlers/public_event_type_handler.go` | Public API clients, highest traffic | High | 104 lines, 3 endpoints |
| `handlers/booking_handler.go` | Owner dashboard, bookings UI | Medium | 97 lines |
| `handlers/time_slot_handler.go` | Slot management UI | High | 144 lines, complex query parsing |
| `handlers/event_type_handler.go` | Event type management UI | Medium | 120 lines |
| `handlers/owner_handler.go` | Owner management | Low | 56 lines |
| `handlers/helpers.go` | ALL handlers | High | 90 lines, shared response utilities |

**Why Priority 1:** Handlers translate HTTP to service calls. No tests means:
- No validation that error codes are correct
- No validation that request parsing works
- No integration-level coverage of API contracts

---

## Priority 2: High Impact - Models (UNIVERSAL DEPENDENTS)

All 5 model files are untested. Models are referenced everywhere.

| File | Dependents | Complexity | Notes |
|------|------------|------------|-------|
| `models/common.go` | EVERYTHING | Low | Shared types (Pagination, ErrorResponse) |
| `models/booking.go` | Handlers, services, repositories | Medium | Core domain |
| `models/event_type.go` | Handlers, services, repositories | Medium | Core domain |
| `models/time_slot.go` | Handlers, services, repositories | Medium | Core domain |
| `models/owner.go` | Handlers, services, repositories | Low | Core domain |

**Why Priority 2:** Models define data structures. No tests means:
- No validation that JSON tags are correct
- No validation that binding constraints work
- No protection against breaking changes

---

## Priority 3: Medium Impact - Memory Repositories

Memory repository implementations are used in tests but lack explicit unit tests.

| File | Has Implicit Tests | Complexity | Notes |
|------|-------------------|------------|-------|
| `repositories/memory/booking_repository.go` | Partial | Medium | Covered by service tests via mocks |
| `repositories/memory/event_type_repository.go` | Partial | Medium | Covered by service tests via mocks |
| `repositories/memory/time_slot_repository.go` | Partial | Medium | Covered by service tests via mocks |
| `repositories/memory/slot_config_repository.go` | None | Medium | No dedicated tests |
| `repositories/memory/owner_repository.go` | Partial | Low | Covered by service tests via mocks |
| `repositories/memory/store.go` | None | Low | Core data structure |

**Why Priority 3:** Memory repos are tested implicitly through service tests, but explicit unit tests would:
- Isolate repository logic testing
- Provide faster feedback than integration tests
- Document expected behavior

---

## Current Test Coverage Summary

### Backend
- **Services**: Fully tested ✓ (booking, event_type, owner, time_slot)
- **Repositories**: Integration tests exist ✓
- **Config**: Tested ✓
- **DB**: Tested ✓
- **Middleware**: Tested ✓
- **App (router/container)**: Tested ✓
- **Handlers**: NOT TESTED ✗
- **Models**: NOT TESTED ✗

### Frontend
- **All pages**: Tested ✓
- **All components**: Tested ✓
- **API client**: Tested ✓
- **Utils (slots, validation)**: Tested ✓

---

## Recommendations

1. **Add handler tests first** - Start with `helpers.go` and `public_booking_handler.go` as they have highest impact
2. **Add model tests** - Focus on `common.go` and `booking.go`
3. **Add memory repository tests** - Focus on `slot_config_repository.go` as it's untested
4. **Consider table-driven tests** for handlers to cover multiple scenarios efficiently

---

*Analysis date: 2026-05-14*