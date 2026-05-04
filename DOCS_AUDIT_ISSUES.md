# Documentation Audit Issues - Needs Human Attention

*Generated: 2026-05-04*

## Summary

Documentation audit completed. Some issues were fixed in the accompanying PR, others need human attention.

---

## Fixed in PR (merged):

| File | Issue | Fix |
|------|-------|-----|
| `backend/AGENTS.md` | Coverage floor said 18%, actual CI requires 30% | Changed to 30% |
| `frontend/README.md` | API method listed as PUT for event type update | Changed to PATCH |
| `frontend/README.md` | Missing `SlotGeneration.tsx` in project structure | Added component |

---

## Issues Needing Human Attention:

### 1. Root README.md is minimal (needs rewrite)
- **Current**: Only 4 lines showing Hexlet badge and URL
- **Expected**: Project overview, structure, tech stack, quick start commands

### 2. QWEN.md is a placeholder
- **File**: `QWEN.md` in root
- **Content**: Only `@AGENTS.md` reference
- **Action**: Either remove or clarify purpose

### 3. AGENTS.md stale references
- **File**: `AGENTS.md` line 35 mentions "older docs still say 'no implementation code exists'"
- **Action**: This meta-reference should be removed once root README is updated

### 4. IMPLEMENTATION.md verification
- **File**: `backend/IMPLEMENTATION.md`
- **Items to verify**:
  - Line 172: "35MB server binary" - still accurate?
  - "Next Steps" section (lines 186-198) - still relevant?

### 5. Documentation files that are correct (no action needed)
- ✅ `SPEC.MD` - exists and is comprehensive
- ✅ `backend/README.md` - accurate
- ✅ `frontend/README.md` - after fixes, accurate
- ✅ `frontend/e2e/README.md` - accurate
- ✅ API endpoints - all verified consistent between backend/frontend/typespec
- ✅ Project structure - matches actual files

---

## Verification Commands Used

```bash
# Coverage floor verification
grep -r "coverage" .github/workflows/backend-ci.yml
# Found: line 81 requires 30%

# API endpoint verification
# Checked router.go vs client.ts - all aligned

# File structure verification
ls -la backend/internal/
ls -la frontend/src/
```

---

## How to Create GitHub Issue

Since gh CLI is not authenticated, create an issue manually at:
https://github.com/ProfMalina/ai-for-developers-project-386/issues/new

Use this title:
```
Documentation audit: outdated API descriptions and references
```

Use the content above (minus this "How to Create" section) as the body.