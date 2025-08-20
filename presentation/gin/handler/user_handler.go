package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"1litw/application"
	"1litw/domain"

	"github.com/gin-gonic/gin"
)

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
