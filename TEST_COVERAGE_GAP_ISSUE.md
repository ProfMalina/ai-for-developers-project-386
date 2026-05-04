## Test Coverage Gaps - Prioritized List

Based on codebase analysis, here are modules lacking test coverage, prioritized by complexity, dependents, and business impact.

### PRIORITY 1 - High (High Complexity + High Dependents)

| Module | Lines | Untested Methods | Complexity | Dependents |
|--------|-------|-----------------|------------|-------------|
| `backend/internal/handlers/booking_handler.go` | 97 | List, GetByID, Cancel | High | BookingService, Owner API |
| `backend/internal/handlers/time_slot_handler.go` | 144 | List, GenerateSlots | High | TimeSlotService, Slot generation |
| `backend/internal/app/container.go` | ~150 | DI wiring | Very High | All handlers, router |

**Rationale:** These handle core booking operations, payment-relevant cancellation logic, and slot generation - critical business workflows.

### PRIORITY 2 - Medium (Moderate Complexity)

| Module | Lines | Untested Methods | Complexity | Dependents |
|--------|-------|-----------------|------------|-------------|
| `backend/internal/handlers/public_event_type_handler.go` | 104 | List, GetByID, GetSlots | Medium | Guest API, public booking |
| `frontend/src/components/owner/BookingsList.tsx` | ~200 | Edge cases, error states | Medium | Owner dashboard |

**Rationale:** Public API endpoints for guest booking - important for conversion funnel.

### PRIORITY 3 - Lower (Utility/Low Risk)

| Module | Notes |
|--------|-------|
| `backend/internal/app/mode.go` | App mode config |
| `frontend/src/utils/validation.ts` | Expand edge case tests |
| `frontend/src/utils/slots.ts` | Expand boundary tests |

---

## Coverage Summary

- **Backend Handlers**: 6 total, 3 fully tested, 3 partially/untested
- **Backend Services**: All have unit tests
- **Frontend Components**: 4 owner components, 3 have tests
- **Frontend Utils**: Partial coverage

## Recommendations

1. Start with `booking_handler.go` - covers owner booking management (cancel, list, get)
2. Then `time_slot_handler.go` - slot generation is core value proposition
3. Then `public_event_type_handler.go` - guest experience matters for bookings
4. Finally `container.go` - integration tests for DI wiring

These gaps represent ~40% of HTTP handler surface area and key guest-facing flows.