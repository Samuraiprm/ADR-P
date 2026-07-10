package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/adr-p/response/db"
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

		// Parse callback data: "confirm:event_id" or "reject:event_id"
		data := request.CallbackQuery.Data
		if len(data) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid callback data"})
			return
		}

		action := data[:7]
		eventID := data[8:]

		var verdict string
		if action == "confirm" {
			verdict = "CONFIRMED"
		} else {
			verdict = "REJECTED"
		}

		err := db.UpdateEventVerdict(database, eventID, verdict)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "action applied"})
	}
}
