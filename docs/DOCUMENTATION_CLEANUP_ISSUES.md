# Documentation Cleanup Issue - Needs Human Attention

## Summary

Documentation review found several discrepancies between README/docs and actual codebase.
Some fixed via PR, others need manual verification.

## Fixed in PR

- `backend/README.md`: Changed Go version from 1.25+ to 1.24+ (line 7)
- `frontend/e2e/README.md`: Changed test count from 138 to 43 (line 94)

## Issues Needing Human Attention

### 1. Backend Coverage Floor Mismatch
- `backend/README.md` (line 158) mentions "30% coverage floor"
- `AGENTS.md` says "Current CI coverage floor is 18%"
- **Action**: Verify actual CI config in `.github/workflows/` and align docs

### 2. Go Version in Makefile Compatibility Checks
- Makefile lines 104-111 check for `golangci-lint built with go1.(25|26|27|28)`
- Current environment is Go 1.24
- **Action**: Update Makefile to support Go 1.24 or verify golangci-lint works

### 3. Placeholder Feature Tests
Frontend E2E README marks these as "not yet implemented" but tests exist:
- Cookie consent banner (test at line 206 in common.spec.ts)
- Language switching (test at line 241 in common.spec.ts)  
- Theme switching (test at line 219 in common.spec.ts)

**Action**: Run tests to verify if they pass or fail:
```bash
make frontend-e2e-chromium
```

### 4. Route Parameter Naming
- `backend/README.md` uses `{id}` (e.g., `/api/event-types/{id}`)
- TypeSpec contract uses `{eventTypeId}`, `{bookingId}` for specificity

**Action**: Decide whether to standardize or keep as-is

### 5. API Contract Validation
- SPEC.MD mentions features that may not be implemented: owner route, timezone support, etc.
- **Action**: Verify these features work as documented

## Verification Commands

```bash
go version                              # Current Go version
make backend-test-coverage              # Check actual coverage
npx playwright test --project=chromium  # Run E2E tests
make backtest                          # Run backend tests + lint
```