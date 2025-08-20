package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"1litw/application"
	"1litw/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUseCase *application.UserUseCase
}

func NewAuthHandler(userUC *application.UserUseCase) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUC,
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

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true) // Clear the cookie
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// NOTE: 需要先修改登入頁面讓他可以支援轉跳回原本的頁面，或是把判斷是否登入丟到前端處理（感覺不太好）
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

	err := h.userUseCase.LinkTelegram(c.Request.Context(), req.Token, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true, "message": "account linked successfully"})
}
