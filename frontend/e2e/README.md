# E2E Tests with Playwright

End-to-end integration tests for the Calendar Booking App frontend.

## Quick Start

```bash
# Install browsers (first time only)
make install-e2e

# Run all E2E tests
make frontend-e2e

# Run with UI
make frontend-e2e-ui

# Run on Chromium only (faster)
make frontend-e2e-chromium

# Run in headed mode (for debugging)
make frontend-e2e-headed

# Run in debug mode
make frontend-e2e-debug

# Open HTML report
make frontend-e2e-report
```

## Directory Structure

```
e2e/
├── fixtures/
│   └── test-data.ts       # Test data and mock fixtures
├── pages/
│   ├── BasePage.ts         # Base page object with common methods
│   ├── GuestHomePage.ts    # Guest home page (/)
│   ├── BookingPage.ts      # Booking flow (/book/:id)
│   └── OwnerDashboard.ts   # Owner dashboard (/owner)
├── specs/
│   ├── guest-flow.spec.ts      # Guest booking flow tests
│   ├── owner-flow.spec.ts      # Owner management flow tests
│   ├── common.spec.ts          # Shared features (i18n, theme, navigation)
│   └── api-integration.spec.ts # API error handling & mocking tests
├── utils/
│   └── helpers.ts          # Shared test utilities
└── tsconfig.json           # TypeScript config for E2E tests
```

## Test Coverage

Coverage notes in this document reflect the current suite inventory and the freshly verified Chromium run.

Current audit baseline:

- `npm run test:e2e -- --list` lists `135 tests in 4 files` across Chromium, Firefox, and WebKit.
- The freshly verified Chromium run currently reports `42 passed`, `3 skipped`, `0 failed`.
- Product-gated placeholders remain out of closure scope until the app implements them.

### Guest Flow Tests (`guest-flow.spec.ts`)
- ✅ Suite includes runnable coverage for viewing event types, opening booking pages, successful booking, navigation, and validation.
- ✅ Guest-flow Chromium coverage is currently green.
- ⚠️ Booking conflict coverage exists in the suite, but it still depends on targeted API-mock setup rather than a real backend data fixture.

### Owner Flow Tests (`owner-flow.spec.ts`)
- ✅ Dashboard, create, edit, delete, slot generation, bookings list, validation, pagination, and cancellation happy-path all have runnable Chromium coverage.
- ✅ Owner-flow Chromium coverage is currently green.

### Common Features Tests (`common.spec.ts`)
- ✅ The suite includes runnable checks for guest and owner navigation, invalid routes, default Russian copy, and smoke coverage that the guest/owner shells remain reachable across the audited viewports.
- ✅ System color-scheme coverage in this suite is limited to load-smoke checks under dark/light emulation, not proof of a user-facing theme switcher.
- ⏳ Cookie consent banner is not implemented in the app yet.
- ⏳ RU/EN switching is not implemented in the app yet.
- ⏳ Explicit theme switching is not implemented in the app yet. Current coverage only checks system preference handling.

### API Integration Tests (`api-integration.spec.ts`)
- ✅ The suite includes API mocking coverage for booking request submission, 400/404/409 handling, guest empty-state fallbacks for failed list loads, a hanging-request shell check, blank-shell known-broken failure-mode coverage for aborted/malformed responses, missing-field tolerance, and generic owner create/cancel error surfaces.
- ✅ The current Chromium run passes this API-integration subset.
- ⚠️ This coverage should still be read as targeted behaviour verification, not as proof that every backend failure path in the product is comprehensively exercised.

**Current listed inventory: 135 tests** (across Chromium, Firefox, and WebKit)

## Page Object Model

All tests follow the Page Object Model pattern for maintainability:

```typescript
// Example usage
const guestHome = new GuestHomePage(page);
await guestHome.goto();
await guestHome.expectLoaded();
await guestHome.bookEventType('Консультация');
```

### Available Page Objects

| Page Object | Methods |
|-------------|---------|
| `BasePage` | `goto()`, `expectHeading()`, `acceptCookies()`, `expectNotification()` |
| `GuestHomePage` | `goto()`, `expectLoaded()`, `bookEventType()`, `expectEventTypeVisible()` |
| `BookingPage` | `selectDate()`, `selectTimeSlot()`, `fillGuestDetails()`, `submitBooking()` |
| `OwnerDashboard` | `switchTab()`, `openAddEventType()`, `fillEventTypeForm()`, `cancelBooking()` |

## Configuration

Playwright is configured in `playwright.config.ts`:

- **Test directory**: `e2e/specs/`
- **Base URL**: `http://localhost:5173` (Vite dev server)
- **Browsers**: Chromium, Firefox, WebKit
- **Retries**: 2 (CI), 0 (local)
- **Reporter**: HTML
- **Screenshots**: On failure
- **Video**: On failure
- **Trace**: On first retry

## CI/CD Integration

The repository already has a Playwright workflow. It runs on push and pull request to `main` or `master` when files under `frontend/**` change.

GitHub Actions workflow: `.github/workflows/playwright.yml`

This workflow is existing infrastructure, not proof that the suite is currently healthy.

Artifacts uploaded:
- `playwright-report` — HTML test report (30 days retention)
- `test-screenshots` — Screenshots on failure (7 days retention)

## Tips

### Debugging a failing test
```bash
# Run in debug mode
make frontend-e2e-debug

# Or run a specific test file
npx playwright test guest-flow.spec.ts --debug

# Run with UI mode
make frontend-e2e-ui
```

### Run specific test
```bash
# By file
npx playwright test guest-flow.spec.ts

# By test name
npx playwright test -g "should create a booking"

# By project (browser)
npx playwright test --project=chromium
```

### Update snapshots
```bash
npx playwright test --update-snapshots
```

## Requirements Not Yet Implemented

The following areas remain placeholders or blocked-by-product because the frontend does not implement them yet:

- **Cookie consent banner** — SPEC.MD requirement
- **Language switching (i18n)** — SPEC.MD requirement (Russian/English)
- **Theme switching** — SPEC.MD requirement (light/dark/system)

These checks should not be counted as implemented passing coverage until the frontend ships the corresponding behavior.

## Troubleshooting

### Browsers not installed
```bash
npx playwright install chromium
```

### System dependencies missing (Linux)
```bash
npx playwright install-deps
# or manually install:
sudo apt-get install libnss3 libatk-bridge2.0-0 libdrm2 libxkbcommon0 libxcomposite1 libxdamage1 libxrandr2 libgbm1 libpango-1.0-0 libcairo2 libasound2 libxshmfence1 libglib2.0-0 libx11-xcb1
```

### Tests fail because dev server is not running
The `playwright.config.ts` automatically starts the Vite dev server via `webServer` config. If you want to reuse an existing server:
```bash
npx playwright test --reuse-existing-server
```
