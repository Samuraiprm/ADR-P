package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/adr-p/ingestion/redis"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	os.Setenv("REDIS_ADDR", "localhost:6379")
	
	rdb := redis.NewClient()
	if !rdb.IsConnected(context.Background()) {
		t.Skip("Redis not available, skipping integration test")
	}
	return rdb
}

func TestHandleEvent_ValidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	router.POST("/api/v1/events", HandleEvent(redisClient))

	event := Event{
		Source:    "telegram",
		EventType: "message",
		Timestamp: 1234567890,
		Payload:   map[string]interface{}{"text": "test"},
	}
	body, _ := json.Marshal(event)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleEvent_InvalidSignature(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	redisClient := setupTestRedis(t)
	defer redisClient.Close()

	os.Setenv("WEBHOOK_SECRET", "test-secret")
	router.POST("/api/v1/events", HandleEvent(redisClient))

	event := Event{
		Source:    "telegram",
		EventType: "message",
		Timestamp: 1234567890,
		Payload:   map[string]interface{}{"text": "test"},
	}
	body, _ := json.Marshal(event)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", "invalid-signature")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
