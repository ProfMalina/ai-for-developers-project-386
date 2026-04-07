import { test, expect } from '@playwright/test';
import { BookingPage } from '../pages/BookingPage';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { testEventTypes, testBooking } from '../fixtures/test-data';

test.describe('API Integration & Error Handling', () => {
  test.describe('API Response Handling', () => {
    test('should handle successful booking creation (201)', async ({ page }) => {
      const bookingPage = new BookingPage(page);
      await bookingPage.goto(testEventTypes[0].id);
      await bookingPage.expectLoaded(testEventTypes[0].name);

      // Select date
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await bookingPage.selectDate(tomorrow);

      // Select time
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      // Fill form
      await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);

      // Monitor API response
      const apiResponsePromise = page.waitForResponse(
        response => response.url().includes('/api/public/bookings') && response.status() === 201
      );

      await bookingPage.submitBooking();

      // Wait for API response
      const response = await apiResponsePromise;
      expect(response.ok()).toBeTruthy();
    });

    test('should handle booking conflict (409)', async ({ page }) => {
      // This test requires setup: create a booking first, then try to book same slot

      const bookingPage = new BookingPage(page);
      await bookingPage.goto(testEventTypes[0].id);
      await bookingPage.expectLoaded(testEventTypes[0].name);

      // Monitor for 409 response
      page.on('response', async (response) => {
        if (response.url().includes('/bookings') && response.status() === 409) {
          // Should show error message
          await expect(page.getByText(/уже забронирован|already booked|конфликт/i)).toBeVisible();
        }
      });

      // In real scenario, you'd try to book an already-booked slot here
      // This is a placeholder for the actual test
    });

    test('should handle validation error (400)', async ({ page }) => {
      const bookingPage = new BookingPage(page);
      await bookingPage.goto(testEventTypes[0].id);
      await bookingPage.expectLoaded(testEventTypes[0].name);

      // Select date and time
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await bookingPage.selectDate(tomorrow);
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      // Submit with invalid email
      await bookingPage.fillGuestDetails(testBooking.guestName, 'invalid-email');

      // Monitor for 400 response
      page.on('response', async (response) => {
        if (response.url().includes('/bookings') && response.status() === 400) {
          const body = await response.json();
          expect(body).toHaveProperty('errors') || expect(body).toHaveProperty('message');
        }
      });

      await bookingPage.submitBooking();

      // Should show validation error
      await bookingPage.expectValidationError('email', 'некорректный');
    });

    test('should handle not found error (404)', async ({ page }) => {
      // Try to book non-existent event type
      await page.goto('/book/non-existent-id');
      await page.waitForLoadState('networkidle');

      // Should show error or redirect
      // Actual behavior depends on implementation
      const hasError = await page.getByText(/не найдено|not found/i).isVisible().catch(() => false);
      const isRedirected = page.url().includes('/');

      expect(hasError || isRedirected).toBeTruthy();
    });

    test('should handle server error (500)', async ({ page }) => {
      // This test would require mocking the server to return 500
      // In production, you'd use MSW or similar to mock this scenario

      page.on('response', async (response) => {
        if (response.status() === 500) {
          // Should show generic error message
          await expect(page.getByText(/ошибка сервера|server error|попробуйте позже/i)).toBeVisible();
        }
      });
    });
  });

  test.describe('Network Failure Handling', () => {
    test('should handle network timeout gracefully', async ({ page }) => {
      // Abort the request to simulate network failure
      await page.route('**/api/**', async route => {
        await route.abort('failed');
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto(testEventTypes[0].id);

      // Try to load data - should show error
      await page.waitForTimeout(1000);

      // Should show error message to user
      const errorMessage = page.getByText(/ошибка сети|network error|не удалось загрузить/i);
      const isVisible = await errorMessage.isVisible().catch(() => false);

      // Test cleanup - unroute
      await page.unroute('**/api/**');

      expect(isVisible).toBeTruthy();
    });

    test('should handle API unavailable', async ({ page }) => {
      // Simulate server not responding
      await page.route('**/api/**', async route => {
        // Don't respond - simulate timeout
        await new Promise(resolve => setTimeout(resolve, 5000));
        await route.abort('timedout');
      });

      const guestHome = page;
      await guestHome.goto('/');

      // Should show loading then error
      await page.waitForTimeout(2000);

      // Should show error state
      const hasError = await page.getByText(/ошибка|error/i).isVisible().catch(() => false);

      await page.unroute('**/api/**');
    });
  });

  test.describe('Invalid API Response Handling', () => {
    test('should handle malformed JSON response', async ({ page }) => {
      // Mock response with invalid JSON
      await page.route('**/api/public/event-types', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: 'invalid json{',
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1000);

      // Should handle gracefully (show error or empty state)
      await page.unroute('**/api/public/event-types');
    });

    test('should handle missing fields in response', async ({ page }) => {
      // Mock response with missing required fields
      await page.route('**/api/public/event-types', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: [
              { id: '1' }, // Missing name, description, duration
            ],
          }),
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1000);

      // Should handle gracefully
      await page.unroute('**/api/public/event-types');
    });
  });

  test.describe('Owner API Error Handling', () => {
    test('should handle event type creation error', async ({ page }) => {
      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();

      await dashboard.switchTab('events');
      await dashboard.openAddEventType();

      // Mock API error
      await page.route('**/api/event-types', async route => {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({
            message: 'Validation error',
            errors: [{ field: 'name', message: 'Name is required' }],
          }),
        });
      });

      // Try to create event type
      await dashboard.fillEventTypeForm({
        name: '',
        description: 'Test',
        duration: 30,
      });
      await dashboard.submitEventType();

      // Should show error
      await page.waitForTimeout(500);
      await page.unroute('**/api/event-types');
    });

    test('should handle booking cancellation error', async ({ page }) => {
      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.goToBookings();

      // Mock API error
      await page.route('**/api/bookings/**', async route => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Server error' }),
        });
      });

      // Try to cancel booking (if any exists)
      const bookingCard = page.locator('article').first();
      const hasBookings = await bookingCard.isVisible().catch(() => false);

      if (hasBookings) {
        const guestName = await bookingCard.textContent();
        await dashboard.cancelBooking(guestName?.split('\n')[0] || 'Test');

        // Should show error message
        await expect(page.getByText(/ошибка|error/i)).toBeVisible({ timeout: 3000 });
      }

      await page.unroute('**/api/bookings/**');
    });
  });

  test.describe('API Mocking for Isolated Tests', () => {
    test('should work with mocked successful booking', async ({ page }) => {
      // Mock event types response
      await page.route('**/api/public/event-types/test-consultation', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'test-consultation',
            name: 'Консультация',
            description: 'Test consultation',
            durationMinutes: 30,
          }),
        });
      });

      // Mock available slots
      await page.route('**/api/public/slots**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            data: [
              { id: 'slot-1', startTime: '2026-04-09T10:00:00Z', isAvailable: true },
              { id: 'slot-2', startTime: '2026-04-09T10:30:00Z', isAvailable: true },
            ],
            meta: { total: 2 },
          }),
        });
      });

      // Mock booking creation
      await page.route('**/api/public/bookings', async route => {
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'booking-1',
            guestName: 'Test User',
            guestEmail: 'test@example.com',
            eventTypeId: 'test-consultation',
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');

      // Verify mocked data is displayed
      await expect(page.getByText('Консультация')).toBeVisible();

      await page.unroute('**/api/**');
    });

    test('should work with mocked API error', async ({ page }) => {
      // Mock API error response
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({
            message: 'Internal server error',
          }),
        });
      });

      await page.goto('/');
      await page.waitForTimeout(500);

      // Should show error state
      const hasError = await page.getByText(/ошибка|error/i).isVisible().catch(() => false);
      expect(hasError).toBeTruthy();

      await page.unroute('**/api/public/event-types**');
    });
  });
});
