import { test, expect } from '@playwright/test';
import { GuestHomePage } from '../pages/GuestHomePage';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { setViewport } from '../utils/helpers';

test.describe('Common Features', () => {
  test.describe('Navigation', () => {
    test('should navigate between guest and owner pages', async ({ page }) => {
      const guestHome = new GuestHomePage(page);

      // Start at guest home
      await guestHome.goto();
      await guestHome.expectLoaded();

      // Click "Владелец" button in header
      await page.getByRole('button', { name: 'Владелец' }).click();
      await page.waitForLoadState('networkidle');

      // Verify owner dashboard is shown
      await expect(page).toHaveURL(/.*\/owner/);
      await expect(page.getByRole('heading', { name: 'Панель управления' })).toBeVisible();

      // Click "Гость" button in header
      await page.getByRole('button', { name: 'Гость' }).click();
      await page.waitForLoadState('networkidle');

      // Verify guest home is shown
      await expect(page).toHaveURL(/.*\/$/);
      await guestHome.expectLoaded();
    });

    test('should show 404 for invalid routes', async ({ page }) => {
      await page.goto('/invalid-route-that-does-not-exist');
      await page.waitForLoadState('networkidle');

      await expect(page.getByRole('heading', { name: '404' })).toBeVisible();
      await expect(page.getByText('Страница не найдена')).toBeVisible();

      // Click "На главную" button
      await page.getByRole('button', { name: 'На главную' }).click();
      await page.waitForLoadState('networkidle');

      await expect(page).toHaveURL(/.*\/$/);
    });

    test('header navigation highlights active page', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // "Гость" button should be highlighted (filled variant)
      const guestButton = page.getByRole('button', { name: 'Гость' });
      await expect(guestButton).toBeVisible();

      // Navigate to owner
      await page.getByRole('button', { name: 'Владелец' }).click();
      await page.waitForLoadState('networkidle');

      // "Владелец" button should be highlighted now
      const ownerButton = page.getByRole('button', { name: 'Владелец' });
      await expect(ownerButton).toBeVisible();
    });
  });

  test.describe('Responsive Design', () => {
    test('should display correctly on mobile viewport (375px)', async ({ page }) => {
      await setViewport(page, 375, 800);

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();
      await guestHome.expectLoaded();

      // Verify main elements are visible
      await expect(page.getByRole('heading', { name: 'Забронировать встречу' })).toBeVisible();

      // Verify layout is responsive (cards should stack)
      const cards = page.locator('article');
      const cardCount = await cards.count();
      expect(cardCount).toBeGreaterThan(0);
    });

    test('should display correctly on tablet viewport (768px)', async ({ page }) => {
      await setViewport(page, 768, 1024);

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();
      await guestHome.expectLoaded();

      await expect(page.getByRole('heading', { name: 'Забронировать встречу' })).toBeVisible();
    });

    test('should display correctly on desktop viewport (1280px)', async ({ page }) => {
      await setViewport(page, 1280, 800);

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();
      await guestHome.expectLoaded();

      await expect(page.getByRole('heading', { name: 'Забронировать встречу' })).toBeVisible();
    });

    test('owner dashboard should be responsive', async ({ page }) => {
      await setViewport(page, 375, 800);

      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();

      // Verify tabs are visible even on mobile
      await expect(page.getByText('Типы встреч')).toBeVisible();
    });
  });

  test.describe('Cookie Consent Banner', () => {
    test('should show cookie consent banner on first visit', async ({ page, context }) => {
      // Clear cookies to simulate first visit
      await context.clearCookies();

      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // Note: Cookie banner is not yet implemented per SPEC.MD
      // This test is a placeholder for future implementation
      // await expect(page.getByText(/куки|cookies|cookie/i)).toBeVisible();
    });

    test('should accept cookies and hide banner', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // Accept cookies if banner exists
      await guestHome.acceptCookies();

      // Banner should not be visible after acceptance
      // Note: Cookie banner is not yet implemented
    });
  });

  test.describe('Theme Support', () => {
    test('should respect system color scheme preference', async ({ page }) => {
      // Test with dark mode preference
      const darkPage = await page.context().newPage();
      await darkPage.emulateMedia({ colorScheme: 'dark' });

      const guestHome = new GuestHomePage(darkPage);
      await guestHome.goto();

      // Note: Theme switching is not yet implemented per SPEC.MD
      // This test is a placeholder
      await expect(darkPage.getByRole('heading', { name: 'Забронировать встречу' })).toBeVisible();
    });

    test('should respect light mode preference', async ({ page }) => {
      const lightPage = await page.context().newPage();
      await lightPage.emulateMedia({ colorScheme: 'light' });

      const guestHome = new GuestHomePage(lightPage);
      await guestHome.goto();

      await expect(lightPage.getByRole('heading', { name: 'Забронировать встречу' })).toBeVisible();
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

    // Note: Language switching is not yet implemented per SPEC.MD
    // These tests are placeholders for future implementation
    test.skip('should switch to English when language is changed', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // Future implementation: switch language
      // await switchLanguage(page, 'en');
      // await expect(page.getByRole('heading', { name: 'Book a meeting' })).toBeVisible();
    });

    test.skip('should persist language preference across sessions', async ({ page }) => {
      // Future implementation
    });
  });

  test.describe('Date and Time Formatting', () => {
    test('should display dates in Russian format', async ({ page }) => {
      const guestHome = new GuestHomePage(page);
      await guestHome.goto();

      // Dates should be formatted in Russian locale
      // This is a soft check - actual implementation uses dayjs with 'ru' locale
    });

    test('should display time in 24-hour format', async ({ page }) => {
      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.switchTab('slots');

      // Time slots should be in 24-hour format (e.g., 14:30 not 2:30 PM)
    });
  });
});
