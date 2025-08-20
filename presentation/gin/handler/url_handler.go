package handler

import (
	"net/http"
	"strconv"
	"time"

	"1litw/application"
	"1litw/domain"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	urlUseCase       *application.URLUseCase
	analyticsUseCase *application.AnalyticsUseCase
}

func NewURLHandler(urlUseCase *application.URLUseCase, analyticsUseCase *application.AnalyticsUseCase) *URLHandler {
	return &URLHandler{urlUseCase: urlUseCase, analyticsUseCase: analyticsUseCase}
}

func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req struct {
		OriginalURL string `json:"original_url" binding:"required"`
		CustomPath  string `json:"custom_path"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user") // From JWT middleware

	shortURL, err := h.urlUseCase.CreateShortURL(c.Request.Context(), user.(*domain.User), req.OriginalURL, req.CustomPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shortURL)
}

func (h *URLHandler) Redirect(c *gin.Context) {
	path := c.Param("short_path")
	shortURL, err := h.urlUseCase.GetByPath(c.Request.Context(), path)
	if err != nil || shortURL == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	h.urlUseCase.RecordClick(c.Request.Context(), shortURL.ID, c.Request.UserAgent(), c.ClientIP())
	c.Redirect(http.StatusFound, shortURL.OriginalURL)
}

func (h *URLHandler) GetMyURLs(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	urls, err := h.urlUseCase.ListByUser(c.Request.Context(), user.(*domain.User))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, urls)
}

func (h *URLHandler) DeleteShortURL(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.urlUseCase.DeleteShortURLByID(c.Request.Context(), user.(*domain.User), id)
	if err != nil {
		// Handle specific errors like ErrDeleteNotAllowed, ErrShortURLNotFound
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *URLHandler) GetStats(c *gin.Context) {
	user, _ := c.Get("user") // Can be nil for public stats if we allow it

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	stats, err := h.analyticsUseCase.GetOverviewByID(c.Request.Context(), user.(*domain.User), id, time.Time{}, time.Time{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *URLHandler) GetAllURLs(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if !user.(*domain.User).Permissions.Has(domain.PermViewAnyStats) && !user.(*domain.User).Permissions.Has(domain.PermDeleteAny) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	urls, err := h.urlUseCase.GetAllURLs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, urls)
}
