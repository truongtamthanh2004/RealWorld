package handlers

import (
	"Netlfy/services"
	"Netlfy/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ProfileHandler struct {
	service services.ProfileService
}

func NewProfileHandler(service services.ProfileService) *ProfileHandler {
	return &ProfileHandler{service}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	username := c.Param("username")
	
	var currentUserID *uint
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if userID, err := utils.ParseToken(tokenString); err == nil {
			currentUserID = &userID
		}
	}

	user, following, err := h.service.GetProfile(username, currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"username":  user.Username,
			"bio":       user.Bio,
			"image":     user.Image,
			"following": following,
		},
	})
}

func (h *ProfileHandler) FollowUser(c *gin.Context) {
	currentUserID := c.MustGet("userID").(uint)
	username := c.Param("username")
	user, following, err := h.service.FollowUser(currentUserID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"username":  user.Username,
			"bio":       user.Bio,
			"image":     user.Image,
			"following": following,
		},
	})
}

func (h *ProfileHandler) UnfollowUser(c *gin.Context) {
	currentUserID := c.MustGet("userID").(uint)
	username := c.Param("username")
	user, following, err := h.service.UnfollowUser(currentUserID, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"username":  user.Username,
			"bio":       user.Bio,
			"image":     user.Image,
			"following": following,
		},
	})
}
