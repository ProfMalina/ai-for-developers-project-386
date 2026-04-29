## Prioritized Test Coverage Gaps

### Priority 1: Handler Methods (High Dependents + Recent Changes)

| File | Function/Method | Priority | Reason |
|------|-----------------|----------|--------|
| `handlers/public_event_type_handler.go` | `List()`, `GetByID()`, `GetSlots()` | HIGH | Public API, high dependents, no integration test |
| `handlers/booking_handler.go` | `List()`, `GetByID()`, `Cancel()` | HIGH | Core booking flow, called by frontend |
| `handlers/event_type_handler.go` | `GetByID()`, `List()`, `Update()`, `Delete()` | HIGH | CRUD operations, owner API |
| `handlers/owner_handler.go` | `GetByID()` | HIGH | Owner lookup, used across flows |
| `handlers/time_slot_handler.go` | `List()`, `GenerateSlots()` | HIGH | Slot management, core feature |

### Priority 2: Helper Functions

| File | Function | Priority | Reason |
|------|----------|----------|--------|
| `handlers/helpers.go` | `ValidationError()` | MEDIUM | Used in validation responses |
| `handlers/helpers.go` | `InvalidTime()` | MEDIUM | Called by booking handlers |
| `handlers/helpers.go` | `Conflict()` | MEDIUM | Booking conflict handling |

### Priority 3: Data Models (Low)

| File | Priority | Notes |
|--------|----------|-------|
| `models/*.go` | LOW | Usually covered by integration tests |

### Coverage Status

- **Backend services**: Good coverage (18%+ from CI)
- **Frontend**: Good coverage
- **Handlers**: Partial - constructor tests exist but integration/method tests missing