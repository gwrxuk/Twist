package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twist/api-gateway/internal/models"
)

var startTime = time.Now()

// HealthCheck handles the health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	uptime := time.Since(startTime).Seconds()

	response := models.HealthResponse{
		Status:  "ok",
		Version: "1.0.0", // TODO: Get this from a build-time variable
		Uptime:  int64(uptime),
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(response, "Service is healthy"))
}

// Metrics handles the metrics endpoint
func (h *Handler) Metrics(c *gin.Context) {
	// This endpoint is handled by the Prometheus middleware,
	// but we need to define it to register it in the router
	c.String(http.StatusOK, "")
}
