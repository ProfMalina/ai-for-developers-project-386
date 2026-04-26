## Summary

Documentation review complete. Some issues were fixed, others require human attention.

## Fixed in this PR

1. **backend/README.md**: Updated Go version from 1.25+ to 1.24+ (actual installed version)
2. **frontend/README.md**: Fixed API endpoint method from PUT to PATCH for event type updates
3. **frontend/e2e/README.md**: Fixed test count from "138 tests" to "44 tests" (actual count)

## Issues requiring human attention

### 1. Model name mismatch between TypeSpec and Go implementation

The TypeSpec defines:
- `model SlotGenerationConfig` (line 197 in typespec/main.tsp)

The Go backend uses:
- `model SlotGenerationRequest` (backend/internal/models/time_slot.go:39)

**Issue**: Naming drift between API spec and implementation. Should decide which name is canonical and update the other.

### 2. E2E tests with missing backend integration

The e2e README (lines 59, 69, 84) marks tests as requiring backend:
- "⚠️ Booking conflicts (requires backend setup)"
- "⚠️ Cancel booking (requires test data)"
- "⚠️ Booking conflict (409)"

These backend features ARE implemented. Tests need proper setup or mocking.

### 3. Placeholder feature tests in common.spec.ts

Lines 77-80 show placeholder indicators for unimplemented features:
- Cookie consent banner (⏳)
- Language switching i18n (⏳)
- Theme switching (⏳)

Check if these features exist in frontend code and update tests accordingly.

### 4. Browser version claims need verification

frontend/README.md lines 186-190 claim browser support for Chrome 90+, Firefox 88+, Safari 14+, Edge 90+. Verify against playwright.config.ts browser versions.

### 5. Docker Compose path discrepancy

- frontend/README.md line 68: mentions `../typespec/tsp-output/schema/openapi.yaml`
- backend/README.md suggests running from backend/ dir: `cd backend && docker-compose up -d`
- Root Makefile (line 83) uses: `cd $(BACKEND_DIR) && docker-compose up -d`

Ensure docker-compose is at correct path.

## Verification commands

```bash
# Check actual test count
cd frontend && npx playwright test --list | grep -c "test:"

# Check API endpoints in handlers
cd backend && grep -r "@route\|c.Param" internal/handlers/

# Check models
grep -r "type.*Generation" backend/internal/models/
```