# Documentation vs Codebase Review - Issues Requiring Human Attention

## Summary
This issue tracks feature gaps documented in SPEC.MD that are not yet implemented in the codebase. These are **not** documentation bugs, but rather documented requirements that need implementation work.

## SPEC.MD Requirements Not Yet Implemented

### 1. Language Switching (i18n)
- **Reference**: SPEC.MD line 100: "Сайт должен поддерживать выбор языка между русским и английским"
- **Current Status**: Not implemented
- **Impact**: Tests in `e2e/specs/common.spec.ts` (lines 241, 261, 270) are placeholders expecting this feature

### 2. Theme Switching
- **Reference**: SPEC.MD line 101: "Сайт должен поддерживать выбор темы и использовать по умолчанию цветовую тему пользователя"
- **Current Status**: Only system color scheme detection works; manual light/dark toggle not implemented
- **Impact**: Tests in `e2e/specs/common.spec.ts` (lines 219, 229) show only system preference detection, not manual switching

### 3. Cookie Consent Banner
- **Reference**: SPEC.MD line 103: "Необходимы предупреждения о куки согласно локальному законодательству"
- **Current Status**: Not implemented
- **Impact**: Test in `e2e/specs/common.spec.ts` (line 206) is a placeholder

## Recommendations

These features are documented in SPEC.MD but have no implementation. Options:

1. **Implement the features** - Add i18n, theme switching, and cookie consent to the frontend
2. **Update SPEC.MD** - Remove these requirements if they're no longer desired
3. **Create separate tracking issues** - Break these into individual implementation tasks

## Notes

- E2E tests correctly mark these as pending/placeholder tests
- The frontend currently works but lacks these specific features
- No other documentation vs code mismatches were found that require human action