package handlers

import (
	"Netlfy/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	service services.CommentService
}

func NewCommentHandler(service services.CommentService) *CommentHandler {
	return &CommentHandler{service}
}

func (h *CommentHandler) AddComment(c *gin.Context) {
	slug := c.Param("slug")

	var req struct {
		Comment struct {
			Body string `json:"body" binding:"required"`
		} `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetUint("userID")

	comment, err := h.service.AddComment(slug, req.Comment.Body, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	slug := c.Param("slug")

	var currentUserID *uint
	if user, exists := c.Get("userID"); exists {
		uid := user.(uint)
		currentUserID = &uid
	}

	comments, err := h.service.GetComments(slug, currentUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	slug := c.Param("slug")
	commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := c.GetUint("userID")

	if err := h.service.DeleteComment(slug, uint(commentID), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
