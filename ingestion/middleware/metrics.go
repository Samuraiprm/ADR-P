package middleware

import (
	"strconv"
	"time"

	"github.com/adr-p/ingestion/metrics"
	"github.com/gin-gonic/gin"
)

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTPRequests.WithLabelValues(c.Request.Method, status).Inc()
		metrics.HTTPDuration.WithLabelValues(c.Request.Method).Observe(duration)
	}
}
