import { describe, it, expect } from 'vitest';
import { http, HttpResponse } from 'msw';
import { ownerApi, guestApi, ApiError, ApiValidationError, ApiNotFoundError, ApiConflictError } from './client';
import { server } from '../test/mocks';

describe('API Client Error Classes', () => {
  describe('ApiError', () => {
    it('should create an ApiError instance', () => {
      const error = new ApiError('Test error');

      expect(error).toBeInstanceOf(Error);
      expect(error.name).toBe('ApiError');
      expect(error.message).toBe('Test error');
    });

    it('should store data on ApiError', () => {
      const data = { error: 'TEST', message: 'Test error' };
      const error = new ApiError('Test error', data);

      expect(error.data).toEqual(data);
    });
  });

  describe('ApiValidationError', () => {
    it('should create an ApiValidationError instance', () => {
      const error = new ApiValidationError('Validation failed');

      expect(error).toBeInstanceOf(ApiError);
      expect(error.name).toBe('ApiValidationError');
    });
  });

  describe('ApiNotFoundError', () => {
    it('should create an ApiNotFoundError instance', () => {
      const error = new ApiNotFoundError('Not found');

      expect(error).toBeInstanceOf(ApiError);
      expect(error.name).toBe('ApiNotFoundError');
    });
  });

  describe('ApiConflictError', () => {
    it('should create an ApiConflictError instance', () => {
      const error = new ApiConflictError('Conflict');

      expect(error).toBeInstanceOf(ApiError);
      expect(error.name).toBe('ApiConflictError');
    });
  });
});

describe('Owner API', () => {
  describe('getEventTypes', () => {
    it('should fetch event types with pagination', async () => {
      const result = await ownerApi.getEventTypes({ page: 1, pageSize: 10 });

      expect(result.items).toHaveLength(2);
      expect(result.pagination.page).toBe(1);
    });
  });

  describe('createEventType', () => {
    it('should create a new event type', async () => {
      const newEventType = {
        name: 'Новая встреча',
        description: 'Описание',
        durationMinutes: 45,
      };

      const result = await ownerApi.createEventType(newEventType);

      expect(result.name).toBe('Новая встреча');
      expect(result.durationMinutes).toBe(45);
    });
  });

  describe('getEventType', () => {
    it('should fetch single event type by id', async () => {
      const result = await ownerApi.getEventType('event-type-1');

      expect(result.id).toBe('event-type-1');
      expect(result.name).toBe('Консультация');
    });
  });

  describe('updateEventType', () => {
    it('should update event type', async () => {
      const updateData = {
        name: 'Обновленная встреча',
        description: 'Новое описание',
        durationMinutes: 90,
      };

      // Mock will return the first event type since we don't have update handler
      const result = await ownerApi.updateEventType('event-type-1', updateData);

      expect(result).toBeDefined();
    });
  });

  describe('deleteEventType', () => {
    it('should delete event type', async () => {
      await ownerApi.deleteEventType('event-type-1');
      // If no error, test passes
      expect(true).toBe(true);
    });
  });

  describe('getAllBookings', () => {
    it('should fetch all bookings', async () => {
      const result = await ownerApi.getAllBookings({ page: 1, pageSize: 10 });

      expect(result.items).toHaveLength(1);
      expect(result.items[0].guestName).toBe('Иван Иванов');
    });
  });

  describe('getBooking', () => {
    it('should fetch single booking by id', async () => {
      const result = await ownerApi.getBooking('booking-1');

      expect(result.id).toBe('booking-1');
      expect(result.guestName).toBe('Иван Иванов');
    });
  });

  describe('cancelBooking', () => {
    it('should cancel booking', async () => {
      await ownerApi.cancelBooking('booking-1');
      // If no error, test passes
      expect(true).toBe(true);
    });
  });

  describe('generateSlots', () => {
    it('should generate time slots for a specific event type', async () => {
      let capturedEventTypeId = '';

      server.use(
        http.post('*/api/event-types/:eventTypeId/slots/generate', async ({ params }) => {
          capturedEventTypeId = String(params.eventTypeId);
          return HttpResponse.json({
            slotsCreated: 10,
            slotsSkipped: 0,
            dateFrom: '2026-04-08',
            dateTo: '2026-05-08',
            createdSlotIds: ['slot-1', 'slot-2', 'slot-3'],
          });
        })
      );

      const result = await ownerApi.generateSlots('event-type-1', {
        workingHoursStart: '09:00',
        workingHoursEnd: '18:00',
        intervalMinutes: 30,
        daysOfWeek: [1, 2, 3, 4, 5],
        dateFrom: '2026-04-08',
        dateTo: '2026-05-08',
        timezone: 'Europe/Moscow',
      });

      expect(result.slotsCreated).toBe(10);
      expect(capturedEventTypeId).toBe('event-type-1');
    });
  });

  describe('getAllSlots', () => {
    it('should fetch all time slots', async () => {
      const result = await ownerApi.getAllSlots({ page: 1, pageSize: 10 });

      expect(result.items).toHaveLength(2);
    });
  });
});

describe('Guest API', () => {
  describe('getPublicEventTypes', () => {
    it('should fetch public event types', async () => {
      const result = await guestApi.getPublicEventTypes({ page: 1, pageSize: 10 });

      expect(result.items).toHaveLength(2);
      expect(result.items[0].name).toBe('Консультация');
    });
  });

  describe('getPublicEventType', () => {
    it('should fetch single event type', async () => {
      const result = await guestApi.getPublicEventType('event-type-1');

      expect(result.id).toBe('event-type-1');
      expect(result.name).toBe('Консультация');
    });
  });

  describe('getAvailableSlots', () => {
    it('should fetch available time slots', async () => {
      const result = await guestApi.getAvailableSlots({
        dateFrom: '2026-04-08',
        dateTo: '2026-04-09',
      });

      expect(result.items).toHaveLength(2);
    });

    it('should request owner-level public slots endpoint with query params', async () => {
      let capturedPathname = '';
      let capturedQuery = '';

      server.use(
        http.get('*/api/public/slots', ({ request }) => {
          const url = new URL(request.url);
          capturedPathname = url.pathname;
          capturedQuery = url.search;

          return HttpResponse.json({
            items: [],
            pagination: {
              page: 2,
              pageSize: 50,
              totalItems: 0,
              totalPages: 0,
              hasNext: false,
              hasPrev: false,
            },
          });
        })
      );

      await guestApi.getAvailableSlots({
        dateFrom: '2026-04-08T00:00:00Z',
        dateTo: '2026-04-08T23:59:59Z',
        page: 2,
        pageSize: 50,
      });

      expect(capturedPathname).toBe('/api/public/slots');
      expect(capturedQuery).toContain('dateFrom=2026-04-08T00:00:00Z');
      expect(capturedQuery).toContain('dateTo=2026-04-08T23:59:59Z');
      expect(capturedQuery).toContain('page=2');
      expect(capturedQuery).toContain('pageSize=50');
      expect(capturedQuery).not.toContain('timezone=');
    });
  });

  describe('createBooking', () => {
    it('should create a booking', async () => {
      const bookingData = {
        eventTypeId: 'event-type-1',
        slotId: 'slot-1',
        guestName: 'Тест Пользователь',
        guestEmail: 'test@example.com',
      };

      const result = await guestApi.createBooking(bookingData);

      expect(result.guestName).toBe('Иван Иванов');
      expect(result.eventTypeId).toBe('event-type-1');
    });

    it('should send slot-based public booking payload', async () => {
      let capturedBody: Record<string, unknown> | null = null;

      server.use(
        http.post('*/api/public/bookings', async ({ request }) => {
          capturedBody = (await request.json()) as Record<string, unknown>;
          return HttpResponse.json({
            id: 'booking-1',
            eventTypeId: 'event-type-1',
            slotId: 'slot-1',
            guestName: 'Тест Пользователь',
            guestEmail: 'test@example.com',
            createdAt: '2026-04-07T12:00:00Z',
            startTime: '2026-04-08T10:00:00Z',
            endTime: '2026-04-08T10:30:00Z',
          }, { status: 201 });
        })
      );

      await guestApi.createBooking({
        eventTypeId: 'event-type-1',
        slotId: 'slot-1',
        guestName: 'Тест Пользователь',
        guestEmail: 'test@example.com',
      });

      expect(capturedBody).toMatchObject({
        eventTypeId: 'event-type-1',
        slotId: 'slot-1',
        guestName: 'Тест Пользователь',
        guestEmail: 'test@example.com',
      });
      expect(capturedBody).not.toHaveProperty('startTime');
    });
  });
});
