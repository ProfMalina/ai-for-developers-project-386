package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// Test handler validation without actual service calls
func TestPublicBookingHandler_Create_InvalidJSON(t *testing.T) {
	r := setupRouter()
	handler := NewPublicBookingHandler()
	r.POST("/api/public/bookings", handler.Create)

	// Send invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/public/bookings", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublicBookingHandler_Create_MissingRequiredFields(t *testing.T) {
	r := setupRouter()
	handler := NewPublicBookingHandler()
	r.POST("/api/public/bookings", handler.Create)

	// Send request with missing required fields
	payload := map[string]interface{}{
		"guestName": "John",
		// Missing eventTypeId, guestEmail
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/public/bookings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublicBookingHandler_Create_InvalidEmail(t *testing.T) {
	r := setupRouter()
	handler := NewPublicBookingHandler()
	r.POST("/api/public/bookings", handler.Create)

	payload := map[string]interface{}{
		"eventTypeId": "123e4567-e89b-12d3-a456-426614174000",
		"guestName":   "John Doe",
		"guestEmail":  "invalid-email",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/public/bookings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPublicBookingHandler_Create_InvalidEventTypeID(t *testing.T) {
	r := setupRouter()
	handler := NewPublicBookingHandler()
	r.POST("/api/public/bookings", handler.Create)

	payload := map[string]interface{}{
		"eventTypeId": "invalid-uuid",
		"guestName":   "John Doe",
		"guestEmail":  "john@example.com",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/public/bookings", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventTypeHandler_Create_InvalidJSON(t *testing.T) {
	r := setupRouter()
	handler := NewEventTypeHandler()
	r.POST("/api/event-types", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/event-types", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventTypeHandler_Create_MissingFields(t *testing.T) {
	r := setupRouter()
	handler := NewEventTypeHandler()
	r.POST("/api/event-types", handler.Create)

	payload := map[string]interface{}{
		"name": "Meeting",
		// Missing description, durationMinutes
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/event-types", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOwnerHandler_Create_InvalidJSON(t *testing.T) {
	r := setupRouter()
	handler := NewOwnerHandler()
	r.POST("/api/owners", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/api/owners", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOwnerHandler_Create_MissingFields(t *testing.T) {
	r := setupRouter()
	handler := NewOwnerHandler()
	r.POST("/api/owners", handler.Create)

	payload := map[string]interface{}{
		"name": "John",
		// Missing email and timezone
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/owners", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOwnerHandler_Create_InvalidEmail(t *testing.T) {
	r := setupRouter()
	handler := NewOwnerHandler()
	r.POST("/api/owners", handler.Create)

	payload := map[string]interface{}{
		"name":     "John Doe",
		"email":    "invalid-email",
		"timezone": "UTC",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/owners", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("SuccessResponse", func(t *testing.T) {
		r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			SuccessResponse(c, http.StatusOK, map[string]string{"key": "value"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "value", response["key"])
	})

	t.Run("BadRequest", func(t *testing.T) {
		r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			BadRequest(c, "Invalid request")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Invalid request", response["message"])
	})

	t.Run("NotFound", func(t *testing.T) {
		r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			NotFound(c, "Resource")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("NoContent", func(t *testing.T) {
		r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			NoContent(c)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			ErrorResponse(c, http.StatusInternalServerError, "Internal error", "Something went wrong")
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Internal error", response["error"])
		assert.Equal(t, "Something went wrong", response["message"])
	})
}
