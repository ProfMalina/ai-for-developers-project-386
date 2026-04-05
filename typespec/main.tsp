import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;

/**
 * Calendar Booking API
 *
 * A calendar booking application where:
 * - Owner can create event types and view upcoming meetings
 * - Guests can browse event types, view calendar, and book free time slots
 * - No overlapping bookings allowed (even for different event types)
 */

@service(#{ title: "Calendar Booking API" })
@server("http://localhost:8080", "Development server")
namespace CalendarBooking;

// ============================================================================
// Domain Entities
// ============================================================================

/**
 * Owner - A single predefined profile who manages the calendar
 * Note: No registration/authentication required
 */
model Owner {
  @minLength(1)
  @maxLength(100)
  id: string;

  @minLength(1)
  @maxLength(100)
  name: string;

  @format("email")
  @minLength(1)
  @maxLength(255)
  email: string;

  /**
   * Owner's preferred timezone for displaying times
   * Format: IANA timezone identifier (e.g., "Europe/Moscow", "America/New_York")
   */
  @minLength(1)
  @maxLength(50)
  timezone: string;
}

/**
 * Event Type - A type of event that can be booked
 * Defined by id, name, description, and duration in minutes
 */
model EventType {
  @minLength(1)
  @maxLength(100)
  id: string;

  @minLength(1)
  @maxLength(100)
  name: string;

  @minLength(1)
  @maxLength(500)
  description: string;

  @minValue(5)
  @maxValue(1440)
  durationMinutes: int32;
}

/**
 * Time Slot - A bookable time slot derived from an event type's duration
 * Represents a specific date and time range
 */
model TimeSlot {
  @minLength(1)
  @maxLength(100)
  id: string;

  @minLength(1)
  @maxLength(100)
  eventTypeId: string;

  startTime: utcDateTime;
  endTime: utcDateTime;
  isAvailable: boolean;
}

/**
 * Booking - A reservation made by a guest for a specific time slot
 */
model Booking {
  @minLength(1)
  @maxLength(100)
  id: string;

  @minLength(1)
  @maxLength(100)
  eventTypeId: string;

  startTime: utcDateTime;
  endTime: utcDateTime;

  @minLength(1)
  @maxLength(100)
  guestName: string;

  @format("email")
  @minLength(1)
  @maxLength(255)
  guestEmail: string;

  createdAt: utcDateTime;
}

/**
 * Booking creation request - Guest provides their details and selects a slot
 */
model CreateBookingRequest {
  @minLength(1)
  @maxLength(100)
  eventTypeId: string;

  /**
   * Start time of the booking
   * Must be in the future (cannot book past or current slots)
   * ISO 8601 format with timezone info (e.g., "2026-04-05T15:30:00+03:00")
   */
  startTime: utcDateTime;

  @minLength(1)
  @maxLength(100)
  guestName: string;

  @format("email")
  @minLength(1)
  @maxLength(255)
  guestEmail: string;

  /**
   * Optional timezone for the booking.
   * If not provided, owner's default timezone is used.
   * Format: IANA timezone identifier (e.g., "Europe/Moscow")
   */
  @maxLength(50)
  timezone?: string;
}

/**
 * Event type creation request
 */
model CreateEventTypeRequest {
  /**
   * Optional ID - if not provided, server will generate one
   */
  @maxLength(100)
  id?: string;

  @minLength(1)
  @maxLength(100)
  name: string;

  @minLength(1)
  @maxLength(500)
  description: string;

  @minValue(5)
  @maxValue(1440)
  durationMinutes: int32;
}

/**
 * Event type update request - all fields optional for partial update
 */
model UpdateEventTypeRequest {
  @minLength(1)
  @maxLength(100)
  name?: string;

  @minLength(1)
  @maxLength(500)
  description?: string;

  @minValue(5)
  @maxValue(1440)
  durationMinutes?: int32;
}

/**
 * Pagination metadata for list responses
 */
model PaginationMeta {
  /** Current page number (1-based) */
  page: int32;

  /** Number of items per page */
  pageSize: int32;

  /** Total number of items available */
  totalItems: int32;

  /** Total number of pages */
  totalPages: int32;

  /** Whether there is a next page */
  hasNext: boolean;

  /** Whether there is a previous page */
  hasPrev: boolean;
}

/**
 * Generic paginated response wrapper
 */
model PaginatedResponse<T> {
  /** List of items */
  items: T[];

  /** Pagination metadata */
  pagination: PaginationMeta;
}

/**
 * Common query parameters for pagination and sorting
 */
model PaginationParams {
  /** Page number (1-based, default: 1) */
  @query @minValue(1)
  page?: int32;

  /** Number of items per page (default: 20, max: 100) */
  @query @minValue(1) @maxValue(100)
  pageSize?: int32;

  /** Field to sort by (default: createdAt) */
  @query
  sortBy?: string;

  /** Sort order: asc or desc (default: desc) */
  @query
  sortOrder?: "asc" | "desc";
}

// ============================================================================
// Owner Endpoints
// ============================================================================

/**
 * Owner operations - Manage event types and view bookings
 */
@route("/api")
namespace OwnerApi {

  // --- Event Type Management ---

  /**
   * Create a new event type
   * Owner can define event types with name, description, and duration
   *
   * Validation Rules:
   * - name: Required, 1-100 characters
   * - description: Required, 1-500 characters
   * - durationMinutes: Required, 5-1440 minutes (5 min to 24 hours)
   *
   * @returns 201 Created - Event type successfully created
   * @returns 400 Bad Request - Validation error in request body
   * @returns 409 Conflict - Event type with same ID already exists
   */
  @post
  @route("/event-types")
  @opExample(#{
    parameters: #{
      body: #{
        name: "Workshop",
        description: "Interactive workshop session",
        durationMinutes: 60
      }
    }
  })
  op createEventType(
    @body body: CreateEventTypeRequest
  ): { @statusCode statusCode: 201, @body eventType: EventType } | BadRequestError | ConflictError;

  /**
   * Get all event types with pagination
   * Returns list of all defined event types
   *
   * @returns 200 OK - List of event types with pagination metadata
   */
  @get
  @route("/event-types")
  op getEventTypes(
    ...PaginationParams
  ): { @body body: PaginatedResponse<EventType> };

  /**
   * Get a specific event type by ID
   *
   * @returns 200 OK - Event type details
   * @returns 404 Not Found - Event type with given ID does not exist
   */
  @get
  @route("/event-types/{eventTypeId}")
  @opExample(#{
    parameters: #{ eventTypeId: "workshop-001" }
  })
  op getEventType(
    @path eventTypeId: string
  ): EventType | NotFoundError;

  /**
   * Update an event type
   * Replaces the entire event type resource
   *
   * @returns 200 OK - Updated event type
   * @returns 400 Bad Request - Validation error
   * @returns 404 Not Found - Event type not found
   */
  @put
  @route("/event-types/{eventTypeId}")
  op updateEventType(
    @path eventTypeId: string,
    @body body: CreateEventTypeRequest
  ): EventType | BadRequestError | NotFoundError;

  /**
   * Delete an event type
   *
   * @returns 204 No Content - Event type successfully deleted
   * @returns 404 Not Found - Event type not found
   * @returns 409 Conflict - Cannot delete event type with existing bookings
   */
  @delete
  @route("/event-types/{eventTypeId}")
  op deleteEventType(
    @path eventTypeId: string
  ): void | NotFoundError | ConflictError;

  // --- Booking Management (Owner View) ---

  /**
   * Get all upcoming bookings with pagination and filtering
   * Owner can view all bookings across all event types in a single list
   *
   * Query Parameters:
   * - dateFrom: Filter bookings starting from this date (ISO 8601)
   * - dateTo: Filter bookings ending before this date (ISO 8601)
   * - page: Page number (default: 1)
   * - pageSize: Items per page (default: 20, max: 100)
   * - sortBy: Sort field (default: startTime)
   * - sortOrder: Sort order (asc/desc, default: asc)
   *
   * Note: All times are returned in the owner's configured timezone
   *
   * @returns 200 OK - Paginated list of bookings
   */
  @get
  @route("/bookings")
  op getAllBookings(
    /** Filter bookings by start date (ISO 8601) */
    @query dateFrom?: utcDateTime,
    /** Filter bookings by end date (ISO 8601) */
    @query dateTo?: utcDateTime,
    /** Optional timezone for displaying booking times */
    @query timezone?: string,
    /** Page number (1-based, default: 1) */
    @query @minValue(1) page?: int32,
    /** Items per page (default: 20, max: 100) */
    @query @minValue(1) @maxValue(100) pageSize?: int32,
    /** Sort field (default: startTime) */
    @query sortBy?: string,
    /** Sort order (asc/desc, default: asc) */
    @query sortOrder?: "asc" | "desc"
  ): { @body body: PaginatedResponse<Booking> };

  /**
   * Get a specific booking by ID
   *
   * @returns 200 OK - Booking details
   * @returns 404 Not Found - Booking not found
   */
  @get
  @route("/bookings/{bookingId}")
  @opExample(#{
    parameters: #{ bookingId: "booking-123" }
  })
  op getBooking(
    @path bookingId: string
  ): Booking | NotFoundError;

  /**
   * Cancel a booking (Owner can cancel bookings)
   *
   * @returns 204 No Content - Booking successfully cancelled
   * @returns 404 Not Found - Booking not found
   * @returns 400 Bad Request - Cannot cancel booking (e.g., already started)
   */
  @delete
  @route("/bookings/{bookingId}")
  op cancelBooking(
    @path bookingId: string
  ): void | NotFoundError | BadRequestError;
}

// ============================================================================
// Guest Endpoints (Public API)
// ============================================================================

/**
 * Guest operations - Browse event types and book slots
 * No authentication required
 */
@route("/api/public")
namespace GuestApi {

  /**
   * Get all available event types with pagination
   * Public page showing name, description, and duration for each event type
   *
   * @returns 200 OK - Paginated list of event types
   */
  @get
  @route("/event-types")
  op getPublicEventTypes(
    ...PaginationParams
  ): { @body body: PaginatedResponse<EventType> };

  /**
   * Get a specific event type by ID (public view)
   *
   * @returns 200 OK - Event type details
   * @returns 404 Not Found - Event type not found
   */
  @get
  @route("/event-types/{eventTypeId}")
  @opExample(#{
    parameters: #{ eventTypeId: "consultation-001" },
    returnType: #{
      id: "consultation-001",
      name: "Consultation",
      description: "One-on-one consultation session",
      durationMinutes: 30
    }
  })
  op getPublicEventType(
    @path eventTypeId: string
  ): EventType | NotFoundError;

  /**
   * Get available time slots for an event type with pagination
   * Returns free slots that can be booked
   *
   * Query Parameters:
   * - dateFrom: Filter slots from this date (ISO 8601)
   * - dateTo: Filter slots until this date (ISO 8601)
   * - timezone: Timezone for displaying times (default: owner's timezone)
   * - page: Page number (default: 1)
   * - pageSize: Items per page (default: 20, max: 100)
   *
   * Business Rules:
   * - Only future slots are returned (slots that haven't started yet)
   * - Slots that are currently in progress or have already ended are excluded
   * - Times are returned in the specified or owner's default timezone
   *
   * @returns 200 OK - Paginated list of available time slots
   * @returns 404 Not Found - Event type not found
   */
  @get
  @route("/event-types/{eventTypeId}/slots")
  op getAvailableSlots(
    @path eventTypeId: string,
    /** Filter slots by start date (ISO 8601) */
    @query dateFrom?: utcDateTime,
    /** Filter slots by end date (ISO 8601) */
    @query dateTo?: utcDateTime,
    /**
     * Optional timezone for displaying slot times.
     * If not provided, owner's default timezone is used.
     * Format: IANA timezone identifier (e.g., "Europe/Moscow")
     */
    @query timezone?: string,
    /** Page number (1-based, default: 1) */
    @query @minValue(1) page?: int32,
    /** Items per page (default: 20, max: 100) */
    @query @minValue(1) @maxValue(100) pageSize?: int32
  ): { @body body: PaginatedResponse<TimeSlot> } | NotFoundError;

  /**
   * Create a booking
   * Guest selects a time slot and provides their details
   *
   * Validation Rules:
   * - eventTypeId: Required, must exist
   * - startTime: Required, must be in the future (strictly greater than current time)
   * - guestName: Required, 1-100 characters
   * - guestEmail: Required, valid email format, 1-255 characters
   * - timezone: Optional, valid IANA timezone
   *
   * Business Rules:
   * - No overlapping bookings allowed
   *   - Cannot book a slot that overlaps with an existing booking
   *   - Applies even for different event types
   * - Cannot book a slot that has already started or ended
   *   - Booking startTime must be in the future (strictly greater than current time)
   * - All times are interpreted in the owner's configured timezone
   *
   * @returns 201 Created - Booking successfully created
   * @returns 400 Bad Request - Validation error (invalid time, email, etc.)
   * @returns 404 Not Found - Event type not found
   * @returns 409 Conflict - Slot overlaps with existing booking
   */
  @post
  @route("/bookings")
  op createBooking(
    @body body: CreateBookingRequest
  ): { @statusCode statusCode: 201, @body booking: Booking } | BadRequestError | NotFoundError | ConflictError;
}

// ============================================================================
// Error Responses
// ============================================================================

/**
 * Standard error response with HTTP status code
 */
model ErrorResponse {
  /** Error type identifier */
  error: string;

  /** Human-readable error message */
  message: string;

  /** Additional error details (optional) */
  details?: string;

  /** Validation errors for individual fields */
  fieldErrors?: FieldError[];
}

/**
 * Field-level validation error
 */
model FieldError {
  /** Field name that caused the error */
  field: string;

  /** Error message for this field */
  message: string;
}

/**
 * Not Found error (404)
 * Returned when a requested resource does not exist
 */
model NotFoundError {
  @statusCode
  status: 404;

  error: "NOT_FOUND";
  message: "The requested resource was not found";
  details?: string;
}

/**
 * Bad Request error (400)
 * Returned when the request is malformed or contains invalid data
 */
model BadRequestError {
  @statusCode
  status: 400;

  error: "BAD_REQUEST";
  message: "The request is invalid or malformed";
  details?: string;
  fieldErrors?: FieldError[];
}

/**
 * Conflict error (409)
 * Returned when the request conflicts with existing data
 * Example: Booking overlaps with an existing reservation
 */
model ConflictError {
  @statusCode
  status: 409;

  error: "CONFLICT";
  message: "The request conflicts with existing data";
  details?: string;
}

/**
 * Invalid time error (400)
 * Returned when trying to book a slot in the past
 */
model InvalidTimeError {
  @statusCode
  status: 400;

  error: "INVALID_TIME";
  message: "Cannot book a slot that has already started or ended";
  details?: string;
}

/**
 * Validation error (400)
 * Returned when request validation fails
 */
model ValidationError {
  @statusCode
  status: 400;

  error: "VALIDATION_ERROR";
  message: "Request validation failed";
  fieldErrors: FieldError[];
}

/**
 * Timezone configuration
 * All datetime values should include timezone information (ISO 8601 with offset or Z for UTC)
 */
model TimezoneInfo {
  /**
   * IANA timezone identifier
   * Examples: "Europe/Moscow", "America/New_York", "Asia/Tokyo"
   */
  @minLength(1)
  @maxLength(50)
  timezone: string;

  /**
   * Current UTC offset for the timezone
   * Format: "+HH:MM" or "-HH:MM"
   */
  @pattern("^[+-]\\d{2}:\\d{2}$")
  utcOffset: string;
}
