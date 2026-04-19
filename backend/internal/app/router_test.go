package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryModeRouter_HealthAndCoreFlow(t *testing.T) {
	container, err := NewContainer(ContainerConfig{
		DatabaseURL: "postgres://invalid:invalid@127.0.0.1:1/booking_db?sslmode=disable",
	})
	require.NoError(t, err)
	r := NewRouter(container, "test")

	rootReq := httptest.NewRequest(http.MethodGet, "/", nil)
	rootRes := httptest.NewRecorder()
	r.ServeHTTP(rootRes, rootReq)
	require.Equal(t, http.StatusOK, rootRes.Code)

	healthReq := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthRes := httptest.NewRecorder()
	r.ServeHTTP(healthRes, healthReq)
	require.Equal(t, http.StatusOK, healthRes.Code)
	var health map[string]any
	require.NoError(t, json.Unmarshal(healthRes.Body.Bytes(), &health))
	assert.Equal(t, "memory", health["storageMode"])
	assert.Equal(t, true, health["degradedMode"])

	createEventTypeBody := []byte(`{"name":"Consultation","description":"Memory mode","durationMinutes":30}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/event-types", bytes.NewReader(createEventTypeBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	r.ServeHTTP(createRes, createReq)
	require.Equal(t, http.StatusCreated, createRes.Code)
	var eventType models.EventType
	require.NoError(t, json.Unmarshal(createRes.Body.Bytes(), &eventType))
	require.NotEmpty(t, eventType.ID)

	now := time.Now().Add(48 * time.Hour)
	date := now.Format("2006-01-02")
	day := int(now.Weekday())
	generatePayload := []byte(fmt.Sprintf(`{"workingHoursStart":"09:00","workingHoursEnd":"10:00","intervalMinutes":30,"daysOfWeek":[%d],"dateFrom":"%s","dateTo":"%s"}`,
		day, date, date,
	))
	genReq := httptest.NewRequest(http.MethodPost, "/api/event-types/"+eventType.ID+"/slots/generate", bytes.NewReader(generatePayload))
	genReq.Header.Set("Content-Type", "application/json")
	genRes := httptest.NewRecorder()
	r.ServeHTTP(genRes, genReq)
	require.Equal(t, http.StatusCreated, genRes.Code)

	listReq := httptest.NewRequest(http.MethodGet, "/api/public/event-types", nil)
	listRes := httptest.NewRecorder()
	r.ServeHTTP(listRes, listReq)
	require.Equal(t, http.StatusOK, listRes.Code)
	var eventTypeList models.PaginatedResponse[models.EventType]
	require.NoError(t, json.Unmarshal(listRes.Body.Bytes(), &eventTypeList))
	require.Len(t, eventTypeList.Items, 1)

	dateFrom := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Format(time.RFC3339)
	dateTo := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC).Format(time.RFC3339)
	slotReq := httptest.NewRequest(http.MethodGet, "/api/public/slots?dateFrom="+dateFrom+"&dateTo="+dateTo, nil)
	slotRes := httptest.NewRecorder()
	r.ServeHTTP(slotRes, slotReq)
	require.Equal(t, http.StatusOK, slotRes.Code)
	var slotList models.PaginatedResponse[models.TimeSlot]
	require.NoError(t, json.Unmarshal(slotRes.Body.Bytes(), &slotList))
	require.NotEmpty(t, slotList.Items)
	slotID := slotList.Items[0].ID

	bookingPayload := []byte(fmt.Sprintf(`{"eventTypeId":"%s","slotId":"%s","guestName":"Guest","guestEmail":"guest@example.com"}`,
		eventType.ID, slotID,
	))
	bookReq := httptest.NewRequest(http.MethodPost, "/api/public/bookings", bytes.NewReader(bookingPayload))
	bookReq.Header.Set("Content-Type", "application/json")
	bookRes := httptest.NewRecorder()
	r.ServeHTTP(bookRes, bookReq)
	require.Equal(t, http.StatusCreated, bookRes.Code)

	conflictReq := httptest.NewRequest(http.MethodPost, "/api/public/bookings", bytes.NewReader(bookingPayload))
	conflictReq.Header.Set("Content-Type", "application/json")
	conflictRes := httptest.NewRecorder()
	r.ServeHTTP(conflictRes, conflictReq)
	assert.Equal(t, http.StatusConflict, conflictRes.Code)
}
