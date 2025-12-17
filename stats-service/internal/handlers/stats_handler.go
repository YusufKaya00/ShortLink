package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/urlshortener/stats-service/internal/service"
)

type StatsHandler struct {
	service *service.StatsService
}

func NewStatsHandler(service *service.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

// GetURLStats godoc
// @Summary Get stats for a specific URL
// @Tags stats
// @Produce json
// @Param code path string true "Short code"
// @Success 200 {object} models.URLStats
// @Router /api/stats/{code} [get]
func (h *StatsHandler) GetURLStats(c *gin.Context) {
	code := c.Param("code")

	stats, err := h.service.GetURLStats(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetOverallStats godoc
// @Summary Get overall stats
// @Tags stats
// @Produce json
// @Success 200 {object} models.OverallStats
// @Router /api/stats/overall [get]
func (h *StatsHandler) GetOverallStats(c *gin.Context) {
	stats, err := h.service.GetOverallStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecentClicks godoc
// @Summary Get recent clicks
// @Tags stats
// @Produce json
// @Param limit query int false "Limit"
// @Success 200 {array} models.Click
// @Router /api/stats/recent [get]
func (h *StatsHandler) GetRecentClicks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	clicks, err := h.service.GetRecentClicks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clicks)
}

func (h *StatsHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "stats-service"})
}
