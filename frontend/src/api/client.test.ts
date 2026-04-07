import { describe, it, expect } from 'vitest';
import { ownerApi, guestApi, ApiError, ApiValidationError, ApiNotFoundError, ApiConflictError } from './client';

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

  describe('getAllBookings', () => {
    it('should fetch all bookings', async () => {
      const result = await ownerApi.getAllBookings({ page: 1, pageSize: 10 });

      expect(result.items).toHaveLength(1);
      expect(result.items[0].guestName).toBe('Иван Иванов');
    });
  });

  describe('generateSlots', () => {
    it('should generate time slots', async () => {
      const result = await ownerApi.generateSlots({
        workingHoursStart: '09:00',
        workingHoursEnd: '18:00',
        intervalMinutes: 30,
        daysOfWeek: [1, 2, 3, 4, 5],
        dateFrom: '2026-04-08',
        dateTo: '2026-05-08',
        timezone: 'Europe/Moscow',
      });

      expect(result.slotsCreated).toBe(10);
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
  });

  describe('createBooking', () => {
    it('should create a booking', async () => {
      const bookingData = {
        eventTypeId: 'event-type-1',
        startTime: '2026-04-08T10:00:00Z',
        guestName: 'Тест Пользователь',
        guestEmail: 'test@example.com',
      };

      const result = await guestApi.createBooking(bookingData);

      expect(result.guestName).toBe('Иван Иванов');
      expect(result.eventTypeId).toBe('event-type-1');
    });
  });
});
