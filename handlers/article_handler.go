package handlers

import (
	"Netlfy/services"
	"Netlfy/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ArticleHandler struct {
	service services.ArticleService
}

func NewArticleHandler(service services.ArticleService) *ArticleHandler {
	return &ArticleHandler{service}
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var currentUserID *uint
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if userID, err := utils.ParseToken(tokenString); err == nil {
			currentUserID = &userID
		}
	}

	articles, count, err := h.service.ListArticles(tag, author, favorited, limit, offset, currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":      articles,
		"articlesCount": count,
	})
}
