package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/urlshortener/url-service/internal/models"
	"github.com/urlshortener/url-service/internal/service"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{service: service}
}

// CreateURL godoc
// @Summary Create a short URL
// @Tags urls
// @Accept json
// @Produce json
// @Param request body models.CreateURLRequest true "URL request"
// @Success 201 {object} models.URLResponse
// @Router /api/urls [post]
func (h *URLHandler) CreateURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID if authenticated (optional)
	var userID *uuid.UUID
	if id, exists := c.Get("user_id"); exists {
		uid := id.(uuid.UUID)
		userID = &uid
	}

	response, err := h.service.CreateURL(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Redirect godoc
// @Summary Redirect to original URL
// @Tags urls
// @Param code path string true "Short code"
// @Success 302
// @Router /{code} [get]
func (h *URLHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	url, err := h.service.GetURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	// Record click asynchronously
	go h.service.RecordClick(
		code,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
		c.GetHeader("Referer"),
	)

	c.Redirect(http.StatusFound, url.OriginalURL)
}

// GetURL godoc
// @Summary Get URL info
// @Tags urls
// @Produce json
// @Param code path string true "Short code"
// @Success 200 {object} models.URLResponse
// @Router /api/urls/{code} [get]
func (h *URLHandler) GetURL(c *gin.Context) {
	code := c.Param("code")

	url, err := h.service.GetURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           url.ID,
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,
		"click_count":  url.ClickCount,
		"created_at":   url.CreatedAt,
	})
}

// GetUserURLs godoc
// @Summary Get user's URLs
// @Tags urls
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.URLResponse
// @Router /api/urls [get]
func (h *URLHandler) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	urls, total, err := h.service.GetUserURLs(userID.(uuid.UUID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":  urls,
		"total": total,
	})
}

// GetAllURLs godoc
// @Summary Get all URLs (public)
// @Tags urls
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.URLResponse
// @Router /api/urls/all [get]
func (h *URLHandler) GetAllURLs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	urls, total, err := h.service.GetAllURLs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":  urls,
		"total": total,
	})
}

// DeleteURL godoc
// @Summary Delete a URL
// @Tags urls
// @Security BearerAuth
// @Param id path string true "URL ID"
// @Success 204
// @Router /api/urls/{id} [delete]
func (h *URLHandler) DeleteURL(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteURL(id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *URLHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "url-service"})
}
