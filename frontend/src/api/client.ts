import axios, { AxiosError } from 'axios';
import type {
  EventType,
  Booking,
  TimeSlot,
  PaginatedResponse,
  CreateEventTypeRequest,
  UpdateEventTypeRequest,
  CreateBookingRequest,
  PaginationParams,
  ErrorResponse,
} from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ErrorResponse>) => {
    const status = error.response?.status;
    const data = error.response?.data;

    if (status === 400) {
      throw new ApiValidationError(data?.message || 'Bad request', data);
    } else if (status === 404) {
      throw new ApiNotFoundError(data?.message || 'Not found', data);
    } else if (status === 409) {
      throw new ApiConflictError(data?.message || 'Conflict', data);
    } else {
      throw new ApiError(data?.message || 'Unknown error', data);
    }
  }
);

// Custom error classes
export class ApiError extends Error {
  data?: ErrorResponse;

  constructor(message: string, data?: ErrorResponse) {
    super(message);
    this.name = 'ApiError';
    this.data = data;
  }
}

export class ApiValidationError extends ApiError {
  constructor(message: string, data?: ErrorResponse) {
    super(message, data);
    this.name = 'ApiValidationError';
  }
}

export class ApiNotFoundError extends ApiError {
  constructor(message: string, data?: ErrorResponse) {
    super(message, data);
    this.name = 'ApiNotFoundError';
  }
}

export class ApiConflictError extends ApiError {
  constructor(message: string, data?: ErrorResponse) {
    super(message, data);
    this.name = 'ApiConflictError';
  }
}

// Owner API endpoints

export const ownerApi = {
  // Event Type Management
  async createEventType(data: CreateEventTypeRequest): Promise<EventType> {
    const response = await apiClient.post<EventType>('/api/event-types', data);
    return response.data;
  },

  async getEventTypes(params?: PaginationParams): Promise<PaginatedResponse<EventType>> {
    const response = await apiClient.get<PaginatedResponse<EventType>>('/api/event-types', {
      params,
    });
    return response.data;
  },

  async getEventType(id: string): Promise<EventType> {
    const response = await apiClient.get<EventType>(`/api/event-types/${id}`);
    return response.data;
  },

  async updateEventType(id: string, data: UpdateEventTypeRequest): Promise<EventType> {
    const response = await apiClient.patch<EventType>(`/api/event-types/${id}`, data);
    return response.data;
  },

  async deleteEventType(id: string): Promise<void> {
    await apiClient.delete(`/api/event-types/${id}`);
  },

  // Bookings
  async getAllBookings(params?: PaginationParams & { dateFrom?: string; dateTo?: string }): Promise<PaginatedResponse<Booking>> {
    const response = await apiClient.get<PaginatedResponse<Booking>>('/api/bookings', {
      params,
    });
    return response.data;
  },

  async getBooking(id: string): Promise<Booking> {
    const response = await apiClient.get<Booking>(`/api/bookings/${id}`);
    return response.data;
  },

  async cancelBooking(id: string): Promise<void> {
    await apiClient.delete(`/api/bookings/${id}`);
  },

  // Slots
  async getAllSlots(params?: {
    dateFrom?: string;
    dateTo?: string;
    eventTypeId?: string;
    isAvailable?: boolean;
    page?: number;
    pageSize?: number;
  }): Promise<PaginatedResponse<TimeSlot>> {
    const response = await apiClient.get<PaginatedResponse<TimeSlot>>('/api/slots', {
      params,
    });
    return response.data;
  },

  async generateSlots(
    eventTypeId: string,
    data: {
      workingHoursStart: string;
      workingHoursEnd: string;
      intervalMinutes: number;
      daysOfWeek: number[];
      dateFrom?: string;
      dateTo?: string;
      timezone?: string;
    }
  ): Promise<{ slotsCreated: number; createdSlotIds: string[] }> {
    const response = await apiClient.post<{ slotsCreated: number; createdSlotIds: string[] }>(
      `/api/event-types/${eventTypeId}/slots/generate`,
      data
    );
    return response.data;
  },
};

// Guest API endpoints (Public)

export const guestApi = {
  async getPublicEventTypes(params?: PaginationParams): Promise<PaginatedResponse<EventType>> {
    const response = await apiClient.get<PaginatedResponse<EventType>>('/api/public/event-types', {
      params,
    });
    return response.data;
  },

  async getPublicEventType(id: string): Promise<EventType> {
    const response = await apiClient.get<EventType>(`/api/public/event-types/${id}`);
    return response.data;
  },

  async getAvailableSlots(
    params?: {
      dateFrom?: string;
      dateTo?: string;
      page?: number;
      pageSize?: number;
    }
  ): Promise<PaginatedResponse<TimeSlot>> {
    const response = await apiClient.get<PaginatedResponse<TimeSlot>>(
      '/api/public/slots',
      { params }
    );
    return response.data;
  },

  async createBooking(data: CreateBookingRequest): Promise<Booking> {
    const response = await apiClient.post<Booking>('/api/public/bookings', data);
    return response.data;
  },
};

export default apiClient;
