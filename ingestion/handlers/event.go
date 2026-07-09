package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/adr-p/ingestion/metrics"
	"github.com/adr-p/ingestion/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Event struct {
	Source    string                 `json:"source" binding:"required"`
	EventType string                 `json:"event_type" binding:"required"`
	Timestamp int64                  `json:"timestamp" binding:"required"`
	UserMeta  map[string]interface{} `json:"user_meta"`
	Payload   map[string]interface{} `json:"payload" binding:"required"`
}

func HandleEvent(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event Event
		if err := c.ShouldBindJSON(&event); err != nil {
			metrics.DroppedEvents.Inc()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate HMAC signature
		signature := c.GetHeader("X-Webhook-Signature")
		if !validateSignature(signature, event) {
			metrics.DroppedEvents.Inc()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		// Add metadata
		eventID := uuid.New().String()
		values := map[string]interface{}{
			"id":         eventID,
			"source":     event.Source,
			"event_type": event.EventType,
			"timestamp":  time.Now().Unix(),
			"user_meta":  event.UserMeta,
			"payload":    event.Payload,
		}

		// Push to Redis Stream
		if err := redisClient.XAdd(c.Request.Context(), "abuse:events:raw", values); err != nil {
			metrics.DroppedEvents.Inc()
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "queue unavailable"})
			return
		}

		metrics.ReceivedEvents.Inc()
		metrics.ValidatedEvents.Inc()

		c.JSON(http.StatusOK, gin.H{"event_id": eventID})
	}
}

func validateSignature(signature string, event Event) bool {
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		return true // Disable validation in dev mode
	}

	data := []byte(event.Source + event.EventType)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(data)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}
