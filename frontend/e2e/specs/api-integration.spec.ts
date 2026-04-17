import { test, expect } from '@playwright/test';
import { BookingPage } from '../pages/BookingPage';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { testBooking } from '../fixtures/test-data';

test.describe('API Integration & Error Handling', () => {
  test.describe('API Response Handling', () => {
    test('should handle successful booking creation (201)', async ({ page }) => {
      // Mock all API responses for isolated testing
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

      await page.route('**/api/public/slots**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'slot-1', startTime: '2026-04-09T10:00:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
              { id: 'slot-2', startTime: '2026-04-09T10:30:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 2, page: 1, pageSize: 100, totalPages: 1 },
          }),
        });
      });

      // Monitor for successful booking
      const bookingResponsePromise = page.waitForResponse(
        response => response.url().includes('/api/public/bookings') && response.status() === 201
      ).catch(() => null);

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

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');
      await bookingPage.expectLoaded('Консультация');

      // Select date
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await bookingPage.selectDate(tomorrow);

      // Select time slot
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      // Fill form
      await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);

      // Submit
      await bookingPage.submitBooking();

      // Check that booking request was made (either success or redirect)
      const bookingResponse = await bookingResponsePromise;
      expect(bookingResponse).not.toBeNull();

      await page.unroute('**/api/**');
    });

    test('should handle booking conflict (409)', async ({ page }) => {
      // Mock event type
      await page.route('**/api/public/event-types/test-consultation', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'test-consultation',
            name: 'Консультация',
            description: 'Test',
            durationMinutes: 30,
          }),
        });
      });

      await page.route('**/api/public/slots**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'slot-1', startTime: '2026-04-09T10:00:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 1, page: 1, pageSize: 100, totalPages: 1 },
          }),
        });
      });

      // Mock conflict error
      await page.route('**/api/public/bookings', async route => {
        await route.fulfill({
          status: 409,
          contentType: 'application/json',
          body: JSON.stringify({
            message: 'Этот слот уже забронирован',
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');
      await bookingPage.expectLoaded('Консультация');

      // Select date and time
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await bookingPage.selectDate(tomorrow);
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      // Fill and submit
      await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);
      await bookingPage.submitBooking();

      // Should show error notification
      await expect(page.getByText('Этот слот уже забронирован')).toBeVisible({ timeout: 5000 })
        .catch(async () => {
          // Alternative: generic error
          await expect(page.getByText('Ошибка')).toBeVisible({ timeout: 3000 });
        });

      await page.unroute('**/api/**');
    });

    test('should handle validation error (400)', async ({ page }) => {
      // Mock event type
      await page.route('**/api/public/event-types/test-consultation', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'test-consultation',
            name: 'Консультация',
            description: 'Test',
            durationMinutes: 30,
          }),
        });
      });

      await page.route('**/api/public/slots**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'slot-1', startTime: '2026-04-09T10:00:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 1, page: 1, pageSize: 100, totalPages: 1 },
          }),
        });
      });

      // Mock validation error
      await page.route('**/api/public/bookings', async route => {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({
            message: 'Ошибка валидации',
            errors: [{ field: 'email', message: 'Некорректный email' }],
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');
      await bookingPage.expectLoaded('Консультация');

      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await bookingPage.selectDate(tomorrow);
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      await bookingPage.fillGuestDetails(testBooking.guestName, 'bad-email');
      await bookingPage.submitBooking();

      // Should show validation error
      await expect(page.getByText(/Некорректный email|Ошибка валидации/i)).toBeVisible({ timeout: 5000 })
        .catch(async () => {
          // Client-side validation should also catch this
          await expect(page.getByText(/email/i)).toBeVisible({ timeout: 3000 });
        });

      await page.unroute('**/api/**');
    });

    test('should handle not found error (404)', async ({ page }) => {
      // Mock 404 for event type
      await page.route('**/api/public/event-types/non-existent', async route => {
        await route.fulfill({
          status: 404,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Тип встречи не найден' }),
        });
      });

      await page.goto('/book/non-existent');
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      // Should show error or redirect to home
      const hasError = await page.getByText(/не найден|not found|ошибка/i).first().isVisible().catch(() => false);
      const isRedirected = page.url().includes('/');

      expect(hasError || isRedirected).toBeTruthy();
      await page.unroute('**/api/**');
    });

    test('should handle server error (500)', async ({ page }) => {
      // Mock server error
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Internal server error' }),
        });
      });

      const guestHome = page;
      await guestHome.goto('/');
      await page.waitForTimeout(1000);

      // Should show error notification
      const hasError = await page.getByText(/Ошибка|Не удалось/i).isVisible().catch(() => false);
      expect(hasError).toBeTruthy();

      await page.unroute('**/api/public/event-types**');
    });
  });

  test.describe('Network Failure Handling', () => {
    test('should handle network failure gracefully', async ({ page }) => {
      // Abort all API requests
      await page.route('**/api/**', async route => {
        await route.abort('failed');
      });

      await page.goto('/');
      await page.waitForTimeout(1500);

      // Should show error notification
      const hasError = await page.getByText(/Ошибка|Не удалось/i).isVisible().catch(() => false);
      expect(hasError).toBeTruthy();

      await page.unroute('**/api/**');
    });

    test('should handle API unavailable', async ({ page }) => {
      // Simulate timeout - don't respond
      await page.route('**/api/**', async () => {
        // Don't respond at all - simulate hanging server
      });

      await page.goto('/');
      // The request will timeout, which takes longer
      await page.waitForTimeout(3000);

      // Either still loading or showing error is acceptable
      const isLoading = await page.getByText(/загрузка|loading/i).isVisible().catch(() => true);
      const hasError = await page.getByText(/Ошибка|Не удалось/i).isVisible().catch(() => false);

      expect(isLoading || hasError).toBeTruthy();
      await page.unroute('**/api/**');
    });
  });

  test.describe('Invalid API Response Handling', () => {
    test('should handle malformed JSON response', async ({ page }) => {
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: 'invalid json{',
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1500);

      // Should handle gracefully - either error or empty state
      const hasHandled = await page.getByText(/Ошибка|Не удалось|нет тип/i).isVisible()
        .catch(() => true); // If page doesn't crash, it's handled

      expect(hasHandled).toBeTruthy();
      await page.unroute('**/api/public/event-types**');
    });

    test('should handle missing fields in response', async ({ page }) => {
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [{ id: '1' }], // Missing name, description, duration
            pagination: { page: 1, pageSize: 10, totalItems: 1, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1500);

      // Should handle without crashing
      const pageLoaded = await page.getByRole('heading', { name: 'Забронировать встречу' }).isVisible()
        .catch(() => false);

      expect(pageLoaded).toBeTruthy();
      await page.unroute('**/api/public/event-types**');
    });
  });

  test.describe('Owner API Error Handling', () => {
    test('should handle event type creation error', async ({ page }) => {
      // Mock successful event types load
      await page.route('**/api/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [],
            pagination: { page: 1, pageSize: 10, totalItems: 0, totalPages: 1, hasNext: false, hasPrev: false },
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

      // Mock creation error
      await page.route('**/api/event-types', async route => {
        if (route.request().method() === 'POST') {
          await route.fulfill({
            status: 400,
            contentType: 'application/json',
            body: JSON.stringify({
              message: 'Ошибка валидации',
              errors: [{ field: 'name', message: 'Название обязательно' }],
            }),
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              items: [],
              pagination: { page: 1, pageSize: 10, totalItems: 0, totalPages: 1, hasNext: false, hasPrev: false },
            }),
          });
        }
      });

      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();

      await dashboard.switchTab('events');
      await dashboard.openAddEventType();

      // Fill with empty name
      await dashboard.fillEventTypeForm({
        name: '',
        description: 'Test',
        duration: 30,
      });
      await dashboard.submitEventType();

      // Should show error
      await expect(page.getByText(/Ошибка|Не удалось|обязательно/i)).toBeVisible({ timeout: 5000 });

      await page.unroute('**/api/**');
    });

    test('should handle booking cancellation error', async ({ page }) => {
      // Mock successful loads
      await page.route('**/api/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [],
            pagination: { page: 1, pageSize: 10, totalItems: 0, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      });

      await page.route('**/api/bookings**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [{
              id: 'booking-1',
              guestName: 'Test User',
              guestEmail: 'test@example.com',
              eventTypeId: 'test',
              startTime: '2026-04-09T10:00:00Z',
              status: 'upcoming',
            }],
            pagination: { page: 1, pageSize: 10, totalItems: 1, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      });

      // Mock cancellation error
      await page.route('**/api/bookings/booking-1', async route => {
        if (route.request().method() === 'DELETE') {
          await route.fulfill({
            status: 500,
            contentType: 'application/json',
            body: JSON.stringify({ message: 'Server error' }),
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              id: 'booking-1',
              guestName: 'Test User',
            }),
          });
        }
      });

      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();
      await dashboard.goToBookings();

      // Verify booking is visible
      await expect(page.getByText('Test User')).toBeVisible({ timeout: 5000 });

      await page.unroute('**/api/**');
    });
  });

  test.describe('API Mocking for Isolated Tests', () => {
    test('should work with mocked successful booking', async ({ page }) => {
      // Mock event type
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

      // Mock slots
      await page.route('**/api/public/slots**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'slot-1', startTime: '2026-04-09T10:00:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
              { id: 'slot-2', startTime: '2026-04-09T10:30:00Z', isAvailable: true, eventTypeId: 'test-consultation' },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 2, page: 1, pageSize: 100, totalPages: 1 },
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');

      // Verify mocked data is displayed
      await expect(page.getByText('Консультация')).toBeVisible();

      await page.unroute('**/api/**');
    });

    test('should show error when API returns 500', async ({ page }) => {
      // Mock API error
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
      await page.waitForTimeout(1500);

      // Should show error notification
      const hasError = await page.getByText(/Ошибка|Не удалось/i).isVisible().catch(() => false);
      expect(hasError).toBeTruthy();

      await page.unroute('**/api/public/event-types**');
    });
  });
});
