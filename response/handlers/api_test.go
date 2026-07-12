package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func parseCallbackData(data string) (action, eventID string, ok bool) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func TestParseCallbackData_Confirm(t *testing.T) {
	action, eventID, ok := parseCallbackData("confirm:abc-123-def")
	assert.True(t, ok)
	assert.Equal(t, "confirm", action)
	assert.Equal(t, "abc-123-def", eventID)
}

func TestParseCallbackData_Reject(t *testing.T) {
	action, eventID, ok := parseCallbackData("reject:abc-123-def")
	assert.True(t, ok)
	assert.Equal(t, "reject", action)
	assert.Equal(t, "abc-123-def", eventID)
}

func TestParseCallbackData_EmptyEventID(t *testing.T) {
	_, _, ok := parseCallbackData("confirm:")
	assert.False(t, ok)
}

func TestParseCallbackData_NoColon(t *testing.T) {
	_, _, ok := parseCallbackData("confirmabc123")
	assert.False(t, ok)
}

func TestParseCallbackData_EmptyString(t *testing.T) {
	_, _, ok := parseCallbackData("")
	assert.False(t, ok)
}

func TestParseCallbackData_UnknownAction(t *testing.T) {
	action, eventID, ok := parseCallbackData("unknown:abc-123")
	assert.True(t, ok)
	assert.Equal(t, "unknown", action)
	assert.Equal(t, "abc-123", eventID)
}

func TestTelegramCallback_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/telegram/callback", TelegramCallback(nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/telegram/callback", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTelegramCallback_EmptyBody(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/telegram/callback", TelegramCallback(nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/telegram/callback", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTelegramCallback_EmptyCallbackData(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/telegram/callback", TelegramCallback(nil))

	body := map[string]interface{}{
		"callback_query": map[string]interface{}{
			"data": "",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/telegram/callback", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTelegramCallback_MissingColon(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/telegram/callback", TelegramCallback(nil))

	body := map[string]interface{}{
		"callback_query": map[string]interface{}{
			"data": "confirmabc123",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/telegram/callback", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTelegramCallback_UnknownAction(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/telegram/callback", TelegramCallback(nil))

	body := map[string]interface{}{
		"callback_query": map[string]interface{}{
			"data": "unknown:abc-123",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/telegram/callback", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
