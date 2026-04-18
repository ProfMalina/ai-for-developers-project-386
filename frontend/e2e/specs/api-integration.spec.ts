import { test, expect, type Page } from '@playwright/test';
import { BookingPage } from '../pages/BookingPage';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { testBooking } from '../fixtures/test-data';

const buildPagination = (totalItems: number, pageSize = 100) => ({
  page: 1,
  pageSize,
  totalItems,
  totalPages: Math.max(1, Math.ceil(totalItems / pageSize)),
  hasNext: false,
  hasPrev: false,
});

const createFutureSlot = (eventTypeId = 'test-consultation', daysAhead = 1, hour = 10, minute = 0) => {
  const start = new Date();
  start.setDate(start.getDate() + daysAhead);
  start.setHours(hour, minute, 0, 0);

  const end = new Date(start);
  end.setMinutes(end.getMinutes() + 30);

  return {
    id: `slot-${hour}-${minute}`,
    ownerId: '00000000-0000-0000-0000-000000000001',
    startTime: start.toISOString(),
    endTime: end.toISOString(),
    isAvailable: true,
    eventTypeId,
  };
};

const EMPTY_STATE_TEXT = 'В данный момент нет доступных типов встреч';

const expectBookingPageShell = async (page: Page, eventTypeName: string) => {
  await expect(page.getByText(`Бронирование: ${eventTypeName}`).first()).toBeVisible();
  await expect(page.getByRole('button', { name: /Шаг 1/i }).first()).toBeVisible();
};

const selectBookingDate = async (page: Page, date: Date) => {
  await page.getByRole('button', { name: /Дата встречи/i }).click();
  await page.getByRole('button', { name: new RegExp(`\\b${date.getDate()}\\b`) }).first().click();
  await page.waitForTimeout(500);
};

test.describe('API Integration & Error Handling', () => {
  test.describe('API Response Handling', () => {
    test('submits a booking request and returns to guest home after a mocked 201', async ({ page }) => {
      const slot = createFutureSlot();
      let bookingRequestBody: Record<string, unknown> | null = null;

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
            items: [slot, createFutureSlot('test-consultation', 1, 10, 30)],
            pagination: buildPagination(2),
          }),
        });
      });

      await page.route('**/api/public/bookings', async route => {
        bookingRequestBody = route.request().postDataJSON() as Record<string, unknown>;

        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'booking-1',
            guestName: testBooking.guestName,
            guestEmail: testBooking.guestEmail,
            eventTypeId: 'test-consultation',
            startTime: slot.startTime,
            endTime: slot.endTime,
            createdAt: new Date().toISOString(),
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');
      await expectBookingPageShell(page, 'Консультация');

      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await selectBookingDate(page, tomorrow);

      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);

      await bookingPage.submitBooking();

      expect(bookingRequestBody).toEqual({
        eventTypeId: 'test-consultation',
        slotId: slot.id,
        guestName: testBooking.guestName,
        guestEmail: testBooking.guestEmail,
      });

      await expect(page).toHaveURL(/\/$/);
      await expect(page.getByText('Забронировать встречу').first()).toBeVisible();

      await page.unroute('**/api/**');
    });

    test('shows the API conflict message when booking creation returns 409', async ({ page }) => {
      const slot = createFutureSlot();

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
            items: [slot],
            pagination: buildPagination(1),
          }),
        });
      });

      await page.route('**/api/public/bookings', async route => {
        await route.fulfill({
          status: 409,
          contentType: 'application/json',
          body: JSON.stringify({
            error: 'Conflict',
            message: 'Этот слот уже забронирован',
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');
      await expectBookingPageShell(page, 'Консультация');

      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await selectBookingDate(page, tomorrow);
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);
      await bookingPage.submitBooking();

      await expect(page.getByText('Этот слот уже забронирован')).toBeVisible({ timeout: 5000 });

      await page.unroute('**/api/**');
    });

    test('shows field-level API validation text when booking creation returns 400', async ({ page }) => {
      const slot = createFutureSlot();

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
            items: [slot],
            pagination: buildPagination(1),
          }),
        });
      });

      await page.route('**/api/public/bookings', async route => {
        await route.fulfill({
          status: 400,
          contentType: 'application/json',
          body: JSON.stringify({
            error: 'Bad Request',
            message: 'Ошибка валидации',
            fieldErrors: [{ field: 'guestEmail', message: 'Некорректный email' }],
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');
      await expectBookingPageShell(page, 'Консультация');

      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      await selectBookingDate(page, tomorrow);
      await bookingPage.expectTimeSlotsVisible();
      await page.getByRole('button', { name: /\d{2}:\d{2}/ }).first().click();

      await bookingPage.fillGuestDetails(testBooking.guestName, testBooking.guestEmail);
      await bookingPage.submitBooking();

      await expect(page.getByText('Некорректный email')).toBeVisible({ timeout: 5000 });

      await page.unroute('**/api/**');
    });

    test('redirects back to guest home when event type lookup returns 404', async ({ page }) => {
      await page.route('**/api/public/event-types/non-existent', async route => {
        await route.fulfill({
          status: 404,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Not Found', message: 'Тип встречи не найден' }),
        });
      });

      await page.goto('/book/non-existent');
      await page.waitForLoadState('domcontentloaded');
      await page.waitForTimeout(500);

      await expect(page).toHaveURL(/\/$/);
      await expect(page.getByText('Забронировать встречу').first()).toBeVisible();
      await page.unroute('**/api/**');
    });

    test('falls back to the guest empty state when public event types return 500', async ({ page }) => {
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal Server Error', message: 'Internal server error' }),
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1000);

      await expect(page.getByText(EMPTY_STATE_TEXT)).toBeVisible();

      await page.unroute('**/api/public/event-types**');
    });
  });

  test.describe('Network Failure Handling', () => {
    test('known-broken current behavior: aborted API request leaves a blank app shell', async ({ page }) => {
      await page.route('**/api/**', async route => {
        await route.abort('failed');
      });

      await page.goto('/');
      await page.waitForTimeout(1500);

      // This documents the current broken failure mode rather than healthy recovery.
      await expect(page.locator('#root')).toBeEmpty();

      await page.unroute('**/api/**');
    });

    test('keeps the guest home shell visible while event types are still loading', async ({ page }) => {
      await page.route('**/api/public/event-types**', async route => {
        await new Promise(resolve => setTimeout(resolve, 10000));
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [],
            pagination: buildPagination(0, 10),
          }),
        });
      });

      await page.goto('/', { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(3000);

      await expect(page.getByText('Забронировать встречу').first()).toBeVisible();
      await expect(page.getByText(EMPTY_STATE_TEXT)).not.toBeVisible();
      await page.unroute('**/api/**');
    });
  });

  test.describe('Invalid API Response Handling', () => {
    test('known-broken current behavior: malformed event types JSON leaves a blank app shell', async ({ page }) => {
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: 'invalid json{',
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1500);

      // This documents the current broken failure mode rather than graceful handling.
      await expect(page.locator('#root')).toBeEmpty();
      await page.unroute('**/api/public/event-types**');
    });

    test('keeps the guest home heading visible when event type payload omits fields', async ({ page }) => {
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

      await expect(page.getByText('Забронировать встречу').first()).toBeVisible();
      await page.unroute('**/api/public/event-types**');
    });
  });

  test.describe('Owner API Error Handling', () => {
    test('shows the generic save error when creating an event type fails', async ({ page }) => {
      await page.route('**/api/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [],
            pagination: buildPagination(0, 10),
          }),
        });
      });

      await page.route('**/api/bookings**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [],
            pagination: buildPagination(0, 10),
          }),
        });
      });

      await page.route('**/api/event-types', async route => {
        if (route.request().method() === 'POST') {
          await route.fulfill({
            status: 400,
            contentType: 'application/json',
            body: JSON.stringify({
              error: 'Bad Request',
              message: 'Ошибка валидации',
              fieldErrors: [{ field: 'name', message: 'Название обязательно' }],
            }),
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              items: [],
              pagination: buildPagination(0, 10),
            }),
          });
        }
      });

      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();

      await dashboard.switchTab('events');
      await dashboard.openAddEventType();

      await page.getByRole('textbox', { name: /название/i }).fill('Тестовая консультация');
      await page.getByRole('textbox', { name: /описание/i }).fill('Test');
      await page.getByRole('textbox', { name: /длительность/i }).fill('30');
      await page.getByRole('button', { name: /создать/i }).click();

      await expect(page.getByText('Не удалось сохранить тип встречи')).toBeVisible({ timeout: 5000 });

      await page.unroute('**/api/**');
    });

    test('shows the generic cancellation error when owner booking deletion fails', async ({ page }) => {
      const futureStart = new Date();
      futureStart.setDate(futureStart.getDate() + 1);
      futureStart.setHours(11, 0, 0, 0);
      const futureEnd = new Date(futureStart);
      futureEnd.setMinutes(futureEnd.getMinutes() + 30);

      await page.route('**/api/event-types**', async route => {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [],
            pagination: buildPagination(0, 10),
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
              startTime: futureStart.toISOString(),
              endTime: futureEnd.toISOString(),
              createdAt: new Date().toISOString(),
            }],
            pagination: buildPagination(1, 10),
          }),
        });
      });

      await page.route('**/api/bookings/booking-1', async route => {
        if (route.request().method() === 'DELETE') {
          await route.fulfill({
            status: 500,
            contentType: 'application/json',
            body: JSON.stringify({ error: 'Internal Server Error', message: 'Server error' }),
          });
        } else {
          await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
              id: 'booking-1',
              guestName: 'Test User',
              guestEmail: 'test@example.com',
              eventTypeId: 'test',
              startTime: futureStart.toISOString(),
              endTime: futureEnd.toISOString(),
              createdAt: new Date().toISOString(),
            }),
          });
        }
      });

      const dashboard = new OwnerDashboard(page);
      await dashboard.goto();
      await dashboard.expectLoaded();
      await dashboard.goToBookings();

      await expect(page.getByText('Test User')).toBeVisible({ timeout: 5000 });
      await dashboard.cancelBooking('Test User');
      await expect(page.getByText('Не удалось отменить бронирование')).toBeVisible({ timeout: 5000 });

      await page.unroute('**/api/**');
    });
  });

  test.describe('API Mocking for Isolated Tests', () => {
    test('renders booking page content from mocked event type and slot payloads', async ({ page }) => {
      const slot = createFutureSlot();

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
            items: [slot, createFutureSlot('test-consultation', 1, 10, 30)],
            pagination: buildPagination(2),
          }),
        });
      });

      const bookingPage = new BookingPage(page);
      await bookingPage.goto('test-consultation');

      await expect(page.getByText('Консультация')).toBeVisible();

      await page.unroute('**/api/**');
    });

    test('keeps the guest home shell visible when mocked event types return 500', async ({ page }) => {
      await page.route('**/api/public/event-types**', async route => {
        await route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({
            error: 'Internal Server Error',
            message: 'Internal server error',
          }),
        });
      });

      await page.goto('/');
      await page.waitForTimeout(1500);

      await expect(page.getByText('Забронировать встречу').first()).toBeVisible();
      await expect(page.getByText(EMPTY_STATE_TEXT)).toBeVisible();

      await page.unroute('**/api/public/event-types**');
    });
  });
});
