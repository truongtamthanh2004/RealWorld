package handlers

import (
	"Netlfy/services"
	"Netlfy/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input struct {
		User struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		} `json:"user"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.service.Register(input.User.Username, input.User.Email, input.User.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"email":    user.Email,
			"token":    token,
			"username": user.Username,
			"bio":      user.Bio,
			"image":    user.Image,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input struct {
		User struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		} `json:"user"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.service.Login(input.User.Email, input.User.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"email":    user.Email,
			"username": user.Username,
			"bio":      user.Bio,
			"image":    user.Image,
			"token":    token,
		},
	})
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.service.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user"})
		return
	}

	token, _ := utils.GenerateToken(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"email":    user.Email,
			"username": user.Username,
			"bio":      user.Bio,
			"image":    user.Image,
			"token":    token,
		},
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	userID := userIDInterface.(uint)

	var input struct {
		User struct {
			Email    *string `json:"email"`
			Username *string `json:"username"`
			Password *string `json:"password"`
			Image    *string `json:"image"`
			Bio      *string `json:"bio"`
		} `json:"user"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, token, err := h.service.UpdateUser(userID, input.User.Email, input.User.Username, input.User.Password, input.User.Image, input.User.Bio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"email":    updatedUser.Email,
			"username": updatedUser.Username,
			"bio":      updatedUser.Bio,
			"image":    updatedUser.Image,
			"token":    token,
		},
	})
}
