package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"1litw/application"
	"1litw/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUseCase     *application.UserUseCase
	telegramUseCase *application.TelegramUseCase
}

func NewAuthHandler(userUC *application.UserUseCase, telegramUC *application.TelegramUseCase) *AuthHandler {
	return &AuthHandler{
		userUseCase:     userUC,
		telegramUseCase: telegramUC,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, application.ErrUserExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Println("failed to register user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.userUseCase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.SetCookie("jwt", token, int(time.Hour*24), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"username": user.Username, "permissions": user.Permissions})
}

func (h *AuthHandler) LinkTelegram(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized, please log in first"})
		return
	}
	user := userValue.(*domain.User)

	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	err := h.telegramUseCase.VerifyAndLink(c.Request.Context(), req.Token, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrTokenNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, application.ErrTokenExpired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link account"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "account linked successfully"})
}

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

type UserHandler struct {
	userUseCase *application.UserUseCase
}

func NewUserHandler(userUseCase *application.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) List(c *gin.Context) {
	operator, _ := c.Get("user")

	users, err := h.userUseCase.List(c.Request.Context(), operator.(*domain.User))
	if err != nil {
		log.Println("failed to list users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdatePermissions(c *gin.Context) {
	operator, _ := c.Get("user")

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Permission int `json:"permission"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUseCase.UpdatePermissions(c.Request.Context(), operator.(*domain.User), userID, domain.Permission(req.Permission)); err != nil {
		if errors.Is(err, application.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		log.Println("failed to update user permissions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user permissions"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Delete(c *gin.Context) {
	operator, _ := c.Get("user")

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.userUseCase.Delete(c.Request.Context(), operator.(*domain.User), userID)
	if err != nil {
		if errors.Is(err, application.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		log.Println("failed to delete user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Println(user, exists)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user.(*domain.User))
}
