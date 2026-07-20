package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		_, _ = gin.DefaultWriter.Write([]byte(
			"[GIN] " + time.Now().Format("2006-01-02 15:04:05") + " | " +
				strconv.Itoa(status) + " | " +
				latency.String() + " | " +
				method + " " + path + query + "\n",
		))
	}
}
