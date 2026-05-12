# Prioritized Test Coverage Gaps

## Prioritized by: dependents, recent changes, complexity

---

### P0 - Critical (High Dependents + High Complexity)

| Priority | Module | File | Current Coverage | Rationale |
|----------|--------|------|-------------------|-----------|
| 1 | **TimeSlotHandler** | `backend/internal/handlers/time_slot_handler.go` | 0% direct tests | Handles `/api/slots` (List) and `/api/event-types/:id/slots/generate` - critical for owner slot management. 144 lines, complex query param parsing (5+ date formats), multiple error paths. Used by router, container, and frontend API. |
| 2 | **PublicEventTypeHandler** | `backend/internal/handlers/public_event_type_handler.go` | 0% tests | Guest booking entry point - serves `/api/public/event-types` and `/api/public/slots`. 104 lines, business logic (IsActive filter), two service dependencies. |
| 3 | **Container (NewContainer)** | `backend/internal/app/container.go` | 0% tests | Bootstraps entire application - all services, handlers, repositories. Failure here breaks entire app. |
| 4 | **Frontend API client** | `frontend/src/api/client.ts` | ~7% (1/15 methods) | `ownerApi` (8 methods) + `guestApi` (5 methods) + error classes untested. All frontend modules depend on this. |

---

### P1 - High Priority

| Priority | Module | File | Current Coverage | Notes |
|----------|--------|------|-------------------|-------|
| 5 | **OwnerHandler.GetByID** | `backend/internal/handlers/owner_handler.go` | Partial | `Create` tested, `GetByID` not tested |
| 6 | **EventTypeHandler (Update/Delete)** | `backend/internal/handlers/event_type_handler.go` | Partial | Only `Create` validation tested |
| 7 | **Frontend slots.ts utils** | `frontend/src/utils/slots.ts` | 20% (1/5) | `doSlotsOverlap`, `formatSlotTime`, `formatDate`, `calculateSlotEndTime` untested |
| 8 | **mode.go (StorageMode)** | `backend/internal/app/mode.go` | 0% | Enum used by container, config |

---

### P2 - Medium Priority

| Priority | Module | File | Current Coverage | Notes |
|----------|--------|------|-------------------|-------|
| 9 | **main.go** | `backend/cmd/server/main.go` | 0% | Entry point |
| 10 | **TimeSlotService** | `backend/internal/services/time_slot_service.go` | Partial | Some tests exist but 1 test failing |
| 11 | **App.tsx** | `frontend/src/App.tsx` | Partial | Basic tests exist |
| 12 | **Models** | `backend/internal/models/*.go` | 0% | Data structures |

---

### Current Coverage Stats

```
handlers:    39.8%
app:         87.0%
services:    86.2% (1 test failing)
repositories: 12.1%
middleware:  84.0%
config:      100%
```

---

### Recommended Test Order

1. `TimeSlotHandler` - covers most critical owner workflow
2. `PublicEventTypeHandler` - covers guest booking entry
3. `Container` - covers app wiring
4. `client.ts` API methods - covers frontend integration
5. Other handlers and utilities