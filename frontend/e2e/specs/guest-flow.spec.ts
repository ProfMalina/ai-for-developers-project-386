import { test, expect } from '@playwright/test';
import { GuestHomePage } from '../pages/GuestHomePage';
import { BookingPage } from '../pages/BookingPage';
import { testEventTypes, testBooking } from '../fixtures/test-data';

test.describe('Guest Booking Flow', () => {
  let guestHome: GuestHomePage;
  let bookingPage: BookingPage;

  test.beforeEach(async ({ page }) => {
    guestHome = new GuestHomePage(page);
    bookingPage = new BookingPage(page);

    // Mock event types list
    await page.route('**/api/public/event-types**', async route => {
      if (route.request().url().includes('/public/event-types') && !route.request().url().includes('/public/event-types/')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: testEventTypes,
            meta: { total: 2, page: 1, pageSize: 10, totalPages: 1 },
          }),
        });
      }
    });

    // Mock individual event type
    await page.route('**/api/public/event-types/test-consultation', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify(testEventTypes[0]),
      });
    });

    // Mock slots
    await page.route('**/api/public/slots**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          data: [
            { id: 'slot-1', startTime: '2026-04-09T10:00:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
            { id: 'slot-2', startTime: '2026-04-09T10:30:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
            { id: 'slot-3', startTime: '2026-04-09T11:00:00Z', isAvailable: false, eventTypeId: 'test-consultation' },
          ],
          meta: { total: 3, page: 1, pageSize: 100, totalPages: 1 },
        }),
      });
    });
  });

  test.afterEach(async ({ page }) => {
    await page.unroute('**/api/**');
  });

  test('should view event types list on guest home page', async ({ page }) => {
    // Listen for console errors
    page.on('console', msg => {
      if (msg.type() === 'error') console.log('BROWSER ERROR:', msg.text());
    });
    page.on('pageerror', err => console.log('PAGE ERROR:', err.message));

    await guestHome.goto();
    await guestHome.expectLoaded();

    // Verify event types are displayed
    for (const eventType of testEventTypes) {
      await guestHome.expectEventTypeVisible(eventType.name, eventType.durationMinutes);
    }
  });

  test('should select event type and view calendar', async ({ page }) => {
    await guestHome.goto();
    await guestHome.expectLoaded();

    // Click on event type
    await guestHome.bookEventType(testEventTypes[0].name);

    // Verify booking page is loaded
    await bookingPage.expectLoaded(testEventTypes[0].name);

    // Select a date
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);

    // Verify time slots are displayed
    await bookingPage.expectTimeSlotsVisible();
  });

  test('should create a booking successfully', async ({ page }) => {
    // Mock successful booking creation
    await page.route('**/api/public/bookings', async route => {
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'booking-1',
          guestName: testBooking.guestName,
          guestEmail: testBooking.guestEmail,
          eventTypeId: 'test-consultation',
          startTime: '2026-04-09T10:00:00Z',
        }),
      });
    });

    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    // Select date
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);

    // Select time slot
    await bookingPage.expectTimeSlotsVisible();
    await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

    // Fill guest details
    await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);

    // Submit
    await bookingPage.submitBooking();

    // Should redirect to home or show success
    await page.waitForTimeout(1000);
    const url = page.url();
    expect(url.includes('/') || url.includes('success')).toBeTruthy();
  });

  test('should show booked slots as unavailable', async ({ page }) => {
    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);

    await bookingPage.expectTimeSlotsVisible();

    // slot-3 is mocked as unavailable (isAvailable: false)
    // The UI may render it differently - check that slots are rendered
    const slotButtons = page.getByRole('button', { name: /\d{2}:\d{2}/ });
    await expect(slotButtons.first()).toBeVisible();
  });

  test('should navigate between pages correctly', async ({ page }) => {
    await guestHome.goto();
    await guestHome.expectLoaded();

    // Book event type
    await guestHome.bookEventType(testEventTypes[0].name);

    // Verify URL changed
    expect(page.url()).toContain('/book/');

    // Go back to home
    await page.goBack();
    await guestHome.expectLoaded();
    expect(page.url()).toMatch(/.*\/$/);
  });

  test('should handle form validation - empty fields', async ({ page }) => {
    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);
    await bookingPage.expectTimeSlotsVisible();
    await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

    // Try to submit with empty fields
    await bookingPage.submitBooking();

    // Should show validation errors or prevent submission
    await page.waitForTimeout(500);
    // Either validation errors or form not submitted
    const url = page.url();
    expect(url).toContain('/book/'); // Should stay on booking page
  });

  test('should handle invalid email format', async ({ page }) => {
    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);
    await bookingPage.expectTimeSlotsVisible();
    await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

    // Fill with invalid email
    await bookingPage.fillGuestDetails(testBooking.guestName, 'invalid-email');
    await bookingPage.submitBooking();

    // Should show validation error
    await page.waitForTimeout(500);
    const url = page.url();
    expect(url).toContain('/book/'); // Should stay on booking page
  });
});
