---
title: "Documentation vs Codebase Discrepancies - Items Requiring Human Attention"
---

## Summary
Documentation audit completed. Some issues were fixed in PR, others require human attention.

## Fixed in PR
- ✅ Frontend README: Fixed API method `PUT` → `PATCH` for event-type update
- ✅ Frontend README: Added missing `SlotGeneration.tsx` to project structure
- ✅ e2e README: Corrected test count from 138 to 43

## Issues Requiring Human Attention

### 1. SPEC.MD Requirements Not Fully Implemented
**Location:** `SPEC.MD` lines 100-102

The spec lists these as requirements but they are not yet implemented in the codebase:
- Language switching (Russian/English) - e2e tests marked as placeholder
- Theme switching (light/dark/system) - e2e tests marked as placeholder
- Cookie consent banner - e2e tests marked as placeholder

**Suggested action:** Either implement these features or update SPEC.MD to reflect they are lower priority/optional.

### 2. Backend README - Minor Inconsistency
**Location:** `backend/README.md` line 7

Mentions "Go 1.25+" but go.mod specifies exact version "1.25.0". Minor - could clarify.

### 3. API Contract Alignment Notes
The TypeSpec (`typespec/main.tsp`) is correctly used as the source of truth. All documented endpoints match the actual implementation. No further action needed here.

---
*Generated during documentation audit*