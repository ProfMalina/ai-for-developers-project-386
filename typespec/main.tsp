import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi";
import "@typespec/versioning";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenApi;
using TypeSpec.Versioning;

/**
 * Calendar Booking API
 *
 * A calendar booking application where:
 * - Owner can create event types and view upcoming meetings
 * - Guests can browse event types, view calendar, and book free time slots
 * - No overlapping bookings allowed (even for different event types)
 */

@service({
  title: "Calendar Booking API",
  version: "1.0.0",
})
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
  id: string;
  name: string;
  email: string;
  /**
   * Owner's preferred timezone for displaying times
   * Format: IANA timezone identifier (e.g., "Europe/Moscow", "America/New_York")
   */
  timezone: string;
}

/**
 * Event Type - A type of event that can be booked
 * Defined by id, name, description, and duration in minutes
 */
model EventType {
  id: string;
  name: string;
  description: string;
  durationMinutes: int32;
}

/**
 * Time Slot - A bookable time slot derived from an event type's duration
 * Represents a specific date and time range
 */
model TimeSlot {
  id: string;
  eventTypeId: string;
  startTime: utcDateTime;
  endTime: utcDateTime;
  isAvailable: boolean;
}

/**
 * Booking - A reservation made by a guest for a specific time slot
 */
model Booking {
  id: string;
  eventTypeId: string;
  startTime: utcDateTime;
  endTime: utcDateTime;
  guestName: string;
  guestEmail: string;
  createdAt: utcDateTime;
}

/**
 * Booking creation request - Guest provides their details and selects a slot
 */
model CreateBookingRequest {
  eventTypeId: string;
  /**
   * Start time of the booking
   * Must be in the future (cannot book past or current slots)
   * ISO 8601 format with timezone info (e.g., "2026-04-05T15:30:00+03:00")
   */
  startTime: utcDateTime;
  guestName: string;
  guestEmail: string;
  /**
   * Optional timezone for the booking.
   * If not provided, owner's default timezone is used.
   */
  timezone?: string;
}

/**
 * Event type creation request
 */
model CreateEventTypeRequest {
  id?: string;
  name: string;
  description: string;
  durationMinutes: int32;
}

/**
 * Available slots query parameters
 */
model GetAvailableSlotsParams {
  @query eventTypeId?: string;
  @query dateFrom?: utcDateTime;
  @query dateTo?: utcDateTime;
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
   */
  @post
  @route("/event-types")
  op createEventType(
    @body body: CreateEventTypeRequest
  ): EventType | ErrorResponse;

  /**
   * Get all event types
   * Returns list of all defined event types
   */
  @get
  @route("/event-types")
  op getEventTypes(): EventType[] | ErrorResponse;

  /**
   * Get a specific event type by ID
   */
  @get
  @route("/event-types/{eventTypeId}")
  op getEventType(
    @path eventTypeId: string
  ): EventType | ErrorResponse;

  /**
   * Update an event type
   */
  @put
  @route("/event-types/{eventTypeId}")
  op updateEventType(
    @path eventTypeId: string,
    @body body: CreateEventTypeRequest
  ): EventType | ErrorResponse;

  /**
   * Delete an event type
   */
  @delete
  @route("/event-types/{eventTypeId}")
  op deleteEventType(
    @path eventTypeId: string
  ): void | ErrorResponse;

  // --- Booking Management (Owner View) ---

  /**
   * Get all upcoming bookings
   * Owner can view all bookings across all event types in a single list
   *
   * Note: All times are returned in the owner's configured timezone
   */
  @get
  @route("/bookings")
  op getAllBookings(
    @query dateFrom?: utcDateTime,
    @query dateTo?: utcDateTime,
    /** Optional timezone for displaying booking times */
    @query timezone?: string
  ): Booking[] | ErrorResponse;

  /**
   * Get a specific booking by ID
   */
  @get
  @route("/bookings/{bookingId}")
  op getBooking(
    @path bookingId: string
  ): Booking | ErrorResponse;

  /**
   * Cancel a booking (Owner can cancel bookings)
   */
  @delete
  @route("/bookings/{bookingId}")
  op cancelBooking(
    @path bookingId: string
  ): void | ErrorResponse;
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
   * Get all available event types
   * Public page showing name, description, and duration for each event type
   */
  @get
  @route("/event-types")
  op getPublicEventTypes(): EventType[] | ErrorResponse;

  /**
   * Get available time slots for an event type
   * Returns free slots that can be booked
   *
   * Business Rules:
   * - Only future slots are returned (slots that haven't started yet)
   * - Slots that are currently in progress or have already ended are excluded
   * - Times are returned in the owner's configured timezone
   */
  @get
  @route("/event-types/{eventTypeId}/slots")
  op getAvailableSlots(
    @path eventTypeId: string,
    /** Filter slots by date range (ISO 8601 format) */
    @query dateFrom?: utcDateTime,
    /** Filter slots by date range (ISO 8601 format) */
    @query dateTo?: utcDateTime,
    /**
     * Optional timezone for displaying slot times.
     * If not provided, owner's default timezone is used.
     * Format: IANA timezone identifier (e.g., "Europe/Moscow")
     */
    @query timezone?: string
  ): TimeSlot[] | ErrorResponse;

  /**
   * Create a booking
   * Guest selects a time slot and provides their details
   *
   * Business Rules:
   * - No overlapping bookings allowed
   *   - Cannot book a slot that overlaps with an existing booking
   *   - Applies even for different event types
   * - Cannot book a slot that has already started or ended
   *   - Booking startTime must be in the future (strictly greater than current time)
   * - All times are interpreted in the owner's configured timezone
   *
   * Validation Errors:
   * - CONFLICT: Slot overlaps with existing booking
   * - INVALID_TIME: Slot is in the past or has already started
   */
  @post
  @route("/bookings")
  op createBooking(
    @body body: CreateBookingRequest
  ): Booking | ErrorResponse | ConflictError | InvalidTimeError;
}

// ============================================================================
// Error Responses
// ============================================================================

/**
 * Standard error response
 */
model ErrorResponse {
  error: string;
  message: string;
  details?: string;
}

/**
 * Conflict error - Used when trying to book an already occupied slot
 */
model ConflictError {
  error: "CONFLICT";
  message: "The selected time slot is already booked";
  details?: string;
}

/**
 * Invalid time error - Used when trying to book a slot in the past
 */
model InvalidTimeError {
  error: "INVALID_TIME";
  message: "Cannot book a slot that has already started or ended";
  details?: string;
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
  timezone: string;
  /**
   * Current UTC offset for the timezone
   * Format: "+HH:MM" or "-HH:MM"
   */
  utcOffset: string;
}
