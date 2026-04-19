package handlers

import (
	"net/http"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthService *services.HealthService
}

func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ready, err := h.healthService.Ready()

	if err != nil || !ready {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "service is unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
