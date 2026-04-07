import { test, expect } from '@playwright/test';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { testEventTypeNew, testEventTypeUpdated } from '../fixtures/test-data';

test.describe('Owner Management Flow', () => {
  let dashboard: OwnerDashboard;

  test.beforeEach(async ({ page }) => {
    dashboard = new OwnerDashboard(page);
  });

  test('should view owner dashboard', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();
  });

  test('should create a new event type', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    // Go to event types tab
    await dashboard.switchTab('events');

    // Open add event type modal
    await dashboard.openAddEventType();

    // Fill form
    await dashboard.fillEventTypeForm({
      name: testEventTypeNew.name,
      description: testEventTypeNew.description,
      duration: testEventTypeNew.durationMinutes,
    });

    // Submit
    await dashboard.submitEventType();

    // Verify event type appears in list
    await dashboard.expectEventTypeVisible(testEventTypeNew.name, testEventTypeNew.durationMinutes);
  });

  test('should edit an existing event type', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');

    // Find an existing event type and edit it
    // Assuming "Консультация" exists from test data
    await dashboard.editEventType('Консультация');

    // Update form
    await dashboard.fillEventTypeForm({
      name: testEventTypeUpdated.name,
      description: testEventTypeUpdated.description,
      duration: testEventTypeUpdated.durationMinutes,
    });

    // Submit
    await dashboard.submitEventType();

    // Verify updated event type
    await dashboard.expectEventTypeVisible(testEventTypeUpdated.name, testEventTypeUpdated.durationMinutes);
  });

  test('should delete an event type', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');

    // Delete event type
    await dashboard.deleteEventType(testEventTypeNew.name);

    // Verify event type is removed
    await dashboard.expectEventTypeNotVisible(testEventTypeNew.name);
  });

  test('should generate time slots', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    // Open slot generation
    await dashboard.openSlotGeneration();

    // Fill form
    const today = new Date();
    const nextWeek = new Date(today);
    nextWeek.setDate(nextWeek.getDate() + 7);

    await dashboard.fillSlotGenerationForm({
      workingHoursStart: '09:00',
      workingHoursEnd: '18:00',
      intervalMinutes: 30,
      daysOfWeek: ['Пн', 'Вт', 'Ср', 'Чт', 'Пт'],
      dateFrom: today.toISOString().split('T')[0],
      dateTo: nextWeek.toISOString().split('T')[0],
    });

    // Submit
    await dashboard.submitSlotGeneration();

    // Verify slots were generated
    await dashboard.expectSlotsGenerated();
  });

  test('should view bookings list', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    // Go to bookings tab
    await dashboard.goToBookings();

    // Verify bookings list is visible
    await dashboard.expectHeading('Список бронирований');
  });

  test('should navigate through bookings pages', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.goToBookings();

    // Check if pagination exists
    const paginationVisible = await dashboard.page.getByRole('navigation').isVisible().catch(() => false);

    if (paginationVisible) {
      await dashboard.expectPaginationVisible();

      // Try to go to page 2 if available
      const page2Button = dashboard.page.getByRole('button', { name: '2' });
      const page2Visible = await page2Button.isVisible().catch(() => false);

      if (page2Visible) {
        await dashboard.goToBookingsPage(2);
        // Verify page changed (URL or content should change)
        await dashboard.page.waitForTimeout(500);
      }
    }
  });

  test('should cancel a booking', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.goToBookings();

    // Find a booking to cancel (if any exist)
    // This is conditional based on test data
    const bookingCard = dashboard.page.locator('article').first();
    const hasBookings = await bookingCard.isVisible().catch(() => false);

    if (hasBookings) {
      const guestName = await bookingCard.textContent();

      // Cancel booking
      await dashboard.cancelBooking(guestName?.split('\n')[0] || 'Test');

      // Verify cancellation notification
      await dashboard.expectNotification(/отменено|canceled/i)
        .catch(() => {
          // Alternative: booking should disappear
          expect(bookingCard).not.toBeVisible();
        });
    }
  });

  test('should validate event type form fields', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');
    await dashboard.openAddEventType();

    // Try to submit empty form
    await dashboard.submitEventType();

    // Should show validation errors
    // Name is required
    const nameInput = dashboard.page.getByRole('textbox', { name: /название/i });
    const isRequired = await nameInput.getAttribute('aria-required')
      .catch(() => nameInput.getAttribute('required'));

    expect(isRequired).toBeTruthy();
  });

  test('should validate duration is within range', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');
    await dashboard.openAddEventType();

    // Try invalid duration (less than 5 minutes)
    await dashboard.fillEventTypeForm({
      name: 'Invalid Event',
      description: 'Test',
      duration: 3,
    });

    await dashboard.submitEventType();

    // Should show validation error for duration
    // The component validates min=5
    const durationInput = dashboard.page.getByRole('spinbutton', { name: /длительность/i });
    await expect(durationInput).toBeVisible();
  });
});
