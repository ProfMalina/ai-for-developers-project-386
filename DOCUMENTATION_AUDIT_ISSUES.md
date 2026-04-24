## Summary

Documentation audit completed. Found several inconsistencies between docs and actual code.

## Issues Fixed in PR

| Issue | File | Fix |
|-------|------|-----|
| Wrong HTTP method (PUT vs PATCH) for update | `frontend/README.md:141` | Changed to `PATCH` |

## Issues Requiring Human Attention

### 1. Test count mismatch in E2E README
**File:** `frontend/e2e/README.md:94`  
**Issue:** Claims "138 tests (across 3 browsers)" but actual test count cannot be verified reliably. E2E test suite has ~41 test definitions but the count includes variations across browsers (Chromium, Firefox, WebKit).  
**Action:** Verify actual test count and update documentation.

### 2. Orphaned /api/owners endpoints in code
**Files:** `backend/internal/handlers/owner_handler.go`, `backend/internal/handlers/handler_test.go`, `backend/internal/handlers/handler_success_test.go`  
**Issue:** Test code contains `/api/owners` endpoints (Create, GetByID) that don't exist in router.go. This appears to be dead code from an earlier design.  
**Action:** Review and delete orphaned test handlers or implement the missing endpoints.

### 3. Coverage floor discrepancy
**Files:** 
- `backend/README.md:156` - says "30%"  
- `AGENTS.md` - says "18%"  
- `backend-ci.yml:81` - enforces "30%" in CI  
**Issue:** Backend AGENTS.md contradicts the actual CI enforcement.  
**Action:** Update `AGENTS.md` to match the actual 30% coverage floor.

### 4. Missing docs/ directory
**Issue:** `AGENTS.md` references `docs/` in the "anti-patterns" section but no such directory exists in the repo.  
**Action:** Either create the docs directory or clarify in AGENTS.md that it doesn't exist.

## Minor Notes

- All API endpoints documented in backend/README.md match actual router.go
- Frontend tech stack description is accurate
- TypeSpec contract correctly describes PATCH for updates
- Go version requirement (1.25+) matches go.mod

## Recommendation

1. Merge the PR that fixes the HTTP method
2. Review the 4 pending items above
3. Consider adding a docs/ directory with architecture overview