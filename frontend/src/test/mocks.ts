import { http, HttpResponse } from 'msw';
import { setupServer } from 'msw/node';
import type {
  EventType,
  Booking,
  TimeSlot,
  PaginatedResponse,
  SlotGenerationResult,
} from '../types/api';

// Mock data
export const mockEventTypes: EventType[] = [
  {
    id: 'event-type-1',
    name: 'Консультация',
    description: 'Индивидуальная консультация по проекту',
    durationMinutes: 30,
  },
  {
    id: 'event-type-2',
    name: 'Встреча',
    description: 'Командная встреча',
    durationMinutes: 60,
  },
];

export const mockTimeSlots: TimeSlot[] = [
  {
    id: 'slot-1',
    ownerId: 'owner-1',
    startTime: '2026-04-08T10:00:00Z',
    endTime: '2026-04-08T10:30:00Z',
    isAvailable: true,
  },
  {
    id: 'slot-2',
    ownerId: 'owner-1',
    startTime: '2026-04-08T11:00:00Z',
    endTime: '2026-04-08T11:30:00Z',
    isAvailable: true,
  },
];

export const mockBooking: Booking = {
  id: 'booking-1',
  eventTypeId: 'event-type-1',
  startTime: '2026-04-08T10:00:00Z',
  endTime: '2026-04-08T10:30:00Z',
  guestName: 'Иван Иванов',
  guestEmail: 'ivan@example.com',
  createdAt: '2026-04-07T12:00:00Z',
};

export const createMockPaginatedResponse = <T>(
  items: T[],
  page = 1,
  pageSize = 10
): PaginatedResponse<T> => ({
  items,
  pagination: {
    page,
    pageSize,
    totalItems: items.length,
    totalPages: Math.ceil(items.length / pageSize),
    hasNext: false,
    hasPrev: false,
  },
});

// MSW handlers
export const handlers = [
  // Guest API
  http.get('*/api/public/event-types', ({ request }) => {
    const url = new URL(request.url);
    const page = parseInt(url.searchParams.get('page') || '1');
    const pageSize = parseInt(url.searchParams.get('pageSize') || '10');

    return HttpResponse.json(createMockPaginatedResponse(mockEventTypes, page, pageSize));
  }),

  http.get('*/api/public/event-types/:id', ({ params }) => {
    const { id } = params;
    const eventType = mockEventTypes.find((et) => et.id === id);
    if (!eventType) {
      return new HttpResponse('Not Found', { status: 404 });
    }
    return HttpResponse.json(eventType);
  }),

  http.get('*/api/public/slots', () => {
    return HttpResponse.json(createMockPaginatedResponse(mockTimeSlots));
  }),

  http.post('*/api/public/bookings', async ({ request }) => {
    const body = await request.json();
    const data = body as Record<string, unknown>;

    // Validate required fields
    if (!data.eventTypeId || !data.guestName || !data.guestEmail || !data.startTime) {
      return HttpResponse.json(
        {
          error: 'VALIDATION_ERROR',
          message: 'Validation failed',
          fieldErrors: [
            { field: 'guestName', message: 'Name is required' },
          ],
        },
        { status: 400 }
      );
    }

    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(data.guestEmail as string)) {
      return HttpResponse.json(
        {
          error: 'VALIDATION_ERROR',
          message: 'Validation failed',
          fieldErrors: [
            { field: 'guestEmail', message: 'Invalid email format' },
          ],
        },
        { status: 400 }
      );
    }

    // Check for conflict (slot already booked)
    if (data.slotId === 'booked-slot') {
      return HttpResponse.json(
        {
          error: 'CONFLICT',
          message: 'This slot is already booked',
        },
        { status: 409 }
      );
    }

    return HttpResponse.json(mockBooking, { status: 201 });
  }),

  // Owner API
  http.get('*/api/event-types', () => {
    return HttpResponse.json(createMockPaginatedResponse(mockEventTypes));
  }),

  http.post('*/api/event-types', async ({ request }) => {
    const body = await request.json();
    const data = body as Record<string, unknown>;

    if (!data.name || !data.description || !data.durationMinutes) {
      return HttpResponse.json(
        {
          error: 'VALIDATION_ERROR',
          message: 'Validation failed',
          fieldErrors: [
            { field: 'name', message: 'Name is required' },
          ],
        },
        { status: 400 }
      );
    }

    const newEventType: EventType = {
      id: `event-type-${Date.now()}`,
      name: data.name as string,
      description: data.description as string,
      durationMinutes: data.durationMinutes as number,
    };

    return HttpResponse.json(newEventType, { status: 201 });
  }),

  http.get('*/api/bookings', () => {
    return HttpResponse.json(createMockPaginatedResponse([mockBooking]));
  }),

  http.post('*/api/slots/generate', async () => {
    const result: SlotGenerationResult = {
      slotsCreated: 10,
      slotsSkipped: 0,
      dateFrom: '2026-04-08',
      dateTo: '2026-05-08',
      createdSlotIds: ['slot-1', 'slot-2', 'slot-3'],
    };
    return HttpResponse.json(result);
  }),

  http.get('*/api/slots', () => {
    return HttpResponse.json(createMockPaginatedResponse(mockTimeSlots));
  }),

  // Owner API - additional endpoints
  http.get('*/api/event-types/:id', ({ params }) => {
    const { id } = params;
    const eventType = mockEventTypes.find((et) => et.id === id);
    if (!eventType) {
      return new HttpResponse('Not Found', { status: 404 });
    }
    return HttpResponse.json(eventType);
  }),

  http.patch('*/api/event-types/:id', async ({ request }) => {
    const body = await request.json();
    const data = body as Record<string, unknown>;

    return HttpResponse.json({
      id: 'event-type-1',
      name: data.name || 'Консультация',
      description: data.description || 'Индивидуальная консультация по проекту',
      durationMinutes: (data.durationMinutes as number) || 30,
    });
  }),

  http.delete('*/api/event-types/:id', () => {
    return new HttpResponse(null, { status: 204 });
  }),

  http.get('*/api/bookings/:id', () => {
    return HttpResponse.json(mockBooking);
  }),

  http.delete('*/api/bookings/:id', () => {
    return new HttpResponse(null, { status: 204 });
  }),
];

export const server = setupServer(...handlers);

// Setup and teardown helpers
export const setupMockServer = () => {
  beforeAll(() => server.listen());
  afterEach(() => server.resetHandlers());
  afterAll(() => server.close());
};
