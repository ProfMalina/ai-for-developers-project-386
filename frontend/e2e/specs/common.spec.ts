import { test, expect } from '@playwright/test';
import { GuestHomePage } from '../pages/GuestHomePage';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { setViewport } from '../utils/helpers';

const BOOKING_HEADING = 'Забронировать встречу';

test.describe('Common Features', () => {
  test.beforeEach(async ({ page }) => {
    // Mock API - use correct response structure: { items: [], pagination: {} }
    await page.route('**/api/public/event-types**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: [
            { id: 'consultation', name: 'Консультация', description: 'Individual consultation', durationMinutes: 30, ownerId: '00000000-0000-0000-0000-000000000001', isActive: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-01-01T00:00:00Z' },
            { id: 'meeting', name: 'Встреча', description: 'Group meeting', durationMinutes: 60, ownerId: '00000000-0000-0000-0000-000000000001', isActive: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-01-01T00:00:00Z' },
          ],
          pagination: { page: 1, pageSize: 10, totalItems: 2, totalPages: 1, hasNext: false, hasPrev: false },
        }),
      });
    });

    await page.route('**/api/event-types**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: [
            { id: 'consultation', name: 'Консультация', description: 'Individual', durationMinutes: 30, ownerId: '00000000-0000-0000-0000-000000000001', isActive: true, createdAt: '2026-01-01T00:00:00Z', updatedAt: '2026-01-01T00:00:00Z' },
          ],
          pagination: { page: 1, pageSize: 10, totalItems: 1, totalPages: 1, hasNext: false, hasPrev: false },
        }),
      });
    });

    await page.route('**/api/bookings**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: [],
          pagination: { page: 1, pageSize: 10, totalItems: 0, totalPages: 1, hasNext: false, hasPrev: false },
        }),
      });
    });

    await page.route('**/api/slots**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: [],
          pagination: { page: 1, pageSize: 10, totalItems: 0, totalPages: 1, hasNext: false, hasPrev: false },
        }),
      });
    });
  });

  test.afterEach(async ({ page }) => {
    await page.unroute('**/api/**');
  });

  test.describe('Navigation', () => {
    test('should navigate between guest and owner pages', async ({ page }) => {
      // Debug: log console errors
      page.on('console', msg => { if (msg.type() === 'error') console.log('BROWSER ERROR:', msg.text()); });
      page.on('pageerror', err => console.log('PAGE ERROR:', err.message));

      const guestHome = new GuestHomePage(page);

      // Start at guest home
      console.log('Step 1: Going to guest home...');
      await guestHome.goto();
      console.log('Step 2: After goto, URL:', page.url());

      // Check page content
      const html = await page.content();
      console.log('HTML length:', html.length);
      const rootHTML = await page.locator('#root').innerHTML().catch(() => 'ERROR');
      console.log('Root innerHTML length:', typeof rootHTML === 'string' ? rootHTML.length : 'N/A');

      // Check if Vite is serving correctly
      const hasReact = html.includes('react-refresh') || html.includes('vite');
      console.log('Has Vite/React scripts:', hasReact);

      await guestHome.expectLoaded();
      console.log('Step 3: Guest home loaded successfully');

      // Click "Владелец" link in header (Button with component={Link})
      await page.getByRole('link', { name: 'Владелец' }).click();
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      // Verify owner dashboard is shown
      expect(page.url()).toContain('/owner');
      await expect(page.getByText('Панель управления').first()).toBeVisible();

      // Click "Гость" link in header
      await page.getByRole('link', { name: 'Гость' }).click();
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      // Verify guest home is shown
      expect(page.url()).toMatch(/\/$/);
      await guestHome.expectLoaded();
    });

    test('should show 404 for invalid routes', async ({ page }) => {
      await page.goto('/invalid-route-that-does-not-exist');
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      // Mantine Title renders as div, use getByText
      await expect(page.getByText('404').first()).toBeVisible();
      await expect(page.getByText('Страница не найдена').first()).toBeVisible();

      // Click "На главную" link (Button with component={Link})
      await page.getByRole('link', { name: 'На главную' }).click();
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      expect(page.url()).toMatch(/\/$/);
    });

    test('header navigation links stay visible while switching between guest and owner pages', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // "Гость" link should be visible
      const guestLink = page.getByRole('link', { name: 'Гость' });
      await expect(guestLink).toBeVisible();

      // Navigate to owner
      await page.getByRole('link', { name: 'Владелец' }).click();
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      // "Владелец" link should be visible
      const ownerLink = page.getByRole('link', { name: 'Владелец' });
      await expect(ownerLink).toBeVisible();
    });
  });

  test.describe('Responsive Design', () => {
    test('keeps guest home content visible on a mobile viewport (375px)', async ({ page }) => {
      await setViewport(page, 375, 800);

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();
      await guestHome.expectLoaded();

      await expect(page.getByText('Консультация').first()).toBeVisible();
      await expect(page.getByRole('link', { name: 'Забронировать' }).first()).toBeVisible();
    });

    test('keeps guest home heading visible on a tablet viewport (768px)', async ({ page }) => {
      await setViewport(page, 768, 1024);

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();
      await guestHome.expectLoaded();

      await expect(page.getByText(BOOKING_HEADING).first()).toBeVisible();
    });

    test('keeps guest home heading visible on a desktop viewport (1280px)', async ({ page }) => {
      await setViewport(page, 1280, 800);

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();
      await guestHome.expectLoaded();

      await expect(page.getByText(BOOKING_HEADING).first()).toBeVisible();
    });

    test('keeps owner dashboard tabs reachable on a mobile viewport', async ({ page }) => {
      await setViewport(page, 375, 800);

      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();

      await expect(page.getByRole('tab', { name: 'Типы встреч' })).toBeVisible();
    });
  });

  test.describe('Cookie Consent Banner', () => {
    test('first visit still loads while cookie consent UI is absent', async ({ page, context }) => {
      await context.clearCookies();

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      await guestHome.expectLoaded();
    });

    test.skip('should accept cookies and hide banner', async () => {
      // Cookie consent UI is not implemented yet.
    });
  });

  test.describe('Theme Support', () => {
    test('guest home still loads when dark color scheme is emulated', async ({ page }) => {
      const darkPage = await page.context().newPage();
      await darkPage.emulateMedia({ colorScheme: 'dark' });

      const guestHome = new GuestHomePage(darkPage);
      await guestHome.goto();

      await expect(darkPage.getByText(BOOKING_HEADING).first()).toBeVisible();
    });

    test('guest home still loads when light color scheme is emulated', async ({ page }) => {
      const lightPage = await page.context().newPage();
      await lightPage.emulateMedia({ colorScheme: 'light' });

      const guestHome = new GuestHomePage(lightPage);
      await guestHome.goto();

      await expect(lightPage.getByText(BOOKING_HEADING).first()).toBeVisible();
    });
  });

  test.describe('Internationalization (i18n)', () => {
    test('should display content in Russian by default', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // Verify Russian text is displayed
      await expect(page.getByRole('heading', { name: 'Забронировать встречу' })).toBeVisible();
      await expect(page.getByText('Выберите тип встречи и удобное для вас время')).toBeVisible();
    });

    // Language switching not yet implemented per SPEC.MD
    test.skip('should switch to English when language is changed', async () => {
      // Future implementation
    });

    test.skip('should persist language preference across sessions', async () => {
      // Future implementation
    });
  });

  test.describe('Date and Time Formatting', () => {
    test('guest home uses Russian locale copy by default', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      await expect(page.getByText('Выберите тип встречи и удобное для вас время')).toBeVisible();
    });

    test('owner schedule tab is reachable for time-format coverage', async ({ page }) => {
      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.switchTab('slots');

      await expect(page.getByRole('tab', { name: 'Расписание' })).toBeVisible();
    });
  });
});
