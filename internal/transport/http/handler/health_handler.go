package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status    string `json:"status"`
	Database  string `json:"database"`
	Timestamp string `json:"timestamp"`
}

// HealthCheck godoc
// @Summary      Health check
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /healthz [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connectivity
	dbStatus := "connected"
	if sqlDB, err := h.db.DB(); err != nil {
		dbStatus = "disconnected"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	response := HealthResponse{
		Status:    "ok",
		Database:  dbStatus,
		Timestamp: c.GetString("request_time"), // This would need middleware to set this
	}

	if dbStatus == "disconnected" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
			"error":    "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Liveness godoc
// @Summary      Liveness probe
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health/liveness [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Readiness godoc
// @Summary      Readiness probe
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health/readiness [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	// Check database connectivity for readiness
	dbStatus := "connected"
	if sqlDB, err := h.db.DB(); err != nil {
		dbStatus = "disconnected"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "disconnected"
	}

	if dbStatus == "disconnected" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "connected",
	})
}
