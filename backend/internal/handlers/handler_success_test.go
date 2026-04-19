package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/db"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type bookingRepoStub struct {
	getByIDFn            func(context.Context, string) (*models.Booking, error)
	listFn               func(context.Context, int, int, string, string, *time.Time, *time.Time) ([]models.Booking, int, error)
	cancelFn             func(context.Context, string) error
	createFn             func(context.Context, *models.Booking) error
	createWithReservedFn func(context.Context, *models.Booking) error
	checkOverlapFn       func(context.Context, time.Time, time.Time) (bool, error)
	deleteFn             func(context.Context, string) error
}

func (s bookingRepoStub) Create(ctx context.Context, booking *models.Booking) error {
	if s.createFn != nil {
		return s.createFn(ctx, booking)
	}
	return nil
}

func (s bookingRepoStub) CreateWithReservedSlot(ctx context.Context, booking *models.Booking) error {
	if s.createWithReservedFn != nil {
		return s.createWithReservedFn(ctx, booking)
	}
	return nil
}

func (s bookingRepoStub) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, errors.New("unexpected GetByID call")
}

func (s bookingRepoStub) List(ctx context.Context, page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) ([]models.Booking, int, error) {
	if s.listFn != nil {
		return s.listFn(ctx, page, pageSize, sortBy, sortOrder, dateFrom, dateTo)
	}
	return nil, 0, errors.New("unexpected List call")
}

func (s bookingRepoStub) CheckOverlap(ctx context.Context, startTime, endTime time.Time) (bool, error) {
	if s.checkOverlapFn != nil {
		return s.checkOverlapFn(ctx, startTime, endTime)
	}
	return false, nil
}

func (s bookingRepoStub) Cancel(ctx context.Context, id string) error {
	if s.cancelFn != nil {
		return s.cancelFn(ctx, id)
	}
	return errors.New("unexpected Cancel call")
}

func (s bookingRepoStub) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type timeSlotRepoStub struct {
	getByIDFn              func(context.Context, string) (*models.TimeSlot, error)
	listFn                 func(context.Context, string, string, int, int, *bool, *time.Time, *time.Time) ([]models.TimeSlot, int, error)
	getAvailableSlotsFn    func(context.Context, string, int, int, *time.Time, *time.Time) ([]models.TimeSlot, int, error)
	createFn               func(context.Context, *models.TimeSlot) error
	deleteAvailableRangeFn func(context.Context, string, string, time.Time, time.Time) error
	markUnavailableFn      func(context.Context, string) error
}

func (s timeSlotRepoStub) Create(ctx context.Context, slot *models.TimeSlot) error {
	if s.createFn != nil {
		return s.createFn(ctx, slot)
	}
	return nil
}

func (s timeSlotRepoStub) GetByID(ctx context.Context, id string) (*models.TimeSlot, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, errors.New("unexpected GetByID call")
}

func (s timeSlotRepoStub) List(ctx context.Context, ownerID, eventTypeID string, page, pageSize int, available *bool, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	if s.listFn != nil {
		return s.listFn(ctx, ownerID, eventTypeID, page, pageSize, available, startTime, endTime)
	}
	return nil, 0, errors.New("unexpected List call")
}

func (s timeSlotRepoStub) GetAvailableSlots(ctx context.Context, ownerID string, page, pageSize int, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
	if s.getAvailableSlotsFn != nil {
		return s.getAvailableSlotsFn(ctx, ownerID, page, pageSize, startTime, endTime)
	}
	return nil, 0, errors.New("unexpected GetAvailableSlots call")
}

func (s timeSlotRepoStub) DeleteAvailableInRange(ctx context.Context, ownerID, eventTypeID string, startTime, endTime time.Time) error {
	if s.deleteAvailableRangeFn != nil {
		return s.deleteAvailableRangeFn(ctx, ownerID, eventTypeID, startTime, endTime)
	}
	return nil
}

func (s timeSlotRepoStub) MarkAsUnavailable(ctx context.Context, slotID string) error {
	if s.markUnavailableFn != nil {
		return s.markUnavailableFn(ctx, slotID)
	}
	return nil
}

type eventTypeRepoStub struct {
	createFn  func(context.Context, *models.EventType) error
	getByIDFn func(context.Context, string) (*models.EventType, error)
	listFn    func(context.Context, string, int, int, string, string) ([]models.EventType, int, error)
	patchFn   func(context.Context, string, models.UpdateEventTypeRequest) (*models.EventType, error)
	deleteFn  func(context.Context, string) error
}

func (s eventTypeRepoStub) Create(ctx context.Context, eventType *models.EventType) error {
	if s.createFn != nil {
		return s.createFn(ctx, eventType)
	}
	return nil
}

func (s eventTypeRepoStub) GetByID(ctx context.Context, id string) (*models.EventType, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, errors.New("unexpected GetByID call")
}

func (s eventTypeRepoStub) List(ctx context.Context, ownerID string, page, pageSize int, sortBy, sortOrder string) ([]models.EventType, int, error) {
	if s.listFn != nil {
		return s.listFn(ctx, ownerID, page, pageSize, sortBy, sortOrder)
	}
	return nil, 0, errors.New("unexpected List call")
}

func (s eventTypeRepoStub) Patch(ctx context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error) {
	if s.patchFn != nil {
		return s.patchFn(ctx, id, req)
	}
	return nil, errors.New("unexpected Patch call")
}

func (s eventTypeRepoStub) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

type ownerRepoStub struct {
	createFn  func(context.Context, *models.Owner) error
	getByIDFn func(context.Context, string) (*models.Owner, error)
}

func (s ownerRepoStub) Create(ctx context.Context, owner *models.Owner) error {
	if s.createFn != nil {
		return s.createFn(ctx, owner)
	}
	return nil
}

func (s ownerRepoStub) GetByID(ctx context.Context, id string) (*models.Owner, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, errors.New("unexpected GetByID call")
}

type slotConfigRepoStub struct {
	createFn func(context.Context, *models.SlotGenerationConfig) error
}

func (s slotConfigRepoStub) Create(ctx context.Context, config *models.SlotGenerationConfig) error {
	if s.createFn != nil {
		return s.createFn(ctx, config)
	}
	return nil
}

func TestEventTypeHandler_Create_UsesDefaultOwnerIDAndReturnsCreated(t *testing.T) {
	defer func(previous string) { db.DefaultOwnerID = previous }(db.DefaultOwnerID)
	db.DefaultOwnerID = "default-owner-id"

	handler := &EventTypeHandler{service: services.NewEventTypeService(eventTypeRepoStub{
		createFn: func(_ context.Context, eventType *models.EventType) error {
			assert.Equal(t, "default-owner-id", eventType.OwnerID)
			eventType.ID = "event-1"
			eventType.IsActive = true
			return nil
		},
	})}

	r := setupRouter()
	r.POST("/api/event-types", handler.Create)

	body := bytes.NewBufferString(`{"name":"Consultation","description":"Initial consult","durationMinutes":45}`)
	req := httptest.NewRequest(http.MethodPost, "/api/event-types", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "event-1")
	assert.Contains(t, w.Body.String(), "default-owner-id")
}

func TestBookingHandler_Cancel_MapsAlreadyCancelledToBadRequest(t *testing.T) {
	handler := &BookingHandler{service: services.NewBookingService(
		bookingRepoStub{
			getByIDFn: func(_ context.Context, id string) (*models.Booking, error) {
				return &models.Booking{ID: id, Status: "cancelled"}, nil //nolint:misspell // persisted booking status value
			},
		},
		timeSlotRepoStub{},
		eventTypeRepoStub{},
	)}

	r := setupRouter()
	r.DELETE("/api/bookings/:id", handler.Cancel)

	req := httptest.NewRequest(http.MethodDelete, "/api/bookings/booking-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "already canceled")
}

func TestBookingHandler_List_ReturnsPaginatedBookings(t *testing.T) {
	handler := &BookingHandler{service: services.NewBookingService(
		bookingRepoStub{
			listFn: func(_ context.Context, page, pageSize int, sortBy, sortOrder string, dateFrom, dateTo *time.Time) ([]models.Booking, int, error) {
				require.Equal(t, 2, page)
				require.Equal(t, 5, pageSize)
				require.Equal(t, "createdAt", sortBy)
				require.Equal(t, "desc", sortOrder)
				return []models.Booking{{ID: "booking-1", GuestName: "Jane Doe"}}, 1, nil
			},
		},
		timeSlotRepoStub{},
		eventTypeRepoStub{},
	)}

	r := setupRouter()
	r.GET("/api/bookings", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/bookings?page=2&pageSize=5&sortBy=createdAt&sortOrder=desc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "booking-1")
	assert.Contains(t, w.Body.String(), "Jane Doe")
	assert.Contains(t, w.Body.String(), "totalItems")
}

func TestTimeSlotHandler_List_MapsMissingEventTypeToNotFound(t *testing.T) {
	defer func(previous string) { db.DefaultOwnerID = previous }(db.DefaultOwnerID)
	db.DefaultOwnerID = "default-owner-id"

	handler := &TimeSlotHandler{service: services.NewTimeSlotService(
		timeSlotRepoStub{},
		slotConfigRepoStub{},
		ownerRepoStub{},
		eventTypeRepoStub{
			getByIDFn: func(_ context.Context, id string) (*models.EventType, error) {
				return nil, errors.New("missing event type")
			},
		},
	)}

	r := setupRouter()
	r.GET("/api/slots", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/api/slots?eventTypeId=missing-event-type", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Event type")
}

func TestPublicEventTypeHandler_GetByID_HidesInactiveEventType(t *testing.T) {
	handler := &PublicEventTypeHandler{
		etService: services.NewEventTypeService(eventTypeRepoStub{
			getByIDFn: func(_ context.Context, id string) (*models.EventType, error) {
				return &models.EventType{ID: id, IsActive: false}, nil
			},
		}),
		slotService: services.NewTimeSlotService(timeSlotRepoStub{}, slotConfigRepoStub{}, ownerRepoStub{}, eventTypeRepoStub{}),
	}

	r := setupRouter()
	r.GET("/api/public/event-types/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/api/public/event-types/inactive-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Event type")
}

func TestOwnerHandler_GetByID_ReturnsOwner(t *testing.T) {
	handler := &OwnerHandler{service: services.NewOwnerService(ownerRepoStub{
		getByIDFn: func(_ context.Context, id string) (*models.Owner, error) {
			return &models.Owner{ID: id, Name: "Alice", Email: "alice@example.com", Timezone: "UTC"}, nil
		},
	})}

	r := setupRouter()
	r.GET("/api/owners/:id", handler.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/api/owners/owner-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "owner-1")
	assert.Contains(t, w.Body.String(), "alice@example.com")
}

func TestPublicEventTypeHandler_GetSlots_ReturnsAvailableSlots(t *testing.T) {
	defer func(previous string) { db.DefaultOwnerID = previous }(db.DefaultOwnerID)
	db.DefaultOwnerID = "public-owner-id"

	future := time.Now().Add(2 * time.Hour).UTC()
	handler := &PublicEventTypeHandler{
		etService: services.NewEventTypeService(eventTypeRepoStub{}),
		slotService: services.NewTimeSlotService(
			timeSlotRepoStub{
				getAvailableSlotsFn: func(_ context.Context, ownerID string, page, pageSize int, startTime, endTime *time.Time) ([]models.TimeSlot, int, error) {
					require.Equal(t, "public-owner-id", ownerID)
					return []models.TimeSlot{{ID: "slot-1", StartTime: future, EndTime: future.Add(30 * time.Minute), IsAvailable: true}}, 1, nil
				},
			},
			slotConfigRepoStub{},
			ownerRepoStub{},
			eventTypeRepoStub{},
		),
	}

	r := setupRouter()
	r.GET("/api/public/slots", handler.GetSlots)

	req := httptest.NewRequest(http.MethodGet, "/api/public/slots?page=1&pageSize=20", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "slot-1")
	assert.Contains(t, w.Body.String(), "totalItems")
}

func TestEventTypeHandler_Update_NotFoundMapsTo404(t *testing.T) {
	handler := &EventTypeHandler{service: services.NewEventTypeService(eventTypeRepoStub{
		patchFn: func(_ context.Context, id string, req models.UpdateEventTypeRequest) (*models.EventType, error) {
			return nil, errors.New("event type not found")
		},
	})}

	r := setupRouter()
	r.PATCH("/api/event-types/:id", handler.Update)

	body := bytes.NewBufferString(`{"name":"Updated title"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/event-types/missing-id", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Event type")
}

func TestOwnerHandler_Create_ReturnsCreatedOwner(t *testing.T) {
	handler := &OwnerHandler{service: services.NewOwnerService(ownerRepoStub{
		createFn: func(_ context.Context, owner *models.Owner) error {
			owner.ID = "owner-created"
			return nil
		},
	})}

	r := setupRouter()
	r.POST("/api/owners", handler.Create)

	payload, err := json.Marshal(map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"timezone": "UTC",
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/owners", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "owner-created")
	assert.Contains(t, w.Body.String(), "alice@example.com")
}
