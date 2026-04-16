# FRONTEND E2E KNOWLEDGE BASE

## OVERVIEW
Playwright browser suite for end-to-end and API-integration flows. Uses page objects, fixture data, screenshots/video/trace on failure, and runs across Chromium, Firefox, and WebKit.

## STRUCTURE
```text
e2e/
├── fixtures/   # canned test data
├── pages/      # page objects
├── specs/      # guest, owner, common, api-integration flows
├── utils/      # shared helpers
└── tsconfig.json
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Shared page methods | `pages/BasePage.ts` | common browser actions |
| Guest journeys | `specs/guest-flow.spec.ts` | booking happy-path + validation |
| Owner journeys | `specs/owner-flow.spec.ts` | CRUD, slots, pagination |
| Cross-cutting UI | `specs/common.spec.ts` | navigation, responsive behavior, placeholders |
| API failures | `specs/api-integration.spec.ts` | error handling and mocking |

## CONVENTIONS
- Keep Page Object Model intact; add behavior to page objects before repeating selectors in specs.
- Default command path is via root Make targets: `make frontend-e2e`, `make frontend-e2e-chromium`, `make frontend-e2e-debug`, `make frontend-e2e-report`.
- Local config auto-starts the Vite dev server; CI uses retries and uploads report/screenshot artifacts.
- Placeholder coverage is explicit: cookie banner, i18n, and theme switching are not yet implemented in the app.

## ANTI-PATTERNS
- Do not edit `playwright-report/` or `test-results/`; they are disposable outputs.
- Do not convert placeholder tests into fake passing checks without the app feature existing.
- Do not scatter raw selectors everywhere when a page object already owns that screen.
- Do not assume backend-dependent cases work in isolated mode; some conflict/cancel flows still require real data/setup.

## COMMANDS
```bash
make install-e2e
make frontend-e2e
make frontend-e2e-chromium
make frontend-e2e-debug
make frontend-e2e-report
```

## NOTES
- Current README claims 138 tests across 3 browsers.
- CI workflow lives at `.github/workflows/playwright.yml` and runs only for frontend changes on `main`/`master`.
