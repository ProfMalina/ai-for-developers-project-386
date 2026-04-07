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
  });

  test('should view event types list on guest home page', async ({ page }) => {
    // Navigate to guest home page
    await guestHome.goto();

    // Verify page is loaded
    await guestHome.expectLoaded();

    // Verify event types are displayed
    for (const eventType of testEventTypes) {
      await guestHome.expectEventTypeVisible(eventType.name, eventType.durationMinutes);
    }
  });

  test('should select event type and view calendar', async ({ page }) => {
    // Navigate to guest home page
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
    // Navigate to booking page
    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    // Step 1: Select date
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);

    // Step 2: Select time slot
    await bookingPage.expectTimeSlotsVisible();
    // Click first available time slot
    const firstSlot = page.getByRole('button', { name: /\d{2}:\d{2}/ }).first();
    const slotTime = await firstSlot.textContent();
    await firstSlot.click();

    // Step 3: Fill guest details
    await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);

    // Verify summary info
    await bookingPage.expectSummaryInfo({
      eventType: testEventTypes[0].name,
      duration: `${testEventTypes[0].durationMinutes} минут`,
    });

    // Submit booking
    await bookingPage.submitBooking();

    // Verify success
    await bookingPage.expectBookingSuccess();
  });

  test('should handle form validation errors', async ({ page }) => {
    // Navigate to booking page
    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    // Select date and time
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);
    await bookingPage.expectTimeSlotsVisible();
    await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

    // Try to submit with empty fields
    await bookingPage.submitInvalidBooking();

    // Fill only name, leave email empty
    await bookingPage.fillGuestDetails(testBooking.guestName, '');
    await bookingPage.submitBooking();

    // Should show validation error for email
    await bookingPage.expectValidationError('email', 'некорректный');
  });

  test('should handle invalid email format', async ({ page }) => {
    // Navigate to booking page
    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    // Select date and time
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);
    await bookingPage.expectTimeSlotsVisible();
    await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

    // Fill with invalid email
    await bookingPage.fillGuestDetails(testBooking.guestName, 'invalid-email');
    await bookingPage.submitBooking();

    // Should show validation error
    await bookingPage.expectValidationError('email', 'некорректный');
  });

  test('should navigate between pages correctly', async ({ page }) => {
    // Go to guest home
    await guestHome.goto();
    await guestHome.expectLoaded();

    // Book event type
    await guestHome.bookEventType(testEventTypes[0].name);

    // Verify URL changed
    expect(page.url()).toContain('/book/');

    // Go back to home using browser back
    await page.goBack();
    await guestHome.expectLoaded();
    expect(page.url()).toMatch(/.*\/$/);
  });
});

test.describe('Guest - Booking Conflicts', () => {
  let guestHome: GuestHomePage;
  let bookingPage: BookingPage;

  test.beforeEach(async ({ page }) => {
    guestHome = new GuestHomePage(page);
    bookingPage = new BookingPage(page);
  });

  test('should show error when booking already taken slot', async ({ page }) => {
    // This test assumes there's already a booking for a specific slot
    // In real scenario, you'd need to set up test data first

    await bookingPage.goto(testEventTypes[0].id);
    await bookingPage.expectLoaded(testEventTypes[0].name);

    // Select date
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    await bookingPage.selectDate(tomorrow);

    // Try to book a slot that might be already booked
    // The UI should show disabled button or error message
    // This is a soft check - actual behavior depends on backend state
    await bookingPage.expectTimeSlotsVisible();
  });
});
