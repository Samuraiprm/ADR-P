package handlers

import (
	"net/http"

	"github.com/adr-p/ingestion/redis"
	"github.com/gin-gonic/gin"
)

func HealthCheck(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := redisClient.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"redis":  "disconnected",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"redis":  "connected",
		})
	}
}
