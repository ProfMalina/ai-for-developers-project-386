# Documentation vs Codebase Audit - Issues Requiring Human Attention

## Summary
Audit of documentation against actual codebase revealed several discrepancies that require human attention.

## Issues Requiring Human Action

### 1. Root README.md is extremely minimal (4 lines)
**File:** `README.md`
**Problem:** Only contains 4 lines - hexlet badge and a URL `https://aiqwen.77.37.240.50.sslip.io/`. This URL may be outdated, broken, or refer to a deployed instance that's no longer available.
**Recommendation:** Either remove the URL, update it to current deployment, or expand the README with actual project information.

### 2. SPEC.MD features not implemented
**Files:** `SPEC.MD`, `frontend/e2e/README.md`
**Problem:** The SPEC.MD lists three requirements that are explicitly marked as not implemented:
- Cookie consent banner (line 100-101 in SPEC.MD)
- Language switching i18n - Russian/English (line 100)
- Theme switching - light/dark/system (line 101)

The `frontend/e2e/README.md` also documents these as placeholder tests (lines 175-180).
**Recommendation:** Either implement these features or update SPEC.MD to reflect the current MVP scope.

### 3. RULES.md referenced but doesn't exist
**File:** `AGENTS.md:79`
**Problem:** AGENTS.md line 79 says "RULES.md is referenced by older guidance but does not exist in the repo." This is already noted but creates confusion.
**Recommendation:** Remove references to RULES.md from documentation or create the file if it's actually needed.

### 4. Backend documentation mentions additional Make targets not in Makefile
**File:** `backend/README.md:115-134`
**Problem:** Backend README mentions targets like `backend-build`, `backend-run`, `backend-fmt`, etc. Most are in the root Makefile, but the documentation structure doesn't clearly indicate this.
**Recommendation:** This is minor - just ensure users know to use root `make` targets.

## Already Fixed (in associated PR)
- Go version in backend/README.md changed from "1.25+" to "1.24+" to match installed Go version
- API method in frontend/README.md corrected from "PUT" to "PATCH" to match actual router