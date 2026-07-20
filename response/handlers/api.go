package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/adr-p/response/db"
	"github.com/adr-p/response/metrics"
	"github.com/gin-gonic/gin"
)

func GetRules(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rules, err := db.GetRules(database)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"rules": rules})
	}
}

func CreateRule(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Name          string `json:"name" binding:"required"`
			ConditionJSON string `json:"condition_json" binding:"required"`
			Action        string `json:"action" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.CreateRule(database, request.Name, []byte(request.ConditionJSON), request.Action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		metrics.RulesCreated.Inc()
		c.JSON(http.StatusCreated, gin.H{"message": "rule created"})
	}
}

func GetStats(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fromStr := c.DefaultQuery("from", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
		toStr := c.DefaultQuery("to", time.Now().Format(time.RFC3339))

		from, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from parameter"})
			return
		}

		to, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to parameter"})
			return
		}

		stats, err := db.GetStats(database, from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	}
}

func TelegramCallback(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			CallbackQuery struct {
				Data string `json:"data"`
			} `json:"callback_query"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		data := request.CallbackQuery.Data
		parts := strings.SplitN(data, ":", 2)
		if len(parts) != 2 || parts[1] == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid callback data, expected 'action:event_id'"})
			return
		}

		action := parts[0]
		eventID := parts[1]

		var verdict string
		switch action {
		case "confirm":
			verdict = "CONFIRMED"
		case "reject":
			verdict = "REJECTED"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action, expected 'confirm' or 'reject'"})
			return
		}

		err := db.UpdateEventVerdict(database, eventID, verdict)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		metrics.TelegramCallbacksHandled.Inc()
		c.JSON(http.StatusOK, gin.H{"message": "action applied", "verdict": verdict, "event_id": eventID})
	}
}
