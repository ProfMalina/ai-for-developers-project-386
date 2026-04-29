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

### Guest Flow Tests (`guest-flow.spec.ts`)
- ✅ View event types list
- ✅ Select event type and view calendar
- ✅ Create a booking successfully
- ✅ Form validation (empty fields, invalid email)
- ✅ Navigation between pages
- ⚠️ Booking conflicts (requires backend setup)

### Owner Flow Tests (`owner-flow.spec.ts`)
- ✅ View dashboard
- ✅ Create event type
- ✅ Edit event type
- ✅ Delete event type
- ✅ Generate time slots
- ✅ View bookings list
- ✅ Navigate through bookings (pagination)
- ⚠️ Cancel booking (requires test data)
- ✅ Form validation

### Common Features Tests (`common.spec.ts`)
- ✅ Navigation between guest/owner pages
- ✅ 404 page for invalid routes
- ✅ Header navigation highlighting
- ✅ Responsive design (mobile, tablet, desktop)
- ⏳ Cookie consent banner (not yet implemented in app)
- ✅ System color scheme detection
- ⏳ Language switching (not yet implemented in app)
- ⏳ Theme switching (not yet implemented in app)

### API Integration Tests (`api-integration.spec.ts`)
- ✅ Successful booking creation (201)
- ⏳ Booking conflict (409) - requires test data setup
- ✅ Validation error (400)
- ✅ Not found error (404)
- ✅ Server error (500) handling
- ✅ Network timeout
- ✅ API unavailable
- ✅ Malformed JSON response
- ✅ Missing fields in response
- ✅ API mocking for isolated tests

**Total: ~41 tests** (across 3 browsers: Chromium, Firefox, WebKit)

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

Tests run automatically on push/PR to `main`/`master` when frontend files change.

GitHub Actions workflow: `.github/workflows/playwright.yml`

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

## Requirements Not Yet Fully Implemented in App

The following features from SPEC.MD are not fully implemented:

- **Cookie consent banner** — SPEC.MD requirement
- **Theme switching** — SPEC.MD requirement (light/dark/system)
- **Language switching (i18n)** — Only Russian locale is hardcoded; English option required but not implemented

These tests will become functional once the corresponding features are added to the frontend.

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
