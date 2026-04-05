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
  startTime: utcDateTime;
  guestName: string;
  guestEmail: string;
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
   */
  @get
  @route("/bookings")
  op getAllBookings(
    @query dateFrom?: utcDateTime,
    @query dateTo?: utcDateTime
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
   */
  @get
  @route("/event-types/{eventTypeId}/slots")
  op getAvailableSlots(
    @path eventTypeId: string,
    @query dateFrom?: utcDateTime,
    @query dateTo?: utcDateTime
  ): TimeSlot[] | ErrorResponse;

  /**
   * Create a booking
   * Guest selects a time slot and provides their details
   *
   * Business Rule: No overlapping bookings allowed
   * - Cannot book a slot that overlaps with an existing booking
   * - Applies even for different event types
   */
  @post
  @route("/bookings")
  op createBooking(
    @body body: CreateBookingRequest
  ): Booking | ErrorResponse;
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
