package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twist/api-gateway/pkg/metrics"
)

// Metrics middleware records request metrics
func Metrics(metrics *metrics.PrometheusClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Stop timer
		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		// Record metrics
		metrics.RecordRequest(path, method, status)
		metrics.RecordRequestDuration(path, method, duration)
	}
}
