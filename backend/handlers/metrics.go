package handlers

import (
	"net/http"
	"time"

	"serverpoint/database"
	"serverpoint/models"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct{}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	id := c.Param("id")

	// Default to last 1 hour
	duration := c.DefaultQuery("duration", "1h")
	var since time.Time
	switch duration {
	case "1h":
		since = time.Now().Add(-1 * time.Hour)
	case "6h":
		since = time.Now().Add(-6 * time.Hour)
	case "24h":
		since = time.Now().Add(-24 * time.Hour)
	default:
		since = time.Now().Add(-1 * time.Hour)
	}

	var metrics []models.Metric
	database.DB.Where("server_id = ? AND created_at > ?", id, since).
		Order("created_at asc").
		Find(&metrics)

	c.JSON(http.StatusOK, metrics)
}

func (h *MetricsHandler) GetLatestMetric(c *gin.Context) {
	id := c.Param("id")

	var metric models.Metric
	result := database.DB.Where("server_id = ?", id).
		Order("created_at desc").
		First(&metric)

	if result.Error != nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	c.JSON(http.StatusOK, metric)
}

func (h *MetricsHandler) GetProcesses(c *gin.Context) {
	id := c.Param("id")

	var processes []models.Process
	database.DB.Where("server_id = ?", id).
		Order("cpu desc").
		Find(&processes)

	c.JSON(http.StatusOK, processes)
}
