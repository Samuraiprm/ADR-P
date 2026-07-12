package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/adr-p/ingestion/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(redisClient *redis.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/events", HandleEvent(redisClient))
	router.GET("/healthz", HealthCheck(redisClient))
	return router
}

func TestHandleEvent_ValidEvent(t *testing.T) {
	os.Setenv("WEBHOOK_SECRET", "")
	redisClient := redis.NewClient()
	if !redisClient.IsConnected(t.Context()) {
		t.Skip("Redis not available")
	}
	defer redisClient.Close()

	router := setupTestRouter(redisClient)

	event := map[string]interface{}{
		"source":     "telegram",
		"event_type": "message",
		"timestamp":  1234567890,
		"payload":    map[string]interface{}{"text": "hello"},
	}
	body, _ := json.Marshal(event)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["event_id"])
}

func TestHandleEvent_MissingFields(t *testing.T) {
	os.Setenv("WEBHOOK_SECRET", "")
	redisClient := redis.NewClient()
	if !redisClient.IsConnected(t.Context()) {
		t.Skip("Redis not available")
	}
	defer redisClient.Close()

	router := setupTestRouter(redisClient)

	event := map[string]interface{}{
		"source": "telegram",
	}
	body, _ := json.Marshal(event)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleEvent_InvalidSignature(t *testing.T) {
	os.Setenv("WEBHOOK_SECRET", "test-secret-123")
	defer os.Unsetenv("WEBHOOK_SECRET")

	redisClient := redis.NewClient()
	if !redisClient.IsConnected(t.Context()) {
		t.Skip("Redis not available")
	}
	defer redisClient.Close()

	router := setupTestRouter(redisClient)

	event := map[string]interface{}{
		"source":     "telegram",
		"event_type": "message",
		"timestamp":  1234567890,
		"payload":    map[string]interface{}{"text": "hello"},
	}
	body, _ := json.Marshal(event)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", "wrong-signature")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHealthCheck_Healthy(t *testing.T) {
	redisClient := redis.NewClient()
	if !redisClient.IsConnected(t.Context()) {
		t.Skip("Redis not available")
	}
	defer redisClient.Close()

	router := setupTestRouter(redisClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "healthy", resp["status"])
}

func TestHandleEvent_EmptyBody(t *testing.T) {
	os.Setenv("WEBHOOK_SECRET", "")
	redisClient := redis.NewClient()
	if !redisClient.IsConnected(t.Context()) {
		t.Skip("Redis not available")
	}
	defer redisClient.Close()

	router := setupTestRouter(redisClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleEvent_InvalidJSON(t *testing.T) {
	os.Setenv("WEBHOOK_SECRET", "")
	redisClient := redis.NewClient()
	if !redisClient.IsConnected(t.Context()) {
		t.Skip("Redis not available")
	}
	defer redisClient.Close()

	router := setupTestRouter(redisClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateSignature(t *testing.T) {
	os.Setenv("WEBHOOK_SECRET", "")
	event := Event{Source: "telegram", EventType: "message"}
	require.True(t, validateSignature("", event))

	os.Setenv("WEBHOOK_SECRET", "secret")
	defer os.Unsetenv("WEBHOOK_SECRET")
	require.False(t, validateSignature("wrong", event))
}
