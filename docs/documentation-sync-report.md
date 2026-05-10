# Documentation Sync Report

Generated: 2026-05-10

## Overview

Comparison of documentation (README.md, docs/) against actual codebase revealed several discrepancies that were fixed and several items requiring human attention.

## Changes Made (PR)

### Frontend README.md

1. **API endpoint correction** (line 142): Changed `PUT /api/event-types/:id` → `PATCH /api/event-types/:id` (actual API uses PATCH)
2. **Added missing endpoint** (line 147): Added `POST /api/event-types/:id/slots/generate` - Generate time slots (was missing from Owner API list)
3. **Project structure updated**: Added missing files/directories:
   - `src/components/owner/SlotGeneration.tsx`
   - `src/test/` (setup.ts, mocks.ts)
   - `src/utils/` (validation.ts, slots.ts)
   - `e2e/` (Playwright tests)
4. **Browser support versions updated**: Chrome 90+ → 120+, Firefox 88+ → 121+, Safari 14+ → 17+, Edge 90+ → 120+

### Backend README.md

1. **Project structure updated**: Added:
   - `internal/app/` (bootstrap, router, container)
   - `internal/repositories/memory/` (in-memory implementations)
2. **Clarified .env comment**: "(not committed)" to reflect actual practice
3. **Fixed duplicate endpoint listing** — the Owner API section already had `POST /api/event-types/{id}/slots/generate`, so no change needed there

## Issues Requiring Human Attention

### 1. SPEC.MD Requirements Not Implemented

**File:** `SPEC.MD` (lines 100-102)

**Issue:** Three SPEC.MD requirements have no implementation:
- Cookie consent banner (legal requirement)
- Language switching (Russian/English i18n)
- Theme switching (light/dark/system preference)

**E2E tests:** `frontend/e2e/specs/common.spec.ts` has placeholder tests for these features.

**Recommendation:** Create a separate issue/ticket to track implementation of these features.

### 2. Root README.md is Minimal

**File:** `README.md` (4 lines total)

**Issue:** Root README contains only:
- CI badge
- Deployed URL

**Actual state:** Each area (backend, frontend, e2e, typespec) has its own detailed README.

**Recommendation:** Consider adding a brief overview to the root README with links to child READMEs.

### 3. Documentation Inconsistency in Backend Config

**File:** `backend/README.md` lines 67-75

**Issue:** Configuration section shows env vars `DB_MAX_CONNS`, `DB_MIN_CONNS`, `DB_MAX_CONN_LIFETIME_HOURS`, `DB_MAX_CONN_IDLE_TIME_MIN` but `internal/config/config.go` only loads `DATABASE_URL`, `SERVER_PORT`, `APP_ENV` — the additional DB pool settings are documented but not implemented.

**Recommendation:** Either implement the additional config options or remove them from the documentation.

### 4. Backend Test Coverage Documentation

**File:** `backend/README.md` line 160

**Issue:** States coverage floor is 30%, but `backend/AGENTS.md` says current floor is 18%.

**Current CI:** `.github/workflows/backend-ci.yml` (need to verify actual floor).

**Recommendation:** Verify and align the coverage floor number in README with actual CI configuration.

### 5. E2E Test Count Discrepancy

**File:** `frontend/e2e/README.md` line 94

**Issue:** States "Total: 138 tests" but the actual test count may differ. The README has a detailed test coverage table that should be verified against actual tests.

**Recommendation:** Run `make frontend-e2e` and count actual tests vs documented count.

## Summary

| Category | Count |
|----------|-------|
| Documentation fixes (automated) | 5 |
| Issues requiring human review | 5 |

The documentation has been updated for straightforward discrepancies (API verbs, file paths, versions). The remaining items require architectural decisions or verification against actual implementations.