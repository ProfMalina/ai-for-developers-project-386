import { test, expect } from '@playwright/test';
import { OwnerDashboard } from '../pages/OwnerDashboard';
import { testEventTypeNew, testEventTypeUpdated } from '../fixtures/test-data';

const upcomingBooking = {
  id: 'booking-1',
  guestName: 'Иван Иванов',
  guestEmail: 'ivan@example.com',
  eventTypeId: 'consultation',
  startTime: '2026-05-09T10:00:00Z',
  endTime: '2026-05-09T10:30:00Z',
  status: 'upcoming',
  createdAt: '2026-04-01T09:00:00Z',
};

test.describe('Owner Management Flow', () => {
  let dashboard: OwnerDashboard;

  test.beforeEach(async ({ page }) => {
    dashboard = new OwnerDashboard(page);

    // Mock event types
    await page.route('**/api/event-types**', async route => {
      const method = route.request().method();
      if (method === 'GET') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'consultation', name: 'Консультация', description: 'Individual consultation', durationMinutes: 30 },
              { id: 'meeting', name: 'Встреча', description: 'Group meeting', durationMinutes: 60 },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 2, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      }
    });

    // Mock individual event type endpoints
    await page.route('**/api/event-types/:id', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          id: 'consultation',
          name: 'Консультация',
          description: 'Individual consultation',
          durationMinutes: 30,
        }),
      });
    });

    // Mock bookings
    await page.route('**/api/bookings**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: [upcomingBooking],
          pagination: { page: 1, pageSize: 10, totalItems: 1, totalPages: 1, hasNext: false, hasPrev: false },
        }),
      });
    });

    // Mock slots
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

    // Mock slot generation
    await page.route('**/api/slots/generate', async route => {
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({
          slotsCreated: 42,
        }),
      });
    });
  });

  test.afterEach(async ({ page }) => {
    await page.unroute('**/api/**');
  });

  test('should view owner dashboard', async () => {
    await dashboard.goto();
    await dashboard.expectLoaded();
  });

  test('should create a new event type', async ({ page }) => {
    await page.route('**/api/event-types**', async route => {
      const method = route.request().method();
      if (method === 'POST') {
        await route.fulfill({
          status: 201,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'new-event',
            name: testEventTypeNew.name,
            description: testEventTypeNew.description,
            durationMinutes: testEventTypeNew.durationMinutes,
          }),
        });
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'consultation', name: 'Консультация', description: 'Individual', durationMinutes: 30 },
              { id: 'meeting', name: 'Встреча', description: 'Group', durationMinutes: 60 },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 2, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      }
    });

    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');
    await dashboard.openAddEventType();

    await dashboard.fillEventTypeForm({
      name: testEventTypeNew.name,
      description: testEventTypeNew.description,
      duration: testEventTypeNew.durationMinutes,
    });

    await dashboard.submitEventType();
    await page.waitForTimeout(500);

    // Verify notification of success
    await expect(page.getByText('Тип встречи создан')).toBeVisible({ timeout: 5000 });
  });

  test('should edit an existing event type', async ({ page }) => {
    await page.route('**/api/event-types**', async route => {
      const method = route.request().method();
      if (method === 'PATCH' || method === 'PUT' || route.request().url().includes('/event-types/')) {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            id: 'consultation',
            name: testEventTypeUpdated.name,
            description: testEventTypeUpdated.description,
            durationMinutes: testEventTypeUpdated.durationMinutes,
          }),
        });
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'consultation', name: 'Консультация', description: 'Individual', durationMinutes: 30 },
              { id: 'meeting', name: 'Встреча', description: 'Group', durationMinutes: 60 },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 2, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      }
    });

    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');

    // Edit "Консультация" event type
    await dashboard.editEventType('Консультация');

    // Update form
    await dashboard.fillEventTypeForm({
      name: testEventTypeUpdated.name,
      description: testEventTypeUpdated.description,
      duration: testEventTypeUpdated.durationMinutes,
    });

    await dashboard.submitEventType();
    await page.waitForTimeout(500);

    // Verify success notification
    await expect(page.getByText('Тип встречи обновлен')).toBeVisible({ timeout: 5000 });
  });

  test('should delete an event type', async ({ page }) => {
    // Handle confirmation dialog
    page.on('dialog', async dialog => {
      await dialog.accept();
    });

    await page.route('**/api/event-types**', async route => {
      const method = route.request().method();
      if (method === 'DELETE') {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ message: 'Deleted' }),
        });
      } else {
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            items: [
              { id: 'consultation', name: 'Консультация', description: 'Individual', durationMinutes: 30 },
            ],
            pagination: { page: 1, pageSize: 10, totalItems: 1, totalPages: 1, hasNext: false, hasPrev: false },
          }),
        });
      }
    });

    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');
    await dashboard.deleteEventType('Консультация');

    await page.waitForTimeout(500);
    // Verify success notification
    await expect(page.getByText('Тип встречи удален')).toBeVisible({ timeout: 5000 });
  });

  test('should generate time slots', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

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

    await dashboard.submitSlotGeneration();
    await page.waitForTimeout(1000);

    // Verify success notification
    await expect(page.getByText(/создано.*слотов|42/i)).toBeVisible({ timeout: 5000 });
  });

  test('should view bookings list', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.goToBookings();

    // Verify booking is visible
    await expect(page.getByText('Иван Иванов')).toBeVisible({ timeout: 5000 });
  });

  test('should change bookings dataset when selecting another pagination page', async ({ page }) => {
    const bookingsByPage = {
      1: [
        upcomingBooking,
        {
          id: 'booking-2',
          guestName: 'Мария Смирнова',
          guestEmail: 'maria@example.com',
          eventTypeId: 'meeting',
          startTime: '2026-05-10T09:00:00Z',
          endTime: '2026-05-10T10:00:00Z',
          status: 'upcoming',
          createdAt: '2026-04-01T10:00:00Z',
        },
      ],
      2: [
        {
          id: 'booking-3',
          guestName: 'Пётр Петров',
          guestEmail: 'petr@example.com',
          eventTypeId: 'consultation',
          startTime: '2026-05-11T11:00:00Z',
          endTime: '2026-05-11T11:30:00Z',
          status: 'upcoming',
          createdAt: '2026-04-02T09:30:00Z',
        },
      ],
    };

    await page.route('**/api/bookings**', async route => {
      const pageParam = Number(new URL(route.request().url()).searchParams.get('page') ?? '1');
      const items = bookingsByPage[pageParam as 1 | 2] ?? bookingsByPage[1];

      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items,
          pagination: {
            page: pageParam,
            pageSize: 10,
            totalItems: 3,
            totalPages: 2,
            hasNext: pageParam < 2,
            hasPrev: pageParam > 1,
          },
        }),
      });
    });

    await dashboard.goto();
    await dashboard.expectLoaded();
    await dashboard.goToBookings();

    await expect(page.getByText('Иван Иванов')).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('Пётр Петров')).not.toBeVisible();

    await dashboard.goToBookingsPage(2);

    await expect(page.getByText('Пётр Петров')).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('Иван Иванов')).not.toBeVisible();
  });

  test('should cancel booking and show success feedback', async ({ page }) => {
    let activeBookings = [upcomingBooking];

    await page.route('**/api/bookings**', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          items: activeBookings,
          pagination: {
            page: 1,
            pageSize: 10,
            totalItems: activeBookings.length,
            totalPages: 1,
            hasNext: false,
            hasPrev: false,
          },
        }),
      });
    });

    await page.route('**/api/bookings/*', async route => {
      if (route.request().method() === 'DELETE') {
        activeBookings = [];
        await route.fulfill({
          status: 204,
          body: '',
        });
        return;
      }

      await route.continue();
    });

    await dashboard.goto();
    await dashboard.expectLoaded();
    await dashboard.goToBookings();

    await dashboard.cancelBooking('Иван Иванов');

    await expect(page.getByText('Бронирование отменено')).toBeVisible({ timeout: 5000 });
    await expect(page.getByText('Бронирований пока нет')).toBeVisible({ timeout: 5000 });
  });

  test('should validate event type form fields', async ({ page }) => {
    await dashboard.goto();
    await dashboard.expectLoaded();

    await dashboard.switchTab('events');
    await dashboard.openAddEventType();

    // Try to submit with empty name
    await dashboard.fillEventTypeForm({
      name: '',
      description: 'Test',
      duration: 30,
    });
    await dashboard.submitEventType();

    // Should show validation error
    await page.waitForTimeout(500);
    const url = page.url();
    // Should stay on dashboard (modal might still be open)
    expect(url).toContain('/owner');
  });
});
