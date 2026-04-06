// Domain entities from TypeSpec API contract

export interface Owner {
  id: string;
  name: string;
  email: string;
  timezone: string;
}

export interface EventType {
  id: string;
  name: string;
  description: string;
  durationMinutes: number;
}

export interface TimeSlot {
  id: string;
  ownerId: string;
  startTime: string; // UTC datetime
  endTime: string; // UTC datetime
  isAvailable: boolean;
}

export interface Booking {
  id: string;
  eventTypeId: string;
  startTime: string; // UTC datetime
  endTime: string; // UTC datetime
  guestName: string;
  guestEmail: string;
  createdAt: string; // UTC datetime
}

// Request/Response models

export interface CreateEventTypeRequest {
  id?: string;
  name: string;
  description: string;
  durationMinutes: number;
}

export interface UpdateEventTypeRequest {
  name?: string;
  description?: string;
  durationMinutes?: number;
}

export interface CreateBookingRequest {
  eventTypeId: string;
  slotId?: string;
  startTime: string; // UTC datetime
  guestName: string;
  guestEmail: string;
  timezone?: string;
}

export interface SlotGenerationConfig {
  workingHoursStart: string; // HH:MM format
  workingHoursEnd: string; // HH:MM format
  daysOfWeek?: number[]; // 0=Sunday, 1=Monday, ..., 6=Saturday
  dateFrom?: string; // ISO 8601 date
  dateTo?: string; // ISO 8601 date
  timezone?: string;
}

export interface SlotGenerationResult {
  slotsCreated: number;
  slotsSkipped: number;
  dateFrom: string;
  dateTo: string;
  createdSlotIds?: string[];
}

// Pagination

export interface PaginationMeta {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
  hasNext: boolean;
  hasPrev: boolean;
}

export interface PaginatedResponse<T> {
  items: T[];
  pagination: PaginationMeta;
}

export interface PaginationParams {
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// Error models

export interface FieldError {
  field: string;
  message: string;
}

export interface ErrorResponse {
  error: string;
  message: string;
  details?: string;
  fieldErrors?: FieldError[];
}

// API Error types

export type ApiError =
  | { status: 400; data: ErrorResponse }
  | { status: 404; data: ErrorResponse }
  | { status: 409; data: ErrorResponse };
